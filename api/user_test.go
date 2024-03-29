package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	mock_db "github.com/asdsec/thenut/db/mock"
	db "github.com/asdsec/thenut/db/sqlc"
	"github.com/asdsec/thenut/token"
	mock_token "github.com/asdsec/thenut/token/mock"
	"github.com/asdsec/thenut/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestLoginUserAPI(t *testing.T) {
	user, password := randomUser(t)
	rsp := randomAuthResponse(user)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Eq(user.Username), gomock.Any()).
					Times(1).
					Return(rsp.AccessToken, &token.TokenPayload{}, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Eq(user.Username), gomock.Any()).
					Times(1).
					Return(rsp.RefreshToken, &token.TokenPayload{}, nil)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAuthResponse(t, recorder.Body, rsp)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username": "invld", // invalid username
				"password": password,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerErrorGetUser",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalServerErrorCreateAccessToken",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Eq(user.Username), gomock.Any()).
					Times(1).
					Return("", &token.TokenPayload{}, errors.New("internal error"))

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalServerErrorCreateAccessToken",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Eq(user.Username), gomock.Any()).
					Times(1).
					Return(rsp.AccessToken, &token.TokenPayload{}, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", &token.TokenPayload{}, errors.New("internal error"))

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalServerErrorCreateSession",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Eq(user.Username), gomock.Any()).
					Times(1).
					Return(rsp.AccessToken, &token.TokenPayload{}, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return(rsp.RefreshToken, &token.TokenPayload{}, nil)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
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

			url := "/auth/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

type eqRegisterUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqRegisterUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	isBirthDateEqual := e.arg.BirthDate.Equal(arg.BirthDate)

	e.arg.BirthDate = arg.BirthDate
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg) && isBirthDateEqual
}

func (e eqRegisterUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v, and password %v", e.arg, e.password)
}

func EqRegisterUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqRegisterUserParamsMatcher{arg, password}
}

func TestRegisterUserAPI(t *testing.T) {
	user, password := randomUser(t)
	rsp := randomAuthResponse(user)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"username":     user.Username,
				"password":     password,
				"email":        user.Email,
				"full_name":    user.FullName,
				"phone_number": user.PhoneNumber,
				"gender":       user.Gender,
				"birth_date":   user.BirthDate,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				arg := db.CreateUserParams{
					Username:    user.Username,
					FullName:    user.FullName,
					Email:       user.Email,
					PhoneNumber: user.PhoneNumber,
					Gender:      user.Gender,
					BirthDate:   user.BirthDate,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqRegisterUserParams(arg, password)).
					Times(1).
					Return(user, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Eq(user.Username), gomock.Any()).
					Times(1).
					Return(rsp.AccessToken, &token.TokenPayload{}, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Eq(user.Username), gomock.Any()).
					Times(1).
					Return(rsp.RefreshToken, &token.TokenPayload{}, nil)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAuthResponse(t, recorder.Body, rsp)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":     "invld", // invalid username
				"password":     password,
				"email":        user.Email,
				"full_name":    user.FullName,
				"phone_number": user.PhoneNumber,
				"gender":       user.Gender,
				"birth_date":   user.BirthDate,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPassword",
			body: gin.H{
				"username":     user.Username,
				"password":     "invld", // invalid password
				"email":        user.Email,
				"full_name":    user.FullName,
				"phone_number": user.PhoneNumber,
				"gender":       user.Gender,
				"birth_date":   user.BirthDate,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":     user.Username,
				"password":     password,
				"email":        user.Email,
				"full_name":    user.FullName,
				"phone_number": user.PhoneNumber,
				"gender":       user.Gender,
				"birth_date":   user.BirthDate,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"username":     user.Username,
				"password":     password,
				"email":        user.Email,
				"full_name":    user.FullName,
				"phone_number": user.PhoneNumber,
				"gender":       user.Gender,
				"birth_date":   user.BirthDate,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalServerErrorCreateAccessToken",
			body: gin.H{
				"username":     user.Username,
				"password":     password,
				"email":        user.Email,
				"full_name":    user.FullName,
				"phone_number": user.PhoneNumber,
				"gender":       user.Gender,
				"birth_date":   user.BirthDate,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", &token.TokenPayload{}, errors.New("internal error"))

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalServerErrorCreateRefreshToken",
			body: gin.H{
				"username":     user.Username,
				"password":     password,
				"email":        user.Email,
				"full_name":    user.FullName,
				"phone_number": user.PhoneNumber,
				"gender":       user.Gender,
				"birth_date":   user.BirthDate,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return(rsp.AccessToken, &token.TokenPayload{}, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", &token.TokenPayload{}, errors.New("internal error"))

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalServerErrorCreateSession",
			body: gin.H{
				"username":     user.Username,
				"password":     password,
				"email":        user.Email,
				"full_name":    user.FullName,
				"phone_number": user.PhoneNumber,
				"gender":       user.Gender,
				"birth_date":   user.BirthDate,
			},
			buildStubs: func(store *mock_db.MockStore, tokenMaker *mock_token.MockTokenMaker) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Eq(user.Username), gomock.Any()).
					Times(1).
					Return(rsp.AccessToken, &token.TokenPayload{}, nil)

				tokenMaker.EXPECT().
					CreateToken(gomock.Eq(user.Username), gomock.Any()).
					Times(1).
					Return(rsp.RefreshToken, &token.TokenPayload{}, nil)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
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

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			require.NotEmpty(t, data)

			url := "/auth/register"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetUserAPI(t *testing.T) {
	user, _ := randomUser(t)

	testCases := []struct {
		name          string
		uri           string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			uri:  user.Username,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InvalidUsername",
			uri:  "invld", // invalid username
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
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
			uri:  user.Username,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			uri:  user.Username,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "UnauthorizedUser",
			uri:  user.Username,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "invalid", time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			uri:  user.Username,
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

			tokenMaker := newTestTokenMaker(t)
			server := newTestServer(t, store, tokenMaker)
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
	user, _ := randomUser(t)

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
				"username": user.Username,
				"email":    user.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				arg := db.UpdateEmailParams{
					Username: user.Username,
					Email:    user.Email,
				}

				store.EXPECT().
					UpdateEmail(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username": "invld", // invalid username
				"email":    user.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
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
				"username": user.Username,
				"email":    user.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				arg := db.UpdateEmailParams{
					Username: user.Username,
					Email:    user.Email,
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
				"username": user.Username,
				"email":    user.Email,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				arg := db.UpdateEmailParams{
					Username: user.Username,
					Email:    user.Email,
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
				"username": user.Username,
				"email":    user.Email,
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
				"username": user.Username,
				"email":    user.Email,
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

			tokenMaker := newTestTokenMaker(t)
			server := newTestServer(t, store, tokenMaker)
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

type eqUpdateUserParamsMatcher struct {
	arg db.UpdateUserParams
}

func (e eqUpdateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateUserParams)
	if !ok {
		return false
	}

	isBirthDateEqual := e.arg.BirthDate.Time.Equal(arg.BirthDate.Time)

	e.arg.BirthDate = arg.BirthDate
	return reflect.DeepEqual(e.arg, arg) && isBirthDateEqual
}

func (e eqUpdateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v, and time obj", e.arg)
}

func EqUpdateUserParams(arg db.UpdateUserParams) gomock.Matcher {
	return eqUpdateUserParamsMatcher{arg}
}

func TestUpdateUserAPI(t *testing.T) {
	user, _ := randomUser(t)
	fullName := "updated_full_name"
	phoneNumber := "updated_phone_number"
	gender := "updated_gender"
	birthDate := time.Now()
	imageUrl := "updated_image_url"

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		buildStubs    func(store *mock_db.MockStore) db.User
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User)
	}{
		{
			name: "Ok",
			body: gin.H{
				"full_name":    fullName,
				"phone_number": phoneNumber,
				"gender":       gender,
				"birth_date":   birthDate,
				"image_url":    imageUrl,
				"username":     user.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				expectedUser := db.User{
					FullName:          fullName,
					PhoneNumber:       phoneNumber,
					Gender:            gender,
					ImageUrl:          imageUrl,
					BirthDate:         birthDate,
					Username:          user.Username,
					HashedPassword:    user.HashedPassword,
					Email:             user.Email,
					Disabled:          user.Disabled,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
				}

				arg := db.UpdateUserParams{
					FullName: sql.NullString{
						String: expectedUser.FullName,
						Valid:  true,
					},
					PhoneNumber: sql.NullString{
						String: expectedUser.PhoneNumber,
						Valid:  true,
					},
					Gender: sql.NullString{
						String: expectedUser.Gender,
						Valid:  true,
					},
					BirthDate: sql.NullTime{
						Time:  expectedUser.BirthDate,
						Valid: true,
					},
					ImageUrl: sql.NullString{
						String: expectedUser.ImageUrl,
						Valid:  true,
					},
					Username: expectedUser.Username,
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), EqUpdateUserParams(arg)).
					Times(1).
					Return(expectedUser, nil)

				return expectedUser
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, expectedUser)
			},
		},
		{
			name: "Ok With Single Field",
			body: gin.H{
				"full_name": fullName,
				"username":  user.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				expectedUser := db.User{
					FullName:          fullName,
					Username:          user.Username,
					HashedPassword:    user.HashedPassword,
					Email:             user.Email,
					PhoneNumber:       user.PhoneNumber,
					ImageUrl:          user.ImageUrl,
					Gender:            user.Gender,
					Disabled:          user.Disabled,
					BirthDate:         user.BirthDate,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
				}

				arg := db.UpdateUserParams{
					FullName: sql.NullString{
						String: expectedUser.FullName,
						Valid:  true,
					},
					Username: expectedUser.Username,
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), EqUpdateUserParams(arg)).
					Times(1).
					Return(expectedUser, nil)

				return expectedUser
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, expectedUser)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"full_name":    fullName,
				"phone_number": phoneNumber,
				"gender":       gender,
				"birth_date":   birthDate,
				"image_url":    imageUrl,
				"username":     "invld", // invalid username
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(0)
				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidBirthDate",
			body: gin.H{
				"birth_date": "invalid_birth_date",
				"username":   user.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(0)
				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UnAuthorizedUser",
			body: gin.H{
				"full_name":    fullName,
				"phone_number": phoneNumber,
				"gender":       gender,
				"birth_date":   birthDate,
				"image_url":    imageUrl,
				"username":     user.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "other_username", time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(0)
				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthentication",
			body: gin.H{
				"full_name":    fullName,
				"phone_number": phoneNumber,
				"gender":       gender,
				"birth_date":   birthDate,
				"image_url":    imageUrl,
				"username":     user.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				// No Authentication
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(0)
				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"full_name":    fullName,
				"phone_number": phoneNumber,
				"gender":       gender,
				"birth_date":   birthDate,
				"image_url":    imageUrl,
				"username":     user.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				expectedUser := db.User{
					FullName:          fullName,
					PhoneNumber:       phoneNumber,
					Gender:            gender,
					ImageUrl:          imageUrl,
					BirthDate:         birthDate,
					Username:          user.Username,
					HashedPassword:    user.HashedPassword,
					Email:             user.Email,
					Disabled:          user.Disabled,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
				}

				arg := db.UpdateUserParams{
					FullName: sql.NullString{
						String: expectedUser.FullName,
						Valid:  true,
					},
					PhoneNumber: sql.NullString{
						String: expectedUser.PhoneNumber,
						Valid:  true,
					},
					Gender: sql.NullString{
						String: expectedUser.Gender,
						Valid:  true,
					},
					BirthDate: sql.NullTime{
						Time:  expectedUser.BirthDate,
						Valid: true,
					},
					ImageUrl: sql.NullString{
						String: expectedUser.ImageUrl,
						Valid:  true,
					},
					Username: expectedUser.Username,
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), EqUpdateUserParams(arg)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"full_name":    fullName,
				"phone_number": phoneNumber,
				"gender":       gender,
				"birth_date":   birthDate,
				"image_url":    imageUrl,
				"username":     user.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				expectedUser := db.User{
					FullName:          fullName,
					PhoneNumber:       phoneNumber,
					Gender:            gender,
					ImageUrl:          imageUrl,
					BirthDate:         birthDate,
					Username:          user.Username,
					HashedPassword:    user.HashedPassword,
					Email:             user.Email,
					Disabled:          user.Disabled,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
				}

				arg := db.UpdateUserParams{
					FullName: sql.NullString{
						String: expectedUser.FullName,
						Valid:  true,
					},
					PhoneNumber: sql.NullString{
						String: expectedUser.PhoneNumber,
						Valid:  true,
					},
					Gender: sql.NullString{
						String: expectedUser.Gender,
						Valid:  true,
					},
					BirthDate: sql.NullTime{
						Time:  expectedUser.BirthDate,
						Valid: true,
					},
					ImageUrl: sql.NullString{
						String: expectedUser.ImageUrl,
						Valid:  true,
					},
					Username: expectedUser.Username,
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), EqUpdateUserParams(arg)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
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
			expectedUser := tc.buildStubs(store)

			tokenMaker := newTestTokenMaker(t)
			server := newTestServer(t, store, tokenMaker)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			require.NotEmpty(t, data)

			url := "/users"
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder, expectedUser)
		})
	}
}

type eqUpdatePasswordParamsMatcher struct {
	arg      db.UpdatePasswordParams
	password string
}

func (e eqUpdatePasswordParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdatePasswordParams)
	if !ok {
		return false
	}

	err := utils.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	dt := e.arg.PasswordChangedAt.Sub(arg.PasswordChangedAt)
	if dt < -time.Minute || dt > time.Minute {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	e.arg.PasswordChangedAt = arg.PasswordChangedAt
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqUpdatePasswordParamsMatcher) String() string {
	return fmt.Sprintf("matches arg: %v, and password: %v", e.arg, e.password)
}

func EqUpdatePasswordParams(arg db.UpdatePasswordParams, password string) gomock.Matcher {
	return eqUpdatePasswordParamsMatcher{arg, password}
}

func TestUpdatePasswordAPI(t *testing.T) {
	user, oldPassword := randomUser(t)
	updatedPassword := "updated_password"

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		buildStubs    func(store *mock_db.MockStore) db.User
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User)
	}{
		{
			name: "Ok",
			body: gin.H{
				"username":     user.Username,
				"old_password": oldPassword,
				"new_password": updatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				expectedUser := db.User{
					Username:          user.Username,
					FullName:          user.FullName,
					HashedPassword:    user.HashedPassword,
					Email:             user.Email,
					PhoneNumber:       user.PhoneNumber,
					ImageUrl:          user.ImageUrl,
					Gender:            user.Gender,
					Disabled:          user.Disabled,
					BirthDate:         user.BirthDate,
					PasswordChangedAt: time.Now(),
					CreatedAt:         user.CreatedAt,
				}

				arg := db.UpdatePasswordParams{
					Username:          expectedUser.Username,
					PasswordChangedAt: expectedUser.PasswordChangedAt,
				}

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					UpdatePassword(gomock.Any(), EqUpdatePasswordParams(arg, updatedPassword)).
					Times(1).
					Return(expectedUser, nil)

				return expectedUser
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, expectedUser)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":     "invld", // invalid username
				"old_password": oldPassword,
				"new_password": updatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					UpdatePassword(gomock.Any(), gomock.Any()).
					Times(0)

				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPassword",
			body: gin.H{
				"username":     user.Username,
				"old_password": "invld", // invalid password
				"new_password": updatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					UpdatePassword(gomock.Any(), gomock.Any()).
					Times(0)

				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UnAuthorizedUser",
			body: gin.H{
				"username":     user.Username,
				"old_password": oldPassword,
				"new_password": updatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "other_username", time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					UpdatePassword(gomock.Any(), gomock.Any()).
					Times(0)

				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthentication",
			body: gin.H{
				"username":     user.Username,
				"old_password": oldPassword,
				"new_password": updatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				// no authentication
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					UpdatePassword(gomock.Any(), gomock.Any()).
					Times(0)

				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NotFoundInGetUser",
			body: gin.H{
				"username":     user.Username,
				"old_password": oldPassword,
				"new_password": updatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)

				store.EXPECT().
					UpdatePassword(gomock.Any(), gomock.Any()).
					Times(0)

				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerErrorGetUser",
			body: gin.H{
				"username":     user.Username,
				"old_password": oldPassword,
				"new_password": updatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)

				store.EXPECT().
					UpdatePassword(gomock.Any(), gomock.Any()).
					Times(0)

				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "WrongOldPassword",
			body: gin.H{
				"username":     user.Username,
				"old_password": "wrong_old_password",
				"new_password": updatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					UpdatePassword(gomock.Any(), gomock.Any()).
					Times(0)

				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalErrorHashPassword",
			body: gin.H{
				"username":     user.Username,
				"old_password": oldPassword,
				"new_password": strings.Repeat("x", 73), // too long password
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					UpdatePassword(gomock.Any(), gomock.Any()).
					Times(0)

				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFoundInUpdateUser",
			body: gin.H{
				"username":     user.Username,
				"old_password": oldPassword,
				"new_password": updatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				expectedUser := db.User{
					Username:          user.Username,
					FullName:          user.FullName,
					HashedPassword:    user.HashedPassword,
					Email:             user.Email,
					PhoneNumber:       user.PhoneNumber,
					ImageUrl:          user.ImageUrl,
					Gender:            user.Gender,
					Disabled:          user.Disabled,
					BirthDate:         user.BirthDate,
					PasswordChangedAt: time.Now(),
					CreatedAt:         user.CreatedAt,
				}

				arg := db.UpdatePasswordParams{
					Username:          expectedUser.Username,
					PasswordChangedAt: expectedUser.PasswordChangedAt,
				}

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					UpdatePassword(gomock.Any(), EqUpdatePasswordParams(arg, updatedPassword)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)

				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerErrorUpdateUser",
			body: gin.H{
				"username":     user.Username,
				"old_password": oldPassword,
				"new_password": updatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				expectedUser := db.User{
					Username:          user.Username,
					FullName:          user.FullName,
					HashedPassword:    user.HashedPassword,
					Email:             user.Email,
					PhoneNumber:       user.PhoneNumber,
					ImageUrl:          user.ImageUrl,
					Gender:            user.Gender,
					Disabled:          user.Disabled,
					BirthDate:         user.BirthDate,
					PasswordChangedAt: time.Now(),
					CreatedAt:         user.CreatedAt,
				}

				arg := db.UpdatePasswordParams{
					Username:          expectedUser.Username,
					PasswordChangedAt: expectedUser.PasswordChangedAt,
				}

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					UpdatePassword(gomock.Any(), EqUpdatePasswordParams(arg, updatedPassword)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)

				return db.User{}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, expectedUser db.User) {
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
			expectedUser := tc.buildStubs(store)

			tokenMaker := newTestTokenMaker(t)
			server := newTestServer(t, store, tokenMaker)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			require.NotEmpty(t, data)

			url := "/users/password"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder, expectedUser)
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
	require.WithinDuration(t, user.CreatedAt, gotUser.CreatedAt, time.Second)
	require.WithinDuration(t, user.BirthDate, gotUser.BirthDate, time.Second)
	require.Equal(t, user.Disabled, gotUser.Disabled)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Gender, gotUser.Gender)
	require.Equal(t, user.ImageUrl, gotUser.ImageUrl)
	require.Equal(t, user.PhoneNumber, gotUser.PhoneNumber)
	require.Equal(t, user.Username, gotUser.Username)
}

func randomAuthResponse(user db.User) authResponse {
	return authResponse{
		AccessToken:  utils.RandomString(32),
		Email:        user.Email,
		FullName:     user.FullName,
		Username:     user.Username,
		RefreshToken: utils.RandomString(32),
		ExpiresIn:    int(testConfig.AccessTokenDuration.Seconds()),
	}
}

func requireBodyMatchAuthResponse(t *testing.T, body *bytes.Buffer, expectedAuthResponse authResponse) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var gotAuthResponse authResponse
	err = json.Unmarshal(data, &gotAuthResponse)
	require.NoError(t, err)
	require.Equal(t, expectedAuthResponse.AccessToken, gotAuthResponse.AccessToken)
	require.Equal(t, expectedAuthResponse.Email, gotAuthResponse.Email)
	require.Equal(t, expectedAuthResponse.ExpiresIn, gotAuthResponse.ExpiresIn)
	require.Equal(t, expectedAuthResponse.FullName, gotAuthResponse.FullName)
	require.Equal(t, expectedAuthResponse.RefreshToken, gotAuthResponse.RefreshToken)
	require.Equal(t, expectedAuthResponse.Username, gotAuthResponse.Username)
}
