// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package sql

import (
	"github.com/TheThingsNetwork/ttn/pkg/identityserver/db"
	"github.com/TheThingsNetwork/ttn/pkg/ttnpb"
	"github.com/satori/go.uuid"
)

func (s *OrganizationStore) SaveAPIKey(ids ttnpb.OrganizationIdentifiers, key ttnpb.APIKey) error {
	err := s.transact(func(tx *db.Tx) error {
		orgID, err := s.getOrganizationID(tx, ids)
		if err != nil {
			return err
		}

		err = s.saveAPIKey(tx, orgID, key)
		if err != nil {
			return err
		}

		return s.saveAPIKeyRights(tx, orgID, key.Name, key.Rights)
	})
	return err
}

func (s *OrganizationStore) GetAPIKey(value string) (ids ttnpb.OrganizationIdentifiers, key ttnpb.APIKey, err error) {
	err = s.transact(func(tx *db.Tx) error {
		var orgID uuid.UUID
		orgID, key, err = s.getAPIKey(tx, value)
		if err != nil {
			return err
		}

		key.Rights, err = s.getAPIKeyRights(tx, orgID, key.Name)
		if err != nil {
			return err
		}

		ids, err = s.getOrganizationIdentifiersFromID(tx, orgID)
		if err != nil {
			return err
		}

		return nil
	})
	return
}

func (s *OrganizationStore) GetAPIKeyByName(ids ttnpb.OrganizationIdentifiers, keyName string) (key ttnpb.APIKey, err error) {
	err = s.transact(func(tx *db.Tx) error {
		orgID, err := s.getOrganizationID(tx, ids)
		if err != nil {
			return err
		}

		key, err = s.getAPIKeyByName(tx, orgID, keyName)
		if err != nil {
			return err
		}

		key.Rights, err = s.getAPIKeyRights(tx, orgID, keyName)

		return err
	})
	return
}

func (s *OrganizationStore) ListAPIKeys(ids ttnpb.OrganizationIdentifiers) (keys []ttnpb.APIKey, err error) {
	err = s.transact(func(tx *db.Tx) error {
		orgID, err := s.getOrganizationID(tx, ids)
		if err != nil {
			return err
		}

		keys, err = s.listAPIKeys(tx, orgID)
		if err != nil {
			return err
		}

		for i, key := range keys {
			keys[i].Rights, err = s.getAPIKeyRights(tx, orgID, key.Name)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return
}

func (s *OrganizationStore) UpdateAPIKeyRights(ids ttnpb.OrganizationIdentifiers, keyName string, rights []ttnpb.Right) error {
	err := s.transact(func(tx *db.Tx) error {
		orgID, err := s.getOrganizationID(tx, ids)
		if err != nil {
			return err
		}

		err = s.deleteAPIKeyRights(tx, orgID, keyName)
		if err != nil {
			return err
		}

		return s.saveAPIKeyRights(tx, orgID, keyName, rights)
	})
	return err
}

func (s *OrganizationStore) DeleteAPIKey(ids ttnpb.OrganizationIdentifiers, keyName string) error {
	err := s.transact(func(tx *db.Tx) error {
		orgID, err := s.getOrganizationID(tx, ids)
		if err != nil {
			return err
		}

		err = s.deleteAPIKeyRights(tx, orgID, keyName)
		if err != nil {
			return err
		}

		return s.deleteAPIKey(tx, orgID, keyName)
	})
	return err
}
