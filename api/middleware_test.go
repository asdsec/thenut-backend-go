package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/asdsec/thenut/token"
	mock_token "github.com/asdsec/thenut/token/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.TokenMaker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, payload, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	const blank = " "

	testCases := []struct {
		name          string
		headerKey     string
		headerValue   string
		buildStubs    func(tokenMaker *mock_token.MockTokenMaker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "Ok",
			headerKey:   authorizationHeaderKey,
			headerValue: authorizationTypeBearer + blank + "valid_token",
			buildStubs: func(tokenMaker *mock_token.MockTokenMaker) {
				tokenMaker.EXPECT().
					VerifyToken(gomock.Eq("valid_token")).
					Times(1).
					Return(&token.TokenPayload{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:        "NoAuthorizationHeader",
			headerKey:   authorizationHeaderKey,
			headerValue: "",
			buildStubs: func(tokenMaker *mock_token.MockTokenMaker) {
				tokenMaker.EXPECT().
					VerifyToken(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:        "InvalidAuthorizationHeader",
			headerKey:   authorizationHeaderKey,
			headerValue: "valid_token",
			buildStubs: func(tokenMaker *mock_token.MockTokenMaker) {
				tokenMaker.EXPECT().
					VerifyToken(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:        "UnsupportedAuthorizationType",
			headerKey:   authorizationHeaderKey,
			headerValue: "unsupported_auth_type" + blank + "valid_token",
			buildStubs: func(tokenMaker *mock_token.MockTokenMaker) {
				tokenMaker.EXPECT().
					VerifyToken(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:        "InvalidToken",
			headerKey:   authorizationHeaderKey,
			headerValue: authorizationTypeBearer + blank + "invalid_token",
			buildStubs: func(tokenMaker *mock_token.MockTokenMaker) {
				tokenMaker.EXPECT().
					VerifyToken(gomock.Eq("invalid_token")).
					Times(1).
					Return(&token.TokenPayload{}, token.ErrInvalidToken)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:        "ExpiredToken",
			headerKey:   authorizationHeaderKey,
			headerValue: authorizationTypeBearer + blank + "expired_token",
			buildStubs: func(tokenMaker *mock_token.MockTokenMaker) {
				tokenMaker.EXPECT().
					VerifyToken(gomock.Eq("expired_token")).
					Times(1).
					Return(&token.TokenPayload{}, token.ErrExpiredToken)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tokenMaker := mock_token.NewMockTokenMaker(ctrl)
			tc.buildStubs(tokenMaker)

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)

			handler := authMiddleware(tokenMaker)

			var err error
			ctx.Request, err = http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)
			ctx.Request.Header.Set(tc.headerKey, tc.headerValue)

			handler(ctx)
			tc.checkResponse(t, recorder)
		})
	}

}
