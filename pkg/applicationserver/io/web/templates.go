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

package web

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/fetch"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"gopkg.in/yaml.v2"
)

const yamlFetchErrorCache = 1 * time.Minute

type noopTemplateStore struct {
}

var (
	errTemplateNotFound = errors.DefineNotFound("template_not_found", "template `{template_id}` not found")
)

// GetTemplate implements TemplateStore.
func (ts *noopTemplateStore) GetTemplate(ctx context.Context, req *ttnpb.GetApplicationWebhookTemplateRequest) (*ttnpb.ApplicationWebhookTemplate, error) {
	return nil, errTemplateNotFound.WithAttributes("template_id", req.TemplateID)
}

// ListTemplates implements TemplateStore.
func (ts *noopTemplateStore) ListTemplates(ctx context.Context, req *ttnpb.ListApplicationWebhookTemplatesRequest) (*ttnpb.ApplicationWebhookTemplates, error) {
	return &ttnpb.ApplicationWebhookTemplates{}, nil
}

// templateStore implements TemplateStore using an underlying fetcher.
type templateStore struct {
	fetcher fetch.Interface
	baseURL *url.URL

	templateIDs          []string
	templateIDsMu        sync.Mutex
	templateIDsError     error
	templateIDsErrorTime time.Time

	templates   map[string]queryResult
	templatesMu sync.Mutex
}

// prependBaseURL prepends the base URL and the template ID to the LogoURL, if it is available.
func (ts *templateStore) prependBaseURL(template *ttnpb.ApplicationWebhookTemplate) error {
	if template.LogoURL == "" {
		return nil
	}
	logoURL, err := url.Parse(template.LogoURL)
	if err != nil {
		return err
	}
	template.LogoURL = ts.baseURL.ResolveReference(logoURL).String()
	return nil
}

// GetTemplate implements the TemplateStore interface.
func (ts *templateStore) GetTemplate(ctx context.Context, req *ttnpb.GetApplicationWebhookTemplateRequest) (*ttnpb.ApplicationWebhookTemplate, error) {
	template, err := ts.getTemplate(req.ApplicationWebhookTemplateIdentifiers)
	if err != nil {
		return nil, err
	}
	template, err = applyWebhookTemplateFieldMask(nil, template, appendImplicitWebhookTemplatePaths(req.FieldMask.Paths...)...)
	if err != nil {
		return nil, err
	}
	err = ts.prependBaseURL(template)
	if err != nil {
		return nil, err
	}
	return template, nil
}

// ListTemplates implements the TemplateStore interface.
func (ts *templateStore) ListTemplates(ctx context.Context, req *ttnpb.ListApplicationWebhookTemplatesRequest) (*ttnpb.ApplicationWebhookTemplates, error) {
	ids, err := ts.getAllTemplateIDs()
	if err != nil {
		return nil, err
	}

	var templates ttnpb.ApplicationWebhookTemplates
	for _, id := range ids {
		template, err := ts.getTemplate(ttnpb.ApplicationWebhookTemplateIdentifiers{
			TemplateID: id,
		})
		if err != nil {
			return nil, err
		}

		template, err = applyWebhookTemplateFieldMask(nil, template, appendImplicitWebhookTemplatePaths(req.FieldMask.Paths...)...)
		if err != nil {
			return nil, err
		}

		err = ts.prependBaseURL(template)
		if err != nil {
			return nil, err
		}

		templates.Templates = append(templates.Templates, template)
	}
	return &templates, nil
}

type queryResult struct {
	t    *ttnpb.ApplicationWebhookTemplate
	err  error
	time time.Time
}

var (
	errFetchFailed = errors.Define("fetch", "fetching failed")
	errParseFile   = errors.DefineCorruption("parse_file", "could not parse file")
)

func (ts *templateStore) allTemplateIDs() (ids []string, err error) {
	data, err := ts.fetcher.File("templates.yml")
	if err != nil {
		return nil, errFetchFailed.WithCause(err)
	}
	err = yaml.Unmarshal(data, &ids)
	if err != nil {
		return nil, errParseFile.WithCause(err)
	}
	return ids, nil
}

func (ts *templateStore) getAllTemplateIDs() ([]string, error) {
	ts.templateIDsMu.Lock()
	defer ts.templateIDsMu.Unlock()
	if ts.templateIDs != nil {
		return ts.templateIDs, nil
	}
	if time.Since(ts.templateIDsErrorTime) < yamlFetchErrorCache {
		return nil, ts.templateIDsError
	}
	ids, err := ts.allTemplateIDs()
	if err != nil {
		ts.templateIDsError, ts.templateIDsErrorTime = err, time.Now()
		return nil, err
	}
	ts.templateIDs, ts.templateIDsError, ts.templateIDsErrorTime = ids, nil, time.Time{}
	return ids, err
}

func (ts *templateStore) template(ids ttnpb.ApplicationWebhookTemplateIdentifiers) (*ttnpb.ApplicationWebhookTemplate, error) {
	data, err := ts.fetcher.File(fmt.Sprintf("%s.yml", ids.TemplateID))
	if err != nil {
		return nil, errFetchFailed.WithCause(err)
	}
	var template webhookTemplate
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		return nil, errParseFile.WithCause(err)
	}
	return template.toPB(), nil
}

func (ts *templateStore) getTemplate(ids ttnpb.ApplicationWebhookTemplateIdentifiers) (t *ttnpb.ApplicationWebhookTemplate, err error) {
	ts.templatesMu.Lock()
	defer ts.templatesMu.Unlock()
	if cached, ok := ts.templates[ids.TemplateID]; ok && cached.err == nil && time.Since(cached.time) < yamlFetchErrorCache {
		return cached.t, cached.err
	}
	template, err := ts.template(ids)
	ts.templates[ids.TemplateID] = queryResult{
		t:    template,
		err:  err,
		time: time.Now(),
	}
	return template, err
}

func appendImplicitWebhookTemplatePaths(paths ...string) []string {
	return append(append(make([]string, 0, 2+len(paths)),
		"ids",
		"name",
	), paths...)
}

func applyWebhookTemplateFieldMask(dst, src *ttnpb.ApplicationWebhookTemplate, paths ...string) (*ttnpb.ApplicationWebhookTemplate, error) {
	if dst == nil {
		dst = &ttnpb.ApplicationWebhookTemplate{}
	}
	return dst, dst.SetFields(src, paths...)
}

type webhookTemplateField struct {
	ID           string `yaml:"id"`
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	Secret       bool   `yaml:"secret"`
	DefaultValue string `yaml:"default-value"`
}

func (f webhookTemplateField) toPB() *ttnpb.ApplicationWebhookTemplateField {
	return &ttnpb.ApplicationWebhookTemplateField{
		ID:           f.ID,
		Name:         f.Name,
		Description:  f.Description,
		Secret:       f.Secret,
		DefaultValue: f.DefaultValue,
	}
}

type webhookTemplatePaths struct {
	UplinkMessage            *string `yaml:"uplink-message,omitempty"`
	JoinAccept               *string `yaml:"join-accept,omitempty"`
	DownlinkAck              *string `yaml:"downlink-ack,omitempty"`
	DownlinkNack             *string `yaml:"downlink-nack,omitempty"`
	DownlinkSent             *string `yaml:"downlink-sent,omitempty"`
	DownlinkFailed           *string `yaml:"downlink-failed,omitempty"`
	DownlinkQueued           *string `yaml:"downlink-queued,omitempty"`
	DownlinkQueueInvalidated *string `yaml:"downlink-queue-invalidated,omitempty"`
	LocationSolved           *string `yaml:"location-solved,omitempty"`
	ServiceData              *string `yaml:"service-data,omitempty"`
}

type webhookTemplate struct {
	TemplateID           string                 `yaml:"template-id"`
	Name                 string                 `yaml:"name"`
	Description          string                 `yaml:"description"`
	LogoURL              string                 `yaml:"logo-url"`
	InfoURL              string                 `yaml:"info-url"`
	DocumentationURL     string                 `yaml:"documentation-url"`
	BaseURL              string                 `yaml:"base-url"`
	Headers              map[string]string      `yaml:"headers,omitempty"`
	Format               string                 `yaml:"format"`
	Fields               []webhookTemplateField `yaml:"fields,omitempty"`
	CreateDownlinkAPIKey bool                   `yaml:"create-downlink-api-key"`
	Paths                webhookTemplatePaths   `yaml:"paths,omitempty"`
}

func (webhookTemplate) pathToMessage(s *string) *ttnpb.ApplicationWebhookTemplate_Message {
	if s == nil {
		return nil
	}
	return &ttnpb.ApplicationWebhookTemplate_Message{
		Path: *s,
	}
}

func (t webhookTemplate) pbFields() []*ttnpb.ApplicationWebhookTemplateField {
	var fields []*ttnpb.ApplicationWebhookTemplateField
	for _, f := range t.Fields {
		fields = append(fields, f.toPB())
	}
	return fields
}

func (t webhookTemplate) toPB() *ttnpb.ApplicationWebhookTemplate {
	return &ttnpb.ApplicationWebhookTemplate{
		ApplicationWebhookTemplateIdentifiers: ttnpb.ApplicationWebhookTemplateIdentifiers{
			TemplateID: t.TemplateID,
		},
		Name:                     t.Name,
		Description:              t.Description,
		LogoURL:                  t.LogoURL,
		InfoURL:                  t.InfoURL,
		DocumentationURL:         t.DocumentationURL,
		BaseURL:                  t.BaseURL,
		Headers:                  t.Headers,
		Format:                   t.Format,
		Fields:                   t.pbFields(),
		CreateDownlinkAPIKey:     t.CreateDownlinkAPIKey,
		UplinkMessage:            t.pathToMessage(t.Paths.UplinkMessage),
		JoinAccept:               t.pathToMessage(t.Paths.JoinAccept),
		DownlinkAck:              t.pathToMessage(t.Paths.DownlinkAck),
		DownlinkNack:             t.pathToMessage(t.Paths.DownlinkNack),
		DownlinkSent:             t.pathToMessage(t.Paths.DownlinkSent),
		DownlinkFailed:           t.pathToMessage(t.Paths.DownlinkFailed),
		DownlinkQueued:           t.pathToMessage(t.Paths.DownlinkQueued),
		DownlinkQueueInvalidated: t.pathToMessage(t.Paths.DownlinkQueueInvalidated),
		LocationSolved:           t.pathToMessage(t.Paths.LocationSolved),
		ServiceData:              t.pathToMessage(t.Paths.ServiceData),
	}
}
