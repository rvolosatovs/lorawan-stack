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

package ttnmage

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/magefile/mage/mg"
)

// Dev namespace.
type Dev mg.Namespace

// Certificates generates certificates for development.
func (Dev) Certificates() error {
	if _, err := os.Stat("key.pem"); err == nil {
		if _, err := os.Stat("cert.pem"); err == nil {
			return nil
		}
	}
	return execGo("run", path.Join(runtime.GOROOT(), "src", "crypto", "tls", "generate_cert.go"), "-ca", "-host", "localhost,*.localhost")
}

// Misspell fixes common spelling mistakes in files.
func (Dev) Misspell() error {
	if mg.Verbose() {
		fmt.Printf("Fixing common spelling mistakes in files\n")
	}
	return execGo("run", "github.com/client9/misspell/cmd/misspell", "-w", "-i", "mosquitto",
		".editorconfig",
		".gitignore",
		".goreleaser.yml",
		".mage",
		".make",
		".revive.toml",
		".travis.yml",
		"api",
		"cmd",
		"config",
		"CONTRIBUTING.md",
		"DEVELOPMENT.md",
		"doc",
		"docker-compose.yml",
		"Dockerfile",
		"lorawan-stack.go",
		"magefile.go",
		"Makefile",
		"pkg",
		"README.md",
		"sdk",
		"SECURITY.md",
		"tools.go",
	)
}

// InitStack initializes the Stack.
func (Dev) InitStack() error {
	if mg.Verbose() {
		fmt.Printf("Initializing the Stack\n")
	}
	if err := execGo("run", "./cmd/ttn-lw-stack", "is-db", "init"); err != nil {
		return err
	}
	if err := execGo("run", "./cmd/ttn-lw-stack", "is-db", "create-admin-user",
		"--id", "admin",
		"--email", "admin@localhost",
		"--password", "admin",
	); err != nil {
		return err
	}
	if err := execGo("run", "./cmd/ttn-lw-stack", "is-db", "create-oauth-client",
		"--id", "cli",
		"--name", "Command Line Interface",
		"--owner", "admin",
		"--no-secret",
		"--redirect-uri", "local-callback",
		"--redirect-uri", "code",
	); err != nil {
		return err
	}
	return execGo("run", "./cmd/ttn-lw-stack", "is-db", "create-oauth-client",
		"--id", "console",
		"--name", "Console",
		"--owner", "admin",
		"--secret", "console",
		"--redirect-uri", "https://localhost:8885/console/oauth/callback",
		"--redirect-uri", "http://localhost:1885/console/oauth/callback",
		"--redirect-uri", "/console/oauth/callback",
	)
}

func init() {
	initDeps = append(initDeps, Dev.Certificates)
}
