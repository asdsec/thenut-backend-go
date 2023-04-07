package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
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
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestDeleteCustomerAPI(t *testing.T) {
	user, _ := randomUser(t)
	customer := randomCustomer(user.Username)

	testCases := []struct {
		name          string
		customerID    int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "Ok",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(customer, nil)

				store.EXPECT().
					DeleteCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "InvalidCustomerID",
			customerID: -1,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					DeleteCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "UnauthenticatedUser",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "invalid_username", time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(customer, nil)

				store.EXPECT().
					DeleteCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:       "NoAuthentication",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				// No authentication
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					DeleteCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:       "NotFoundGetCustomer",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(db.Customer{}, sql.ErrNoRows)

				store.EXPECT().
					DeleteCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalServerErrorGetCustomer",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(db.Customer{}, sql.ErrConnDone)

				store.EXPECT().
					DeleteCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "NotFoundUpdateCustomer",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(customer, nil)

				store.EXPECT().
					DeleteCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalServerErrorUpdateCustomer",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(customer, nil)

				store.EXPECT().
					DeleteCustomer(gomock.Any(), gomock.Eq(customer.ID)).
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
			tc.buildStubs(store)

			tokenMaker := newTestTokenMaker(t)
			server := newTestServer(t, store, tokenMaker)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.customerID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateCustomerAPI(t *testing.T) {
	user, _ := randomUser(t)
	customer := randomCustomer(user.Username)
	expected := db.Customer{
		ID:        customer.ID,
		Owner:     customer.Owner,
		ImageUrl:  utils.RandomImageUrl(),
		CreatedAt: customer.CreatedAt,
	}

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
				"id":        customer.ID,
				"image_url": expected.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(customer, nil)

				arg := db.UpdateCustomerParams{
					ID:       customer.ID,
					ImageUrl: expected.ImageUrl,
				}

				store.EXPECT().
					UpdateCustomer(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(expected, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, expected)
			},
		},
		{
			name: "InvalidCustomerID",
			body: gin.H{
				"id":        "invld", // invalid id
				"image_url": expected.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					UpdateCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UnauthenticatedUser",
			body: gin.H{
				"id":        customer.ID,
				"image_url": expected.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "invalid_username", time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(customer, nil)

				store.EXPECT().
					UpdateCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthentication",
			body: gin.H{
				"id":        customer.ID,
				"image_url": expected.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				// No authentication
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					UpdateCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NotFoundGetCustomer",
			body: gin.H{
				"id":        customer.ID,
				"image_url": expected.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(db.Customer{}, sql.ErrNoRows)

				store.EXPECT().
					UpdateCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerErrorGetCustomer",
			body: gin.H{
				"id":        customer.ID,
				"image_url": expected.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(db.Customer{}, sql.ErrConnDone)

				store.EXPECT().
					UpdateCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFoundUpdateCustomer",
			body: gin.H{
				"id":        customer.ID,
				"image_url": expected.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(customer, nil)

				arg := db.UpdateCustomerParams{
					ID:       customer.ID,
					ImageUrl: expected.ImageUrl,
				}

				store.EXPECT().
					UpdateCustomer(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Customer{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerErrorUpdateCustomer",
			body: gin.H{
				"id":        customer.ID,
				"image_url": expected.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(customer, nil)

				arg := db.UpdateCustomerParams{
					ID:       customer.ID,
					ImageUrl: expected.ImageUrl,
				}

				store.EXPECT().
					UpdateCustomer(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Customer{}, sql.ErrConnDone)
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
			tc.buildStubs(store)

			tokenMaker := newTestTokenMaker(t)
			server := newTestServer(t, store, tokenMaker)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetCustomerAPI(t *testing.T) {
	user, _ := randomUser(t)
	customer := randomCustomer(user.Username)

	testCases := []struct {
		name          string
		customerID    int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "Ok",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(customer, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, customer)
			},
		},
		{
			name:       "InvalidCustomerId",
			customerID: -1, // invalid id
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "UnauthenticatedUser",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "invalid_username", time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(customer, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:       "NoAuthentication",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				// No authentication
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:       "NotFound",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(db.Customer{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalServerError",
			customerID: customer.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetCustomer(gomock.Any(), gomock.Eq(customer.ID)).
					Times(1).
					Return(db.Customer{}, sql.ErrConnDone)
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
			tc.buildStubs(store)

			tokenMaker := newTestTokenMaker(t)
			server := newTestServer(t, store, tokenMaker)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.customerID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateCustomerAPI(t *testing.T) {
	user, _ := randomUser(t)
	customer := randomCustomer(user.Username)

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
				"owner": customer.Owner,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					CreateCustomer(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(customer, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, customer)
			},
		},
		{
			name: "InvalidOwner",
			body: gin.H{
				"owner": "invld", // invalid owner
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					CreateCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Forbidden foreign_key_violation",
			body: gin.H{
				"owner": customer.Owner,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					CreateCustomer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Customer{}, &pq.Error{Code: "23503"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "Forbidden unique_violation",
			body: gin.H{
				"owner": customer.Owner,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					CreateCustomer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Customer{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"owner": customer.Owner,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					CreateCustomer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Customer{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "UnauthenticatedUser",
			body: gin.H{
				"owner": customer.Owner,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "invalid_username", time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					CreateCustomer(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthentication",
			body: gin.H{
				"owner": customer.Owner,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				// No authentication
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					CreateCustomer(gomock.Any(), gomock.Any()).
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomCustomer(owner string) db.Customer {
	return db.Customer{
		ID:        utils.RandomInt(1, 1000),
		Owner:     owner,
		ImageUrl:  utils.RandomImageUrl(),
		CreatedAt: utils.RandomBirthDate(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, customer db.Customer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotCustomer db.Customer
	err = json.Unmarshal(data, &gotCustomer)
	require.NoError(t, err)
	require.WithinDuration(t, customer.CreatedAt, gotCustomer.CreatedAt, time.Second)
	require.Equal(t, customer.ID, gotCustomer.ID)
	require.Equal(t, customer.ImageUrl, gotCustomer.ImageUrl)
	require.Equal(t, customer.Owner, gotCustomer.Owner)
}
