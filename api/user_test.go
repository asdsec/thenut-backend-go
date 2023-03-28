package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mock_db "github.com/asdsec/thenut/db/mock"
	db "github.com/asdsec/thenut/db/sqlc"
	"github.com/asdsec/thenut/token"
	"github.com/asdsec/thenut/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetUserAPI(t *testing.T) {
	tUser, _ := randomUser(t)

	testCases := []struct {
		name          string
		uri           string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			uri:  tUser.Username,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(tUser.Username)).
					Times(1).
					Return(tUser, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, tUser)
			},
		},
		{
			name: "InvalidUsername",
			uri:  "invld", // invalid username
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			uri:  tUser.Username,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(tUser.Username)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			uri:  tUser.Username,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(tUser.Username)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "UnauthorizedUser",
			uri:  tUser.Username,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "invalid", time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(tUser.Username)).
					Times(1).
					Return(tUser, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			uri:  tUser.Username,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				// no authentication
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/users/" + tc.uri
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateEmailAPI(t *testing.T) {
	tUser, _ := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"username": tUser.Username,
				"email":    tUser.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				arg := db.UpdateEmailParams{
					Username: tUser.Username,
					Email:    tUser.Email,
				}

				store.EXPECT().
					UpdateEmail(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(tUser, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, tUser)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username": "invld", // invalid username
				"email":    tUser.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					UpdateEmail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"username": tUser.Username,
				"email":    tUser.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				arg := db.UpdateEmailParams{
					Username: tUser.Username,
					Email:    tUser.Email,
				}

				store.EXPECT().
					UpdateEmail(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username": tUser.Username,
				"email":    tUser.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				arg := db.UpdateEmailParams{
					Username: tUser.Username,
					Email:    tUser.Email,
				}

				store.EXPECT().
					UpdateEmail(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "UnauthenticatedUser",
			body: gin.H{
				"username": tUser.Username,
				"email":    tUser.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "invalid_username", time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					UpdateEmail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthentication",
			body: gin.H{
				"username": tUser.Username,
				"email":    tUser.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				// no authentication
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					UpdateEmail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			require.NotEmpty(t, data)

			url := "/users/email"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	user = db.User{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
		PhoneNumber:    utils.RandomPhoneNumber(),
		ImageUrl:       utils.RandomImageUrl(),
		Gender:         utils.RandomGender(),
		Disabled:       false,
		BirthDate:      utils.RandomBirthDate(),
	}

	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.WithinDuration(t, user.BirthDate, gotUser.BirthDate, time.Second)
	require.WithinDuration(t, user.CreatedAt, gotUser.CreatedAt, time.Second)
	require.Equal(t, user.Disabled, gotUser.Disabled)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Gender, gotUser.Gender)
	require.Equal(t, user.ImageUrl, gotUser.ImageUrl)
	require.Equal(t, user.PhoneNumber, gotUser.PhoneNumber)
	require.Equal(t, user.Username, gotUser.Username)
}
