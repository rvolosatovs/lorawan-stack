// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package migrations

func init() {
	const forwards = `
		CREATE TABLE IF NOT EXISTS authorization_codes (
			id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			authorization_code   STRING(64) UNIQUE NOT NULL,
			client_id            UUID NOT NULL REFERENCES clients(id),
			created_at           TIMESTAMP NOT NULL DEFAULT current_timestamp(),
			expires_in           INTEGER NOT NULL,
			scope                STRING NOT NULL,
			redirect_uri         STRING NOT NULL,
			state                STRING NOT NULL,
			user_id              UUID NOT NULL REFERENCES users(id)
		);

		CREATE TABLE IF NOT EXISTS access_tokens (
			id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			access_token    STRING UNIQUE NOT NULL,
			client_id       UUID NOT NULL REFERENCES clients(id),
			user_id         UUID NOT NULL REFERENCES users(id),
			created_at      TIMESTAMP NOT NULL DEFAULT current_timestamp(),
			expires_in      INTEGER NOT NULL,
			scope           STRING NOT NULL,
			redirect_uri    STRING NOT NULL
		);

		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			refresh_token   STRING(64) UNIQUE NOT NULL,
			client_id       UUID NOT NULL REFERENCES clients(id),
			user_id         UUID NOT NULL REFERENCES users(id),
			created_at      TIMESTAMP NOT NULL DEFAULT current_timestamp(),
			scope           STRING NOT NULL,
			redirect_uri    STRING NOT NULL
		);
	`

	const backwards = `
		DROP TABLE IF EXISTS refresh_tokens;
		DROP TABLE IF EXISTS access_tokens;
		DROP TABLE IF EXISTS authorization_codes;
	`

	Registry.Register(5, "5_oauth_initial_schema", forwards, backwards)
}
