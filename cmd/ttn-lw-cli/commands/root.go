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

// Package commands implements the commands for the ttn-lw-cli binary.
package commands

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.thethings.network/lorawan-stack/cmd/internal/commands"
	"go.thethings.network/lorawan-stack/cmd/internal/shared/version"
	"go.thethings.network/lorawan-stack/cmd/ttn-lw-cli/internal/api"
	"go.thethings.network/lorawan-stack/cmd/ttn-lw-cli/internal/io"
	"go.thethings.network/lorawan-stack/cmd/ttn-lw-cli/internal/util"
	conf "go.thethings.network/lorawan-stack/pkg/config"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/log"
	"golang.org/x/oauth2"
)

var (
	logger       *log.Logger
	name         = "ttn-lw-cli"
	mgr          = conf.InitializeWithDefaults(name, "ttn_lw", DefaultConfig)
	config       = &Config{}
	oauth2Config *oauth2.Config
	ctx          = newContext(context.Background())
	cache        util.AuthCache

	inputDecoder io.Decoder

	// Root command is the entrypoint of the program
	Root = &cobra.Command{
		Use:               name,
		SilenceErrors:     true,
		SilenceUsage:      true,
		Short:             "The Things Network Command-line Interface",
		PersistentPreRunE: preRun(checkAuth, refreshToken, requireAuth),
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			// clean up the API
			api.CloseAll()

			err := util.SaveAuthCache(cache)
			if err != nil {
				return err
			}

			return ctx.Err()
		},
	}
)

func preRun(tasks ...func() error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// read in config from file
		err := mgr.ReadInConfig()
		if err != nil {
			return err
		}

		// unmarshal config
		if err = mgr.Unmarshal(config); err != nil {
			return err
		}

		// create input decoder on Stdin
		if io.IsPipe(os.Stdin) {
			inputDecoder, err = getInputDecoder(os.Stdin)
			if err != nil {
				return err
			}
		}

		// get cache
		cache, err = util.GetAuthCache()
		if err != nil {
			return err
		}
		cache = cache.ForID(config.CredentialsID)

		// create logger
		logger, err = log.NewLogger(
			log.WithLevel(config.Log.Level),
			log.WithHandler(log.NewCLI(os.Stderr)),
		)
		if err != nil {
			return err
		}
		ctx = log.NewContext(ctx, logger)

		// prepare the API
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{}
		api.SetLogger(logger)
		if config.Insecure {
			api.SetInsecure(true)
		}
		if config.DumpRequests {
			api.SetDumpRequests(true)
		}
		if config.CA != "" {
			pemBytes, err := ioutil.ReadFile(config.CA)
			if err != nil {
				return err
			}
			rootCAs := http.DefaultTransport.(*http.Transport).TLSClientConfig.RootCAs
			if rootCAs == nil {
				if rootCAs, err = x509.SystemCertPool(); err != nil {
					rootCAs = x509.NewCertPool()
				}
			}
			rootCAs.AppendCertsFromPEM(pemBytes)
			http.DefaultTransport.(*http.Transport).TLSClientConfig.RootCAs = rootCAs
			if err = api.AddCA(pemBytes); err != nil {
				return err
			}
		}

		// OAuth
		oauth2Config = &oauth2.Config{
			ClientID: "cli",
			Endpoint: oauth2.Endpoint{
				AuthURL:   fmt.Sprintf("%s/authorize", config.OAuthServerAddress),
				TokenURL:  fmt.Sprintf("%s/token", config.OAuthServerAddress),
				AuthStyle: oauth2.AuthStyleInParams,
			},
		}

		for _, task := range tasks {
			if err := task(); err != nil {
				return err
			}
		}
		return nil
	}
}

var errUnknownHost = errors.DefineUnauthenticated("unknown_host", "unknown host `{host}` for current credentials", "known")

func checkAuth() error {
	if config.AllowUnknownHosts {
		return nil
	}
	if knownHosts, ok := cache.Get("hosts").([]string); ok && len(knownHosts) > 0 {
	nextHost:
		for _, host := range config.getHosts() {
			for _, knownHost := range knownHosts {
				if host == knownHost {
					continue nextHost
				}
			}
			logger.Errorf("Found an unknown host `%s` that was not configured when you logged in", host)
			logger.Error("You may want to check your configuration, login/logout or use the --allow-unknown-hosts flag")
			return errUnknownHost.WithAttributes("host", host, "known", knownHosts)
		}
	}
	return nil
}

func refreshToken() error {
	if token, ok := cache.Get("oauth_token").(*oauth2.Token); ok && token != nil {
		freshToken, err := oauth2Config.TokenSource(ctx, token).Token()
		if err == nil && freshToken != token {
			cache.Set("oauth_token", freshToken)
			if err := util.SaveAuthCache(cache); err != nil {
				return err
			}
		}
		return err
	}
	return nil
}

var errUnauthenticated = errors.DefineUnauthenticated("unauthenticated", "not authenticated with either API key or OAuth access token")

func optionalAuth() error {
	err := requireAuth()
	if err != nil && !errors.IsUnauthenticated(err) {
		return err
	}
	return nil
}

func requireAuth() error {
	if apiKey, ok := cache.Get("api_key").(string); ok && apiKey != "" {
		logger.Debug("Using API key")
		api.SetAuth("bearer", apiKey)
		return nil
	}
	if token, ok := cache.Get("oauth_token").(*oauth2.Token); ok && token != nil {
		friendlyExpiry := token.Expiry.Truncate(time.Minute).Format(time.Kitchen)
		if time.Until(token.Expiry) > 0 {
			logger.Debugf("Using access token (valid until %s)", friendlyExpiry)
			api.SetAuth(token.TokenType, token.AccessToken)
			return nil
		}
		logger.Warnf("Access token expired at %s", friendlyExpiry)
	}
	logger.Error("Please login with the login command")
	return errUnauthenticated
}

var (
	versionCommand     = version.Print(Root)
	genManPagesCommand = commands.GenManPages(Root)
	genMDDocCommand    = commands.GenMDDoc(Root)
	ganYAMLDocCommand  = commands.GenYAMLDoc(Root)
)

func init() {
	Root.SetGlobalNormalizationFunc(util.NormalizeFlags)
	Root.PersistentFlags().AddFlagSet(mgr.Flags())
	versionCommand.PersistentPreRunE = preRun()
	Root.AddCommand(versionCommand)
	genManPagesCommand.PersistentPreRunE = preRun()
	Root.AddCommand(genManPagesCommand)
	genMDDocCommand.PersistentPreRunE = preRun()
	Root.AddCommand(genMDDocCommand)
	ganYAMLDocCommand.PersistentPreRunE = preRun()
	Root.AddCommand(ganYAMLDocCommand)
}
