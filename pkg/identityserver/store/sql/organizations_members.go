// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
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

package sql

import (
	"fmt"
	"strings"

	"github.com/satori/go.uuid"
	"go.thethings.network/lorawan-stack/pkg/identityserver/db"
	"go.thethings.network/lorawan-stack/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

// ListByUser returns the organizations to which an user is a member of.
func (s *organizationStore) ListByUser(ids ttnpb.UserIdentifiers, specializer store.OrganizationSpecializer) (result []store.Organization, err error) {
	err = s.transact(func(tx *db.Tx) error {
		st := s.store()

		userID, err := st.Users.getUserID(tx, ids)
		if err != nil {
			return err
		}

		organizations, err := s.listUserOrganizations(tx, userID)
		if err != nil {
			return err
		}

		for _, organization := range organizations {
			specialized := specializer(organization.Organization)

			err := s.loadAttributes(tx, organization.ID, specialized)
			if err != nil {
				return err
			}

			result = append(result, specialized)
		}

		return nil
	})
	return
}

func (s *organizationStore) listUserOrganizations(q db.QueryContext, userID uuid.UUID) (organizations []organization, err error) {
	err = q.Select(
		&organizations,
		`SELECT DISTINCT
				organizations.*
			FROM organizations
			INNER JOIN organizations_members
			ON (organizations.id = organizations_members.organization_id)
			WHERE user_id = $1`,
		userID)
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

// HasMemberRights checks whether an user has or not a set of given rights to
// an organization. Returns false if the user is not part of the organization.
func (s *organizationStore) HasMemberRights(ids ttnpb.OrganizationIdentifiers, target ttnpb.UserIdentifiers, rights ...ttnpb.Right) (result bool, err error) {
	err = s.transact(func(tx *db.Tx) error {
		orgID, err := s.getOrganizationID(tx, ids)
		if err != nil {
			return err
		}

		userID, err := s.store().Users.getUserID(tx, target)
		if err != nil {
			return err
		}

		result, err = s.hasMemberRights(tx, orgID, userID, rights...)
		if err != nil {
			return err
		}

		return nil
	})
	return
}

func (s *organizationStore) hasMemberRights(q db.QueryContext, orgID, userID uuid.UUID, rights ...ttnpb.Right) (bool, error) {
	clauses := make([]string, 0, len(rights))
	args := make([]interface{}, 0, len(rights)+2)
	args = append(args, orgID, userID)

	for i, right := range rights {
		args = append(args, right)
		clauses = append(clauses, fmt.Sprintf(`"right" = $%d`, i+3))
	}

	found := 0
	err := q.SelectOne(
		&found,
		fmt.Sprintf(`
			SELECT
					COUNT(user_id)
				FROM organizations_members
				WHERE organization_id = $1 AND user_id = $2 AND (%s)`, strings.Join(clauses, " OR ")),
		args...)
	if db.IsNoRows(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return (found == len(rights)), nil
}

// ListMembers retrieves all the members from an organization. Optionally a list
// of rights can be passed to filter them.
func (s *organizationStore) ListMembers(ids ttnpb.OrganizationIdentifiers, rights ...ttnpb.Right) (members []ttnpb.OrganizationMember, err error) {
	err = s.transact(func(tx *db.Tx) error {
		orgID, err := s.getOrganizationID(tx, ids)
		if err != nil {
			return err
		}

		members, err = s.listMembers(tx, orgID, rights...)
		if err != nil {
			return err
		}

		for i := range members {
			members[i].OrganizationIdentifiers = ids
		}

		return nil
	})
	return
}

func (s *organizationStore) listMembers(q db.QueryContext, orgID uuid.UUID, rights ...ttnpb.Right) ([]ttnpb.OrganizationMember, error) {
	args := make([]interface{}, 0, 2+len(rights))
	args = append(args, orgID)

	var query string
	if len(rights) == 0 {
		query = `
			SELECT
				user_id,
				"right"
			FROM organizations_members
			WHERE organization_id = $1`
	} else {
		args = append(args, len(rights))
		clauses := make([]string, 0, len(rights))

		for i, right := range rights {
			clauses = append(clauses, fmt.Sprintf(`"right" = $%d`, i+3))
			args = append(args, right)
		}

		query = `
			SELECT
				user_id,
				"right"
			FROM organizations_members
			WHERE organization_id = $1 AND user_id IN (
				SELECT
					user_id
				FROM (
					SELECT
						user_id,
						COUNT(user_id) AS count
					FROM organizations_members
					WHERE organization_id = $1 AND (%s)
					GROUP BY user_id
				)
				WHERE count = $2
			)`

		query = fmt.Sprintf(query, strings.Join(clauses, " OR "))
	}

	var rows []*struct {
		UserID uuid.UUID
		Right  ttnpb.Right
	}
	err := q.Select(&rows, query, args...)
	if err != nil {
		return nil, err
	}
	if rows == nil || len(rows) == 0 {
		return make([]ttnpb.OrganizationMember, 0), nil
	}

	// map the rows by User ID
	byUser := make(map[string]*ttnpb.OrganizationMember)
	for _, row := range rows {
		key := row.UserID.String()
		if _, ok := byUser[key]; !ok {
			identifier, err := s.store().Users.getUserIdentifiersFromID(q, row.UserID)
			if err != nil {
				return nil, err
			}

			byUser[key] = &ttnpb.OrganizationMember{
				UserIdentifiers: identifier,
				Rights:          []ttnpb.Right{},
			}
		}

		byUser[key].Rights = append(byUser[key].Rights, row.Right)
	}

	members := make([]ttnpb.OrganizationMember, 0, len(byUser))
	for _, member := range byUser {
		members = append(members, *member)
	}

	return members, nil
}

// SetMember inserts or updates a member within an organization.
// If the list of rights is empty the member will be unset.
func (s *organizationStore) SetMember(member ttnpb.OrganizationMember) (err error) {
	err = s.transact(func(tx *db.Tx) error {
		orgID, err := s.getOrganizationID(tx, member.OrganizationIdentifiers)
		if err != nil {
			return err
		}

		userID, err := s.store().Users.getUserID(tx, member.UserIdentifiers)
		if err != nil {
			return err
		}

		err = s.removeMember(tx, orgID, userID)
		if err != nil {
			return err
		}

		if len(member.Rights) == 0 {
			return nil
		}

		return s.setMember(tx, orgID, userID, member.Rights)
	})
	return
}

func (s *organizationStore) removeMember(q db.QueryContext, orgID, userID uuid.UUID) error {
	_, err := q.Exec(
		`DELETE FROM
			organizations_members
			WHERE organization_id = $1 AND user_id = $2`,
		orgID,
		userID)
	return err
}

func (s *organizationStore) setMember(q db.QueryContext, orgID, userID uuid.UUID, rights []ttnpb.Right) error {
	args := make([]interface{}, 0, len(rights)+2)
	args = append(args, orgID, userID)

	values := make([]string, 0, len(rights))
	for i, right := range rights {
		args = append(args, right)
		values = append(values, fmt.Sprintf("($1, $2, $%d)", i+3))
	}

	query := fmt.Sprintf(
		`INSERT
			INTO organizations_members (organization_id, user_id, "right")
			VALUES %s
			ON CONFLICT (organization_id, user_id, "right")
			DO NOTHING`,
		strings.Join(values, ", "))

	_, err := q.Exec(query, args...)
	return err
}

// ListMemberRights returns the rights a given user has for an entity.
func (s *organizationStore) ListMemberRights(ids ttnpb.OrganizationIdentifiers, target ttnpb.UserIdentifiers) (rights []ttnpb.Right, err error) {
	err = s.transact(func(tx *db.Tx) error {
		orgID, err := s.getOrganizationID(tx, ids)
		if err != nil {
			return err
		}

		userID, err := s.store().Users.getUserID(tx, target)
		if err != nil {
			return err
		}

		return s.listMemberRights(tx, orgID, userID, &rights)
	})
	return
}

func (s *organizationStore) listMemberRights(q db.QueryContext, orgID, userID uuid.UUID, result *[]ttnpb.Right) (err error) {
	err = q.Select(
		result,
		`SELECT
				"right"
			FROM organizations_members
			WHERE organization_id = $1 AND user_id = $2`,
		orgID,
		userID)
	return
}
