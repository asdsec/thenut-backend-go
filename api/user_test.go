package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
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
	"github.com/asdsec/thenut/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// todo: implement register user api test
// todo: implement login user api test

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
	tUser, _ := randomUser(t)
	tFullName := "updated_full_name"
	tPhoneNumber := "updated_phone_number"
	tGender := "updated_gender"
	tBirthDate := time.Now()
	tImageUrl := "updated_image_url"

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
				"full_name":    tFullName,
				"phone_number": tPhoneNumber,
				"gender":       tGender,
				"birth_date":   tBirthDate,
				"image_url":    tImageUrl,
				"username":     tUser.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				expectedUser := db.User{
					FullName:          tFullName,
					PhoneNumber:       tPhoneNumber,
					Gender:            tGender,
					ImageUrl:          tImageUrl,
					BirthDate:         tBirthDate,
					Username:          tUser.Username,
					HashedPassword:    tUser.HashedPassword,
					Email:             tUser.Email,
					Disabled:          tUser.Disabled,
					PasswordChangedAt: tUser.PasswordChangedAt,
					CreatedAt:         tUser.CreatedAt,
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
				"full_name": tFullName,
				"username":  tUser.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				expectedUser := db.User{
					FullName:          tFullName,
					Username:          tUser.Username,
					HashedPassword:    tUser.HashedPassword,
					Email:             tUser.Email,
					PhoneNumber:       tUser.PhoneNumber,
					ImageUrl:          tUser.ImageUrl,
					Gender:            tUser.Gender,
					Disabled:          tUser.Disabled,
					BirthDate:         tUser.BirthDate,
					PasswordChangedAt: tUser.PasswordChangedAt,
					CreatedAt:         tUser.CreatedAt,
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
				"full_name":    tFullName,
				"phone_number": tPhoneNumber,
				"gender":       tGender,
				"birth_date":   tBirthDate,
				"image_url":    tImageUrl,
				"username":     "invld", // invalid username
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
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
				"username":   tUser.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
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
				"full_name":    tFullName,
				"phone_number": tPhoneNumber,
				"gender":       tGender,
				"birth_date":   tBirthDate,
				"image_url":    tImageUrl,
				"username":     tUser.Username,
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
				"full_name":    tFullName,
				"phone_number": tPhoneNumber,
				"gender":       tGender,
				"birth_date":   tBirthDate,
				"image_url":    tImageUrl,
				"username":     tUser.Username,
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
				"full_name":    tFullName,
				"phone_number": tPhoneNumber,
				"gender":       tGender,
				"birth_date":   tBirthDate,
				"image_url":    tImageUrl,
				"username":     tUser.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				expectedUser := db.User{
					FullName:          tFullName,
					PhoneNumber:       tPhoneNumber,
					Gender:            tGender,
					ImageUrl:          tImageUrl,
					BirthDate:         tBirthDate,
					Username:          tUser.Username,
					HashedPassword:    tUser.HashedPassword,
					Email:             tUser.Email,
					Disabled:          tUser.Disabled,
					PasswordChangedAt: tUser.PasswordChangedAt,
					CreatedAt:         tUser.CreatedAt,
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
				"full_name":    tFullName,
				"phone_number": tPhoneNumber,
				"gender":       tGender,
				"birth_date":   tBirthDate,
				"image_url":    tImageUrl,
				"username":     tUser.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				expectedUser := db.User{
					FullName:          tFullName,
					PhoneNumber:       tPhoneNumber,
					Gender:            tGender,
					ImageUrl:          tImageUrl,
					BirthDate:         tBirthDate,
					Username:          tUser.Username,
					HashedPassword:    tUser.HashedPassword,
					Email:             tUser.Email,
					Disabled:          tUser.Disabled,
					PasswordChangedAt: tUser.PasswordChangedAt,
					CreatedAt:         tUser.CreatedAt,
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

			server := newTestServer(t, store)
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
	tUser, tOldPassword := randomUser(t)
	tUpdatedPassword := "updated_password"

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
				"username":     tUser.Username,
				"old_password": tOldPassword,
				"new_password": tUpdatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				expectedUser := db.User{
					Username:          tUser.Username,
					FullName:          tUser.FullName,
					HashedPassword:    tUser.HashedPassword,
					Email:             tUser.Email,
					PhoneNumber:       tUser.PhoneNumber,
					ImageUrl:          tUser.ImageUrl,
					Gender:            tUser.Gender,
					Disabled:          tUser.Disabled,
					BirthDate:         tUser.BirthDate,
					PasswordChangedAt: time.Now(),
					CreatedAt:         tUser.CreatedAt,
				}

				arg := db.UpdatePasswordParams{
					Username:          expectedUser.Username,
					PasswordChangedAt: expectedUser.PasswordChangedAt,
				}

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(tUser.Username)).
					Times(1).
					Return(tUser, nil)

				store.EXPECT().
					UpdatePassword(gomock.Any(), EqUpdatePasswordParams(arg, tUpdatedPassword)).
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
				"old_password": tOldPassword,
				"new_password": tUpdatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
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
				"username":     tUser.Username,
				"old_password": "invld", // invalid password
				"new_password": tUpdatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
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
				"username":     tUser.Username,
				"old_password": tOldPassword,
				"new_password": tUpdatedPassword,
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
				"username":     tUser.Username,
				"old_password": tOldPassword,
				"new_password": tUpdatedPassword,
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
				"username":     tUser.Username,
				"old_password": tOldPassword,
				"new_password": tUpdatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(tUser.Username)).
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
				"username":     tUser.Username,
				"old_password": tOldPassword,
				"new_password": tUpdatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(tUser.Username)).
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
				"username":     tUser.Username,
				"old_password": "wrong_old_password",
				"new_password": tUpdatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(tUser.Username)).
					Times(1).
					Return(tUser, nil)

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
				"username":     tUser.Username,
				"old_password": tOldPassword,
				"new_password": strings.Repeat("x", 73), // too long password
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(tUser.Username)).
					Times(1).
					Return(tUser, nil)

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
				"username":     tUser.Username,
				"old_password": tOldPassword,
				"new_password": tUpdatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				expectedUser := db.User{
					Username:          tUser.Username,
					FullName:          tUser.FullName,
					HashedPassword:    tUser.HashedPassword,
					Email:             tUser.Email,
					PhoneNumber:       tUser.PhoneNumber,
					ImageUrl:          tUser.ImageUrl,
					Gender:            tUser.Gender,
					Disabled:          tUser.Disabled,
					BirthDate:         tUser.BirthDate,
					PasswordChangedAt: time.Now(),
					CreatedAt:         tUser.CreatedAt,
				}

				arg := db.UpdatePasswordParams{
					Username:          expectedUser.Username,
					PasswordChangedAt: expectedUser.PasswordChangedAt,
				}

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(tUser.Username)).
					Times(1).
					Return(tUser, nil)

				store.EXPECT().
					UpdatePassword(gomock.Any(), EqUpdatePasswordParams(arg, tUpdatedPassword)).
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
				"username":     tUser.Username,
				"old_password": tOldPassword,
				"new_password": tUpdatedPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, tUser.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) db.User {
				expectedUser := db.User{
					Username:          tUser.Username,
					FullName:          tUser.FullName,
					HashedPassword:    tUser.HashedPassword,
					Email:             tUser.Email,
					PhoneNumber:       tUser.PhoneNumber,
					ImageUrl:          tUser.ImageUrl,
					Gender:            tUser.Gender,
					Disabled:          tUser.Disabled,
					BirthDate:         tUser.BirthDate,
					PasswordChangedAt: time.Now(),
					CreatedAt:         tUser.CreatedAt,
				}

				arg := db.UpdatePasswordParams{
					Username:          expectedUser.Username,
					PasswordChangedAt: expectedUser.PasswordChangedAt,
				}

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(tUser.Username)).
					Times(1).
					Return(tUser, nil)

				store.EXPECT().
					UpdatePassword(gomock.Any(), EqUpdatePasswordParams(arg, tUpdatedPassword)).
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

			server := newTestServer(t, store)
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
