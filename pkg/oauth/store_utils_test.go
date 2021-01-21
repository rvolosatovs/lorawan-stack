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

package oauth_test

import (
	"context"

	"github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type mockStoreContents struct {
	calls []string
	req   struct {
		ctx               context.Context
		fieldMask         *types.FieldMask
		session           *ttnpb.UserSession
		sessionID         string
		userIDs           *ttnpb.UserIdentifiers
		clientIDs         *ttnpb.ClientIdentifiers
		authorization     *ttnpb.OAuthClientAuthorization
		authorizationCode *ttnpb.OAuthAuthorizationCode
		code              string
		token             *ttnpb.OAuthAccessToken
		previousID        string
		tokenID           string
	}
	res struct {
		session           *ttnpb.UserSession
		user              *ttnpb.User
		client            *ttnpb.Client
		authorization     *ttnpb.OAuthClientAuthorization
		authorizationCode *ttnpb.OAuthAuthorizationCode
		accessToken       *ttnpb.OAuthAccessToken
	}
	err struct {
		getUser                 error
		getSession              error
		deleteSession           error
		getClient               error
		getAuthorization        error
		authorize               error
		createAuthorizationCode error
		getAuthorizationCode    error
		deleteAuthorizationCode error
		createAccessToken       error
		getAccessToken          error
		deleteAccessToken       error
	}
}

type mockStore struct {
	store.UserStore
	store.UserSessionStore
	store.ClientStore
	store.OAuthStore

	mockStoreContents
}

func (s *mockStore) reset() {
	s.mockStoreContents = mockStoreContents{}
}

var (
	mockErrUnauthenticated = grpc.Errorf(codes.Unauthenticated, "Unauthenticated")
	mockErrNotFound        = grpc.Errorf(codes.NotFound, "NotFound")
)

func (s *mockStore) GetUser(ctx context.Context, id *ttnpb.UserIdentifiers, fieldMask *types.FieldMask) (*ttnpb.User, error) {
	s.req.ctx, s.req.userIDs, s.req.fieldMask = ctx, id, fieldMask
	s.calls = append(s.calls, "GetUser")
	return s.res.user, s.err.getUser
}

func (s *mockStore) GetSession(ctx context.Context, userIDs *ttnpb.UserIdentifiers, sessionID string) (*ttnpb.UserSession, error) {
	s.req.ctx, s.req.userIDs, s.req.sessionID = ctx, userIDs, sessionID
	s.calls = append(s.calls, "GetSession")
	return s.res.session, s.err.getSession
}

func (s *mockStore) DeleteSession(ctx context.Context, userIDs *ttnpb.UserIdentifiers, sessionID string) error {
	s.req.ctx, s.req.userIDs, s.req.sessionID = ctx, userIDs, sessionID
	s.calls = append(s.calls, "DeleteSession")
	return s.err.deleteSession
}

func (s *mockStore) GetClient(ctx context.Context, id *ttnpb.ClientIdentifiers, fieldMask *types.FieldMask) (*ttnpb.Client, error) {
	s.req.ctx, s.req.clientIDs, s.req.fieldMask = ctx, id, fieldMask
	s.calls = append(s.calls, "GetClient")
	return s.res.client, s.err.getClient
}

func (s *mockStore) GetAuthorization(ctx context.Context, userIDs *ttnpb.UserIdentifiers, clientIDs *ttnpb.ClientIdentifiers) (*ttnpb.OAuthClientAuthorization, error) {
	s.req.ctx, s.req.userIDs, s.req.clientIDs = ctx, userIDs, clientIDs
	s.calls = append(s.calls, "GetAuthorization")
	return s.res.authorization, s.err.getAuthorization
}

func (s *mockStore) Authorize(ctx context.Context, req *ttnpb.OAuthClientAuthorization) (authorization *ttnpb.OAuthClientAuthorization, err error) {
	s.req.ctx, s.req.authorization = ctx, req
	s.calls = append(s.calls, "Authorize")
	return s.res.authorization, s.err.authorize
}

func (s *mockStore) CreateAuthorizationCode(ctx context.Context, code *ttnpb.OAuthAuthorizationCode) error {
	s.req.ctx, s.req.authorizationCode = ctx, code
	s.calls = append(s.calls, "CreateAuthorizationCode")
	return s.err.createAuthorizationCode
}

func (s *mockStore) GetAuthorizationCode(ctx context.Context, code string) (*ttnpb.OAuthAuthorizationCode, error) {
	s.req.ctx, s.req.code = ctx, code
	s.calls = append(s.calls, "GetAuthorizationCode")
	return s.res.authorizationCode, s.err.getAuthorizationCode
}

func (s *mockStore) DeleteAuthorizationCode(ctx context.Context, code string) error {
	s.req.ctx, s.req.code = ctx, code
	s.calls = append(s.calls, "DeleteAuthorizationCode")
	return s.err.deleteAuthorizationCode
}

func (s *mockStore) CreateAccessToken(ctx context.Context, token *ttnpb.OAuthAccessToken, previousID string) error {
	s.req.ctx, s.req.token, s.req.previousID = ctx, token, previousID
	s.calls = append(s.calls, "CreateAccessToken")
	return s.err.createAccessToken
}

func (s *mockStore) GetAccessToken(ctx context.Context, tokenID string) (*ttnpb.OAuthAccessToken, error) {
	s.req.ctx, s.req.tokenID = ctx, tokenID
	s.calls = append(s.calls, "GetAccessToken")
	return s.res.accessToken, s.err.getAccessToken
}

func (s *mockStore) DeleteAccessToken(ctx context.Context, tokenID string) error {
	s.req.ctx = ctx
	if tokenID != "" {
		s.req.tokenID = tokenID
	}
	s.calls = append(s.calls, "DeleteAccessToken")
	return s.err.deleteAccessToken
}
