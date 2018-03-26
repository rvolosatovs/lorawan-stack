// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package migrations

func init() {
	// The check constraint in the column `id` enforces the table to have only
	// one row at most.
	const forwards = `
		CREATE TABLE IF NOT EXISTS settings (
			id                     INT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
			updated_at             TIMESTAMP NOT NULL DEFAULT current_timestamp(),
			blacklisted_ids        STRING NOT NULL DEFAULT '[]',
			skip_validation        BOOL NOT NULL DEFAULT false,
			invitation_only        BOOL NOT NULL DEFAULT false,
			admin_approval         BOOL NOT NULL DEFAULT false,
			validation_token_ttl   INT NOT NULL,
			allowed_emails         STRING NOT NULL DEFAULT '[]',
			invitation_token_ttl   INT NOT NULL
		);
	`
	const backwards = `
		DROP TABLE IF EXISTS settings;
	`

	Registry.Register(6, "6_settings_initial_schema", forwards, backwards)
}
