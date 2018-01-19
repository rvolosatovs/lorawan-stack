// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package templates

// EmailValidation is the email template used to validate an email address.
type Invitation struct {
	PublicURL        string
	OrganizationName string
	WebUIURL         string
	Token            string
}

// Name implements Template.
func (t *Invitation) GetName() string {
	return "Invitation"
}

// Render implements Template.
func (t *Invitation) Render() (string, string, error) {
	subject := "You had been invited to join {{.OrganizationName}}!"
	message := `<h1>Invitation</h1>

<p>
	You just got invited to create an account
	at <a href='{{.PublicURL}}'>{{.OrganizationName}}</a>.
</p>

<p>
	You can create your account by
	clicking the button below.
</p>

<p>
	<a class='button' href='{{.WebUIURL}}/register?token={{.Token}}'>Create account</a>
</p>`

	return render(subject, message, t)
}
