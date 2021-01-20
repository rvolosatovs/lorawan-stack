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

package emails

import (
	"fmt"

	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// APIKeyChanged is the email that is sent when users updates an API key
type APIKeyChanged struct {
	Data
	Key    *ttnpb.APIKey
	Rights []ttnpb.Right
}

// Identifier returns the pretty name of the API key.
// The naming of this method is for compatibility reasons.
func (a APIKeyChanged) Identifier() string {
	return a.Key.PrettyName()
}

// ConsoleURL returns the URL to the API key in the Console.
func (a APIKeyChanged) ConsoleURL() string {
	if a.Entity.Type == "user" {
		return fmt.Sprintf("%s/user/api-keys/%s", a.Network.ConsoleURL, a.Key.ID)
	}
	return fmt.Sprintf("%s/%ss/%s/api-keys/%s", a.Network.ConsoleURL, a.Entity.Type, a.Entity.ID, a.Key.ID)
}

// TemplateName returns the name of the template to use for this email.
func (APIKeyChanged) TemplateName() string { return "api_key_changed" }

const apiKeyChangedSubject = `An API key has been changed`

const apiKeyChangedText = `Dear {{.User.Name}},

The API key "{{.Identifier}}" for {{.Entity.Type}} "{{.Entity.ID}}" on {{.Network.Name}} has been updated with the following rights:
{{range $right := .Rights}} 
{{$right}} {{end}}

You can go to {{.ConsoleURL}} to view and edit this API key in the Console.

If you prefer to use the command-line interface, you can run the following commands to view or edit this API key:

ttn-lw-cli {{.Entity.Type}}s api-keys get --{{.Entity.Type}}-id {{.Entity.ID}} --api-key-id {{.Key.ID}}
ttn-lw-cli {{.Entity.Type}}s api-keys set --{{.Entity.Type}}-id {{.Entity.ID}} --api-key-id {{.Key.ID}}

For more information on how to use the command-line interface, please refer to the documentation: {{ documentation_url "/getting-started/cli/" }}.
`

// DefaultTemplates returns the default templates for this email.
func (APIKeyChanged) DefaultTemplates() (subject, html, text string) {
	return apiKeyChangedSubject, "", apiKeyChangedText
}
