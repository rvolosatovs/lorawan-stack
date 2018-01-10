// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package sql

import (
	"github.com/TheThingsNetwork/ttn/pkg/errors"
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/db"
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/store"
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/types"
)

// UserStore implements store.UserStore.
type UserStore struct {
	storer
	*extraAttributesStore
	*apiKeysStore
}

func NewUserStore(store storer) *UserStore {
	return &UserStore{
		storer:               store,
		extraAttributesStore: newExtraAttributesStore(store, "user"),
		apiKeysStore:         newAPIKeysStore(store, "user"),
	}
}

// Create creates an user.
func (s *UserStore) Create(user types.User) error {
	err := s.transact(func(tx *db.Tx) error {
		err := s.create(tx, user)
		if err != nil {
			return err
		}

		return s.storeAttributes(tx, user.GetUser().UserID, user, nil)
	})
	return err
}

func (s *UserStore) create(q db.QueryContext, user types.User) error {
	u := user.GetUser()
	_, err := q.NamedExec(
		`INSERT
			INTO users (
				user_id,
				display_name,
				email,
				admin,
				state,
				password,
				validated_at)
			VALUES (
				lower(:user_id),
				:display_name,
				lower(:email),
				:admin,
				:state,
				:password,
				:validated_at)`,
		u)

	if duplicates, yes := db.IsDuplicate(err); yes {
		if email, duplicated := duplicates["email"]; duplicated {
			return ErrUserEmailTaken.New(errors.Attributes{
				"email": email,
			})
		}
		if userID, duplicated := duplicates["user_id"]; duplicated {
			return ErrUserIDTaken.New(errors.Attributes{
				"user_id": userID,
			})
		}
	}

	return err
}

// GetByID finds the user by ID and returns it.
func (s *UserStore) GetByID(userID string, factory store.UserFactory) (types.User, error) {
	result := factory()

	err := s.transact(func(tx *db.Tx) error {
		err := s.getByID(tx, userID, result)
		if err != nil {
			return err
		}

		return s.loadAttributes(tx, result.GetUser().UserID, result)
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *UserStore) getByID(q db.QueryContext, userID string, result types.User) error {
	err := q.SelectOne(
		result,
		`SELECT *
			FROM users
			WHERE user_id = lower($1)`,
		userID)
	if db.IsNoRows(err) {
		return ErrUserNotFound.New(errors.Attributes{
			"user_id": userID,
		})
	}
	return err
}

// GetByEmail finds the user by email address and returns it.
func (s *UserStore) GetByEmail(email string, factory store.UserFactory) (types.User, error) {
	result := factory()

	err := s.transact(func(tx *db.Tx) error {
		err := s.getByEmail(tx, email, result)
		if err != nil {
			return err
		}

		return s.loadAttributes(tx, result.GetUser().UserID, result)
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *UserStore) getByEmail(q db.QueryContext, email string, result types.User) error {
	err := q.SelectOne(
		result,
		`SELECT *
			FROM users
			WHERE email = lower($1)`,
		email)
	if db.IsNoRows(err) {
		return ErrUserEmailNotFound.New(errors.Attributes{
			"email": email,
		})
	}
	return err
}

// Update updates an user.
func (s *UserStore) Update(user types.User) error {
	err := s.transact(func(tx *db.Tx) error {
		err := s.update(tx, user)
		if err != nil {
			return err
		}

		return s.storeAttributes(s.queryer(), user.GetUser().UserID, user, nil)
	})
	return err
}

func (s *UserStore) update(q db.QueryContext, user types.User) error {
	u := user.GetUser()

	_, err := q.NamedExec(
		`UPDATE users
			SET
				display_name = :display_name,
				email = lower(:email),
				validated_at = :validated_at,
				password = :password,
				admin = :admin,
				state = :state,
				updated_at = current_timestamp()
			WHERE user_id = :user_id`,
		u)

	if _, yes := db.IsDuplicate(err); yes {
		return ErrUserEmailTaken.New(errors.Attributes{
			"email": u.Email,
		})
	}

	return err
}

// Delete deletes an user.
func (s *UserStore) Delete(userID string) error {
	err := s.transact(func(tx *db.Tx) error {
		// unset as application collaborator where possible
		applications, ok := s.store().Applications.(*ApplicationStore)
		if !ok {
			return errors.Errorf("Expected ptr to ApplicationStore but got %T", s.store().Applications)
		}

		apps, err := applications.userApplications(tx, userID)
		if err != nil {
			return err
		}

		for _, app := range apps {
			err := applications.unsetCollaborator(tx, app.ApplicationID, userID)
			if err != nil {
				return err
			}
		}

		// unset as gateway collaborator where possible
		gateways, ok := s.store().Gateways.(*GatewayStore)
		if !ok {
			return errors.Errorf("Expected ptr to GatewayStore but got %T", s.store().Gateways)
		}

		gtws, err := gateways.userGateways(tx, userID)
		if err != nil {
			return err
		}

		for _, gtw := range gtws {
			err := gateways.unsetCollaborator(tx, gtw.GatewayID, userID)
			if err != nil {
				return err
			}
		}

		// revoke all authorized clients
		oauth, ok := s.store().OAuth.(*OAuthStore)
		if !ok {
			return errors.Errorf("Expected ptr to OAuthStore but got %T", s.store().OAuth)
		}

		err = oauth.deleteAuthorizationCodesByUser(tx, userID)
		if err != nil {
			return err
		}

		clientIDs, err := oauth.listAuthorizedClients(tx, userID)
		for _, clientID := range clientIDs {
			_, err := oauth.deleteAccessTokensByUserAndClient(tx, userID, clientID)
			if err != nil {
				return err
			}

			_, err = oauth.deleteRefreshTokenByUserAndClient(tx, userID, clientID)
			if err != nil {
				return err
			}
		}

		// delete created clients
		clients, ok := s.store().Clients.(*ClientStore)
		if !ok {
			return errors.Errorf("Expected ptr to ClientStore but got %T", s.store().Clients)
		}

		found, err := clients.userClients(tx, userID)
		if err != nil {
			return err
		}

		for _, client := range found {
			err := clients.delete(tx, client.GetClient().ClientID)
			if err != nil {
				return err
			}
		}

		// delete api keys
		err = s.deleteAPIKeys(tx, userID)
		if err != nil {
			return err
		}

		// delete validation tokens
		err = s.deleteValidationTokens(tx, userID)
		if err != nil {
			return err
		}

		return s.delete(tx, userID)
	})

	return err
}

func (s *UserStore) delete(q db.QueryContext, userID string) error {
	id := new(string)
	err := q.SelectOne(
		id,
		`DELETE
			FROM users
			WHERE user_id = $1
			RETURNING user_id`,
		userID)
	if db.IsNoRows(err) {
		return ErrUserNotFound.New(errors.Attributes{
			"user_id": userID,
		})
	}
	return err
}

func (s *UserStore) deleteValidationTokens(q db.QueryContext, userID string) error {
	_, err := q.Exec(`DELETE FROM validation_tokens WHERE user_id = $1`, userID)
	return err
}

// SaveValidationToken saves the validation token.
func (s *UserStore) SaveValidationToken(userID string, token *types.ValidationToken) error {
	return s.saveValidationToken(s.queryer(), userID, token)
}

func (s *UserStore) saveValidationToken(q db.QueryContext, userID string, token *types.ValidationToken) error {
	_, err := q.Exec(
		`INSERT
			INTO validation_tokens (
				validation_token,
				user_id,
				created_at,
				expires_in
			)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id)
			DO UPDATE SET
				validation_token = excluded.validation_token,
				created_at = excluded.created_at,
				expires_in = excluded.expires_in`,
		token.ValidationToken,
		userID,
		token.CreatedAt,
		token.ExpiresIn)
	return err
}

// GetValidationToken retrieves the validation token.
func (s *UserStore) GetValidationToken(token string) (string, *types.ValidationToken, error) {
	return s.getValidationToken(s.queryer(), token)
}

func (s *UserStore) getValidationToken(q db.QueryContext, token string) (string, *types.ValidationToken, error) {
	var t struct {
		*types.ValidationToken
		UserID string
	}
	err := q.SelectOne(
		&t,
		`SELECT
				user_id,
				validation_token,
				created_at,
				expires_in
			FROM validation_tokens
			WHERE validation_token = $1`,
		token)
	if db.IsNoRows(err) {
		return "", nil, ErrValidationTokenNotFound.New(nil)
	}
	if err != nil {
		return "", nil, err
	}
	return t.UserID, t.ValidationToken, nil
}

// DeleteValidationToken deletes the validation token.
func (s *UserStore) DeleteValidationToken(token string) error {
	return s.deleteValidationToken(s.queryer(), token)
}

func (s *UserStore) deleteValidationToken(q db.QueryContext, token string) error {
	t := new(string)
	err := q.SelectOne(
		t,
		`DELETE
			FROM validation_tokens
			WHERE validation_token = $1
			RETURNING validation_token`,
		token)
	if db.IsNoRows(err) {
		return ErrValidationTokenNotFound.New(nil)
	}
	return err
}

// LoadAttributes loads the extra attributes in user if it is a store.Attributer.
func (s *UserStore) LoadAttributes(userID string, user types.User) error {
	return s.loadAttributes(s.queryer(), userID, user)
}

func (s *UserStore) loadAttributes(q db.QueryContext, userID string, user types.User) error {
	attr, ok := user.(store.Attributer)
	if ok {
		return s.extraAttributesStore.loadAttributes(q, userID, attr)
	}

	return nil
}

// StoreAttributes store the extra attributes of user if it is a store.Attributer
// and writes the resulting user in result.
func (s *UserStore) StoreAttributes(userID string, user, result types.User) error {
	return s.storeAttributes(s.queryer(), userID, user, result)
}

func (s *UserStore) storeAttributes(q db.QueryContext, userID string, user, result types.User) error {
	attr, ok := user.(store.Attributer)
	if ok {
		res, ok := result.(store.Attributer)
		if result == nil || !ok {
			return s.extraAttributesStore.storeAttributes(q, userID, attr, nil)
		}

		return s.extraAttributesStore.storeAttributes(q, userID, attr, res)
	}

	return nil
}
