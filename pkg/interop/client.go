// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interop

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"time"

	"go.thethings.network/lorawan-stack/pkg/config"
	"go.thethings.network/lorawan-stack/pkg/encoding/lorawan"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/types"
)

const (
	// loRaAllianceDomain is the domain of LoRa Alliance.
	loRaAllianceDomain = "lora-alliance.org"

	// LoRaAllianceJoinEUIDomain is the LoRa Alliance domain used for JoinEUI resolution.
	LoRaAllianceJoinEUIDomain = "joineuis." + loRaAllianceDomain

	// LoRaAllianceNetIDDomain is the LoRa Alliance domain used for NetID resolution.
	LoRaAllianceNetIDDomain = "netids." + loRaAllianceDomain

	defaultHTTPSPort = 443
)

type JoinServerProtocol uint8

const (
	LoRaWANJoinServerProtocol1_0 = iota
	LoRaWANJoinServerProtocol1_1
)

func (p JoinServerProtocol) BackendInterfacesVersion() string {
	switch p {
	case LoRaWANJoinServerProtocol1_0:
		return "1.0"
	case LoRaWANJoinServerProtocol1_1:
		return "1.1"
	default:
		panic(fmt.Sprintf("Join Server protocol	`%v` is not compliant with Backend Interfaces specification", p))
	}
}

func serverURL(scheme, fqdn, path string, port uint32) string {
	if scheme == "" {
		scheme = "https"
	}
	if port == 0 {
		port = defaultHTTPSPort
	}
	if path != "" {
		path = fmt.Sprintf("/%s", path)
	}
	return fmt.Sprintf("%s://%s:%d%s", scheme, fqdn, port, path)
}

func newHTTPRequest(url string, pld interface{}, headers map[string]string) (*http.Request, error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(pld); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

func JoinServerFQDN(eui types.EUI64, domain string) string {
	if domain == "" {
		domain = LoRaAllianceJoinEUIDomain
	}
	return fmt.Sprintf(
		"%01x.%01x.%01x.%01x.%01x.%01x.%01x.%01x.%01x.%01x.%01x.%01x.%01x.%01x.%01x.%01x.%s",
		eui[7]&0x0f, eui[7]>>4,
		eui[6]&0x0f, eui[6]>>4,
		eui[5]&0x0f, eui[5]>>4,
		eui[4]&0x0f, eui[4]>>4,
		eui[3]&0x0f, eui[3]>>4,
		eui[2]&0x0f, eui[2]>>4,
		eui[1]&0x0f, eui[1]>>4,
		eui[0]&0x0f, eui[0]>>4,
		domain,
	)
}

type joinServerHTTPClient struct {
	Client         http.Client
	NewRequestFunc func(joinEUI types.EUI64, pld interface{}) (*http.Request, error)
	Protocol       JoinServerProtocol
}

func (cl joinServerHTTPClient) exchange(joinEUI types.EUI64, req, res interface{}) error {
	httpReq, err := cl.NewRequestFunc(joinEUI, req)
	if err != nil {
		return err
	}

	httpRes, err := cl.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()
	return json.NewDecoder(httpRes.Body).Decode(res)
}

func parseResult(r Result) error {
	if r.ResultCode == ResultSuccess {
		return nil
	}

	err, ok := resultErrors[r.ResultCode]
	if ok {
		return err
	}
	return errUnexpectedResult.WithAttributes("code", r.ResultCode)
}

func (cl joinServerHTTPClient) GetAppSKey(ctx context.Context, asID string, req *ttnpb.SessionKeyRequest) (*ttnpb.AppSKeyResponse, error) {
	interopAns := &AppSKeyAns{}
	if err := cl.exchange(req.JoinEUI, &AppSKeyReq{
		AsJsMessageHeader: AsJsMessageHeader{
			MessageHeader: MessageHeader{
				ProtocolVersion: cl.Protocol.BackendInterfacesVersion(),
				MessageType:     MessageTypeAppSKeyReq,
			},
			SenderID:   asID,
			ReceiverID: EUI64(req.JoinEUI),
		},
		DevEUI:       EUI64(req.DevEUI),
		SessionKeyID: Buffer(req.SessionKeyID),
	}, interopAns); err != nil {
		return nil, err
	}
	if err := parseResult(interopAns.Result); err != nil {
		return nil, err
	}

	return &ttnpb.AppSKeyResponse{
		AppSKey: ttnpb.KeyEnvelope(interopAns.AppSKey),
	}, nil
}

func (cl joinServerHTTPClient) HandleJoinRequest(ctx context.Context, netID types.NetID, req *ttnpb.JoinRequest) (*ttnpb.JoinResponse, error) {
	pld := req.Payload.GetJoinRequestPayload()
	if pld == nil {
		return nil, ErrMalformedMessage
	}

	dlSettings, err := lorawan.MarshalDLSettings(req.DownlinkSettings)
	if err != nil {
		return nil, err
	}

	var cfList []byte
	if req.CFList != nil {
		cfList, err = lorawan.MarshalCFList(*req.CFList)
		if err != nil {
			return nil, err
		}
	}

	interopAns := &JoinAns{}
	if err := cl.exchange(pld.JoinEUI, &JoinReq{
		NsJsMessageHeader: NsJsMessageHeader{
			MessageHeader: MessageHeader{
				ProtocolVersion: cl.Protocol.BackendInterfacesVersion(),
				MessageType:     MessageTypeJoinReq,
			},
			SenderID:   NetID(netID),
			ReceiverID: EUI64(pld.JoinEUI),
			SenderNSID: NetID(netID),
		},
		MACVersion: MACVersion(req.SelectedMACVersion),
		PHYPayload: Buffer(req.RawPayload),
		DevEUI:     EUI64(pld.DevEUI),
		DevAddr:    DevAddr(req.DevAddr),
		DLSettings: Buffer(dlSettings),
		RxDelay:    req.RxDelay,
		CFList:     Buffer(cfList),
	}, interopAns); err != nil {
		return nil, err
	}
	if err := parseResult(interopAns.Result); err != nil {
		return nil, err
	}

	fNwkSIntKey := interopAns.FNwkSIntKey
	if req.SelectedMACVersion.Compare(ttnpb.MAC_V1_1) <= 0 {
		fNwkSIntKey = interopAns.NwkSKey
	}
	return &ttnpb.JoinResponse{
		RawPayload: interopAns.PHYPayload,
		SessionKeys: ttnpb.SessionKeys{
			SessionKeyID: []byte(interopAns.SessionKeyID),
			FNwkSIntKey:  (*ttnpb.KeyEnvelope)(fNwkSIntKey),
			SNwkSIntKey:  (*ttnpb.KeyEnvelope)(interopAns.SNwkSIntKey),
			NwkSEncKey:   (*ttnpb.KeyEnvelope)(interopAns.NwkSEncKey),
			AppSKey:      (*ttnpb.KeyEnvelope)(interopAns.AppSKey),
		},
		Lifetime: time.Duration(interopAns.Lifetime) * time.Second,
	}, nil
}

type RemoteServerConfig struct {
	DNS     string            `name:"dns" description:"Domain name under which server address will be resolved by according to LoRaWAN Backend Interfaces 1.1 specification. LoRa Alliance domain is used if unset. FQDN takes precedence if set"`
	FQDN    string            `name:"fqdn" description:"FQDN of the server, DNS lookup is performed if not set"`
	Path    string            `name:"path" description:"URL path to use without the leading slash. Defaults to the paths specified in LoRaWAN Backend Interfaces 1.1 specification"`
	Port    uint32            `name:"port" description:"Port to use, defaults to 443"`
	TLS     config.TLS        `name:"tls" description:"TLS configuration to use"`
	Headers map[string]string `name:"headers" description:"Custom HTTP headers to send as part of the request"`
}

func makeJoinServerHTTPRequestFunc(scheme string, conf RemoteServerConfig) func(types.EUI64, interface{}) (*http.Request, error) {
	port := conf.Port
	if port == 0 {
		port = defaultHTTPSPort
	}
	path := conf.Path
	if path != "" {
		path = fmt.Sprintf("/%s", path)
	}
	return func(joinEUI types.EUI64, pld interface{}) (*http.Request, error) {
		fqdn := conf.FQDN
		if fqdn == "" {
			fqdn = JoinServerFQDN(joinEUI, conf.DNS)
		}
		return newHTTPRequest(serverURL(scheme, fqdn, path, port), pld, conf.Headers)
	}
}

type RemoteJoinServerConfig struct {
	RemoteServerConfig

	JoinEUI  []types.EUI64Prefix `name:"join-eui" description:"List of JoinEUI prefixes Join Server handles"`
	Protocol JoinServerProtocol  `name:"protocol" description:"Join Server protocol to use"`
}

var errUnknownProtocol = errors.DefineInvalidArgument("unknown_protocol", "unknown protocol")

type ClientConfig struct {
	JoinServers []RemoteJoinServerConfig `name:"join-servers" description:"List of Join Servers configured by JoinEUI."`
}

type joinServerClient interface {
	HandleJoinRequest(ctx context.Context, netID types.NetID, req *ttnpb.JoinRequest) (*ttnpb.JoinResponse, error)
	GetAppSKey(ctx context.Context, asID string, req *ttnpb.SessionKeyRequest) (*ttnpb.AppSKeyResponse, error)
}

type prefixJoinServerClient struct {
	joinServerClient
	prefix types.EUI64Prefix
}

type Client struct {
	joinServers []prefixJoinServerClient // Sorted by JoinEUI prefix range length.
}

func NewClient(ctx context.Context, conf ClientConfig, fallbackTLS *tls.Config) (*Client, error) {
	jss := make([]prefixJoinServerClient, 0, len(conf.JoinServers))
	for _, jsConf := range conf.JoinServers {
		var js joinServerClient
		switch jsConf.Protocol {
		case LoRaWANJoinServerProtocol1_0, LoRaWANJoinServerProtocol1_1:
			tlsConfig := fallbackTLS
			if !reflect.DeepEqual(jsConf.TLS, config.TLS{}) {
				var err error
				tlsConfig, err = jsConf.TLS.Config(ctx)
				if err != nil {
					return nil, err
				}
			}

			var tr *http.Transport
			if tlsConfig != nil {
				tr = &http.Transport{
					TLSClientConfig: tlsConfig,
				}
			}
			js = &joinServerHTTPClient{
				Client: http.Client{
					Transport: tr,
				},
				NewRequestFunc: makeJoinServerHTTPRequestFunc("https", jsConf.RemoteServerConfig),
				Protocol:       jsConf.Protocol,
			}
		default:
			return nil, errUnknownProtocol
		}
		for _, pre := range jsConf.JoinEUI {
			jss = append(jss, prefixJoinServerClient{
				joinServerClient: js,
				prefix:           pre,
			})
		}
	}
	sort.Slice(jss, func(i, j int) bool {
		pi, pj := jss[i].prefix, jss[j].prefix
		if pi.Length != pj.Length {
			return pi.Length > pj.Length
		}
		return pi.EUI64.MarshalNumber() > pj.EUI64.MarshalNumber()
	})
	return &Client{
		joinServers: jss,
	}, nil
}

func (cl Client) joinServer(joinEUI types.EUI64) (joinServerClient, bool) {
	// NOTE: joinServers slice is sorted by prefix length and the range start decreasing, hence the first match is the most specific one.
	for _, js := range cl.joinServers {
		fmt.Println(js.prefix.EUI64, js.prefix.Length, js.prefix.EUI64.Mask(js.prefix.Length))
		if js.prefix.Matches(joinEUI) {
			return js.joinServerClient, true
		}
	}
	return nil, false
}

func (cl Client) GetAppSKey(ctx context.Context, asID string, req *ttnpb.SessionKeyRequest) (*ttnpb.AppSKeyResponse, error) {
	js, ok := cl.joinServer(req.JoinEUI)
	if !ok {
		return nil, errNotRegistered
	}
	return js.GetAppSKey(ctx, asID, req)
}

func (cl Client) HandleJoinRequest(ctx context.Context, netID types.NetID, req *ttnpb.JoinRequest) (*ttnpb.JoinResponse, error) {
	pld := req.Payload.GetJoinRequestPayload()
	if pld == nil {
		return nil, ErrMalformedMessage
	}
	js, ok := cl.joinServer(pld.JoinEUI)
	if !ok {
		return nil, errNotRegistered
	}
	return js.HandleJoinRequest(ctx, netID, req)
}
