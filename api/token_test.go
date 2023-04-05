package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mock_db "github.com/asdsec/thenut/db/mock"
	db "github.com/asdsec/thenut/db/sqlc"
	"github.com/asdsec/thenut/token"
	mock_token "github.com/asdsec/thenut/token/mock"
	"github.com/asdsec/thenut/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRenewAccessTokenAPI(t *testing.T) {
	session := randomSession()
	refreshTokenPayload := token.TokenPayload{
		ID:       session.ID,
		Username: session.Username,
	}
	accessToken := utils.RandomString(32)
	rsp := tokenResponse{
		ExpiresIn:    int(testConfig.AccessTokenDuration.Seconds()),
		RefreshToken: session.RefreshToken,
		AccessToken:  accessToken,
		Username:     session.Username,
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"refresh_token": session.RefreshToken,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				tokenMaker.EXPECT().
					VerifyToken(gomock.Eq(session.RefreshToken)).
					Times(1).
					Return(&refreshTokenPayload, nil)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Eq(session.Username), gomock.Eq(testConfig.AccessTokenDuration)).
					Times(1).
					Return(accessToken, &token.TokenPayload{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBdyMatcherTokenResponse(t, recorder.Body, rsp)
			},
		},
		{
			name: "NullRefreshToken",
			body: gin.H{
				"refresh_token": "",
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				tokenMaker.EXPECT().
					VerifyToken(gomock.Any()).
					Times(0)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidToken",
			body: gin.H{
				"refresh_token": session.RefreshToken,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				tokenMaker.EXPECT().
					VerifyToken(gomock.Eq(session.RefreshToken)).
					Times(1).
					Return(&token.TokenPayload{}, token.ErrInvalidToken)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			body: gin.H{
				"refresh_token": session.RefreshToken,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				tokenMaker.EXPECT().
					VerifyToken(gomock.Eq(session.RefreshToken)).
					Times(1).
					Return(&token.TokenPayload{}, token.ErrExpiredToken)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"refresh_token": session.RefreshToken,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				tokenMaker.EXPECT().
					VerifyToken(gomock.Eq(session.RefreshToken)).
					Times(1).
					Return(&refreshTokenPayload, nil)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(db.Session{}, sql.ErrNoRows)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"refresh_token": session.RefreshToken,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				tokenMaker.EXPECT().
					VerifyToken(gomock.Eq(session.RefreshToken)).
					Times(1).
					Return(&refreshTokenPayload, nil)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(db.Session{}, sql.ErrConnDone)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BlockedSession",
			body: gin.H{
				"refresh_token": session.RefreshToken,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				blockedSession := db.Session{
					ID:           session.ID,
					Username:     session.Username,
					RefreshToken: session.RefreshToken,
					UserAgent:    session.UserAgent,
					ClientIp:     session.ClientIp,
					IsBlocked:    true,
					ExpiresAt:    session.ExpiresAt,
					CreatedAt:    session.CreatedAt,
				}

				tokenMaker.EXPECT().
					VerifyToken(gomock.Eq(session.RefreshToken)).
					Times(1).
					Return(&refreshTokenPayload, nil)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(blockedSession, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ImproperUsername",
			body: gin.H{
				"refresh_token": session.RefreshToken,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				differentRefreshTokenPayload := token.TokenPayload{
					ID:       session.ID,
					Username: "improper_username",
				}

				tokenMaker.EXPECT().
					VerifyToken(gomock.Eq(session.RefreshToken)).
					Times(1).
					Return(&differentRefreshTokenPayload, nil)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "MismatchSessionToken",
			body: gin.H{
				"refresh_token": "mismatch_token",
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				tokenMaker.EXPECT().
					VerifyToken(gomock.Any()).
					Times(1).
					Return(&refreshTokenPayload, nil)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredSession",
			body: gin.H{
				"refresh_token": session.RefreshToken,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				expiredSession := db.Session{
					ID:           session.ID,
					Username:     session.Username,
					RefreshToken: session.RefreshToken,
					UserAgent:    session.UserAgent,
					ClientIp:     session.ClientIp,
					IsBlocked:    session.IsBlocked,
					ExpiresAt:    time.Now().Add(-time.Minute),
					CreatedAt:    session.CreatedAt,
				}

				tokenMaker.EXPECT().
					VerifyToken(gomock.Any()).
					Times(1).
					Return(&refreshTokenPayload, nil)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(expiredSession, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalErrorCreateToken",
			body: gin.H{
				"refresh_token": session.RefreshToken,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				tokenMaker.EXPECT().
					VerifyToken(gomock.Eq(session.RefreshToken)).
					Times(1).
					Return(&refreshTokenPayload, nil)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Eq(session.Username), gomock.Eq(testConfig.AccessTokenDuration)).
					Times(1).
					Return("", &token.TokenPayload{}, errors.New("internal error"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockStore(ctrl)
			tokenMaker := mock_token.NewMockTokenMaker(ctrl)
			tc.buildStubs(store, tokenMaker)

			server := newTestServer(t, store, tokenMaker)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			require.NotEmpty(t, data)

			url := "/tokens/renew"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomSession() db.Session {
	return db.Session{
		ID:           uuid.UUID{},
		Username:     utils.RandomOwner(),
		RefreshToken: utils.RandomString(32),
		UserAgent:    utils.RandomString(6),
		ClientIp:     utils.RandomString(6),
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(time.Minute),
		CreatedAt:    utils.RandomBirthDate(),
	}
}

func requireBdyMatcherTokenResponse(t *testing.T, body *bytes.Buffer, expectedTokenResponse tokenResponse) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var gotTokenResponse tokenResponse
	err = json.Unmarshal(data, &gotTokenResponse)
	require.NoError(t, err)
	require.Equal(t, expectedTokenResponse.AccessToken, gotTokenResponse.AccessToken)
	require.Equal(t, expectedTokenResponse.ExpiresIn, gotTokenResponse.ExpiresIn)
	require.Equal(t, expectedTokenResponse.RefreshToken, gotTokenResponse.RefreshToken)
	require.Equal(t, expectedTokenResponse.Username, gotTokenResponse.Username)
}
