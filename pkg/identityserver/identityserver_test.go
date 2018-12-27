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

package identityserver

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/pkg/component"
	"go.thethings.network/lorawan-stack/pkg/config"
	"go.thethings.network/lorawan-stack/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/rpcmetadata"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/util/test"
	"google.golang.org/grpc"
)

var (
	setup        sync.Once
	dbConnString string
	population   = store.NewPopulator(12, 42)
)

var (
	userIndex                                         int
	newUser, newUserIdx                               = getTestUser()
	rejectedUser, rejectedUserIdx                     = getTestUser()
	defaultUser, defaultUserIdx                       = getTestUser()
	suspendedUser, suspendedUserIdx                   = getTestUser()
	adminUser, adminUserIdx                           = getTestUser()
	collaboratorUser, collaboratorUserIdx             = getTestUser()
	applicationAccessUser, applicationAccessUserIdx   = getTestUser()
	clientAccessUser, clientAccessUserIdx             = getTestUser()
	gatewayAccessUser, gatewayAccessUserIdx           = getTestUser()
	organizationAccessUser, organizationAccessUserIdx = getTestUser()
	userAccessUser, userAccessUserIdx                 = getTestUser()
	paginationUser, paginationUserIdx                 = getTestUser()
)

var now = time.Now()

func init() {
	newUser.Admin = false
	newUser.PrimaryEmailAddressValidatedAt = nil
	newUser.State = ttnpb.STATE_REQUESTED

	rejectedUser.Admin = false
	rejectedUser.PrimaryEmailAddressValidatedAt = &now
	rejectedUser.State = ttnpb.STATE_REJECTED

	defaultUser.Admin = false
	defaultUser.PrimaryEmailAddressValidatedAt = &now
	defaultUser.State = ttnpb.STATE_APPROVED

	for id, apiKeys := range population.APIKeys {
		if id.GetUserIDs().GetUserID() == defaultUser.GetUserID() {
			population.APIKeys[id] = append(
				apiKeys,
				&ttnpb.APIKey{
					Name:   "key without rights",
					Rights: []ttnpb.Right{ttnpb.RIGHT_SEND_INVITES},
				},
			)
		}
	}

	suspendedUser.Admin = false
	suspendedUser.PrimaryEmailAddressValidatedAt = &now
	suspendedUser.State = ttnpb.STATE_SUSPENDED

	adminUser.Admin = true
	adminUser.PrimaryEmailAddressValidatedAt = &now
	adminUser.State = ttnpb.STATE_APPROVED

	paginationUser.Admin = false
	paginationUser.PrimaryEmailAddressValidatedAt = &now
	paginationUser.State = ttnpb.STATE_APPROVED
}

func getTestUser() (*ttnpb.User, int) {
	defer func() { userIndex++ }()

	return population.Users[userIndex], userIndex
}

func userCreds(idx int, preferredNames ...string) grpc.CallOption {
	for id, apiKeys := range population.APIKeys {
		if id.GetUserIDs().GetUserID() == population.Users[idx].GetUserID() {
			selectedIdx := 0
			if len(preferredNames) == 0 {
				preferredNames = []string{"default key"}
			}
		findPreferred:
			for _, name := range preferredNames {
				for i, apiKey := range apiKeys {
					if apiKey.Name == name {
						selectedIdx = i
						break findPreferred
					}
				}
			}
			return grpc.PerRPCCredentials(rpcmetadata.MD{
				AuthType:      "bearer",
				AuthValue:     apiKeys[selectedIdx].Key,
				AllowInsecure: true,
			})
		}
	}
	return nil
}

func userAPIKeys(userID *ttnpb.UserIdentifiers) ttnpb.APIKeys {
	for id, apiKeys := range population.APIKeys {
		if id.GetUserIDs().GetUserID() == userID.GetUserID() {
			return ttnpb.APIKeys{
				APIKeys: apiKeys,
			}
		}
	}

	return ttnpb.APIKeys{
		APIKeys: []*ttnpb.APIKey{},
	}
}

func applicationAPIKeys(applicationID *ttnpb.ApplicationIdentifiers) ttnpb.APIKeys {
	for id, apiKeys := range population.APIKeys {
		if id.GetApplicationIDs().GetApplicationID() == applicationID.GetApplicationID() {
			return ttnpb.APIKeys{
				APIKeys: apiKeys,
			}
		}
	}

	return ttnpb.APIKeys{
		APIKeys: []*ttnpb.APIKey{},
	}
}

func gatewayAPIKeys(gatewayID *ttnpb.GatewayIdentifiers) ttnpb.APIKeys {
	for id, apiKeys := range population.APIKeys {
		if id.GetGatewayIDs().GetGatewayID() == gatewayID.GetGatewayID() {
			return ttnpb.APIKeys{
				APIKeys: apiKeys,
			}
		}
	}

	return ttnpb.APIKeys{
		APIKeys: []*ttnpb.APIKey{},
	}
}

func organizationAPIKeys(organizationID *ttnpb.OrganizationIdentifiers) ttnpb.APIKeys {
	for id, apiKeys := range population.APIKeys {
		if id.GetOrganizationIDs().GetOrganizationID() == organizationID.GetOrganizationID() {
			return ttnpb.APIKeys{
				APIKeys: apiKeys,
			}
		}
	}

	return ttnpb.APIKeys{
		APIKeys: []*ttnpb.APIKey{},
	}
}

func userApplications(userID *ttnpb.UserIdentifiers) ttnpb.Applications {
	applications := []*ttnpb.Application{}
	for _, app := range population.Applications {
		for id, collaborators := range population.Memberships {
			if app.EntityIdentifiers().IDString() == id.IDString() {
				for _, collaborator := range collaborators {
					if collaborator.EntityIdentifiers().IDString() == userID.GetUserID() {
						applications = append(applications, app)
					}
				}
			}
		}
	}

	return ttnpb.Applications{
		Applications: applications,
	}
}

func userClients(userID *ttnpb.UserIdentifiers) ttnpb.Clients {
	clients := []*ttnpb.Client{}
	for _, client := range population.Clients {
		for id, collaborators := range population.Memberships {
			if client.EntityIdentifiers().IDString() == id.IDString() {
				for _, collaborator := range collaborators {
					if collaborator.EntityIdentifiers().IDString() == userID.GetUserID() {
						clients = append(clients, client)
					}
				}
			}
		}
	}

	return ttnpb.Clients{
		Clients: clients,
	}
}

func userGateways(userID *ttnpb.UserIdentifiers) ttnpb.Gateways {
	gateways := []*ttnpb.Gateway{}
	for _, gateway := range population.Gateways {
		for id, collaborators := range population.Memberships {
			if gateway.EntityIdentifiers().IDString() == id.IDString() {
				for _, collaborator := range collaborators {
					if collaborator.EntityIdentifiers().IDString() == userID.GetUserID() {
						gateways = append(gateways, gateway)
					}
				}
			}
		}
	}

	return ttnpb.Gateways{
		Gateways: gateways,
	}
}

func userOrganizations(userID *ttnpb.UserIdentifiers) ttnpb.Organizations {
	organizations := []*ttnpb.Organization{}
	for _, organization := range population.Organizations {
		for id, collaborators := range population.Memberships {
			if organization.EntityIdentifiers().IDString() == id.IDString() {
				for _, collaborator := range collaborators {
					if collaborator.EntityIdentifiers().IDString() == userID.GetUserID() {
						organizations = append(organizations, organization)
					}
				}
			}
		}
	}

	return ttnpb.Organizations{
		Organizations: organizations,
	}
}

func getIdentityServer(t *testing.T) (*IdentityServer, *grpc.ClientConn) {
	setup.Do(func() {
		dbName := os.Getenv("TEST_DB_NAME")
		if dbName == "" {
			dbName = "is_integration_test"
		}
		dbConnString = fmt.Sprintf("postgresql://root@localhost:26257/%s?sslmode=disable", dbName)
		db, err := gorm.Open("postgres", dbConnString)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		if err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName)).Error; err != nil {
			panic(err)
		}
		store.AutoMigrate(db)
		if err = store.Clear(db); err != nil {
			panic(err)
		}
		if err = population.Populate(test.Context(), db); err != nil {
			panic(err)
		}
	})
	c := component.MustNew(test.GetLogger(t), &component.Config{ServiceBase: config.ServiceBase{
		Base: config.Base{
			Log: config.Log{
				Level: log.DebugLevel,
			},
		},
	}})
	is, err := New(c, &Config{
		DatabaseURI: dbConnString,
	})
	if err != nil {
		panic(err)
	}
	if err = is.Start(); err != nil {
		panic(err)
	}
	return is, is.LoopbackConn()
}

func testWithIdentityServer(t *testing.T, f func(is *IdentityServer, cc *grpc.ClientConn)) {
	f(getIdentityServer(t))
}
