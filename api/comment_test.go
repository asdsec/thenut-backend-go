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
	"github.com/stretchr/testify/require"
)

func TestCreatePostCommentAPI(t *testing.T) {
	user, _ := randomUser(t)
	merchant := randomMerchant(user.Username)
	post := randomPost(merchant.ID)
	comment := randomPostComment(post.ID)
	expected := newCommentResponse(comment)

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
				"post_id": post.ID,
				"comment": comment.Comment,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				arg := db.CreateCommentParams{
					CommentType: db.CommentTypePost,
					Owner:       user.Username,
					Comment:     comment.Comment,
					PostID: sql.NullInt64{
						Int64: post.ID,
						Valid: true,
					},
				}

				store.EXPECT().
					CreateComment(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(comment, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPostCommentResponse(t, recorder.Body, expected)
			},
		},
		{
			name: "InvalidPostID",
			body: gin.H{
				"post_id": -1, // invalid post id
				"comment": comment.Comment,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					CreateComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"post_id": post.ID,
				"comment": comment.Comment,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				// No Authorization
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					CreateComment(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"post_id": post.ID,
				"comment": comment.Comment,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				arg := db.CreateCommentParams{
					CommentType: db.CommentTypePost,
					Owner:       user.Username,
					Comment:     comment.Comment,
					PostID: sql.NullInt64{
						Int64: post.ID,
						Valid: true,
					},
				}

				store.EXPECT().
					CreateComment(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Comment{}, sql.ErrConnDone)
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

			url := "/posts/comments"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListPostCommentsAPI(t *testing.T) {
	user, _ := randomUser(t)
	merchant := randomMerchant(user.Username)
	post := randomPost(merchant.ID)

	n := 5
	comments := make([]db.Comment, n)
	expected := make([]commentResponse, n)
	for i := 0; i < n; i++ {
		comments[i] = randomPostComment(post.ID)
		expected[i] = newCommentResponse(comments[i])
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			body: gin.H{
				"post_id": post.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				arg := db.ListPostCommentsParams{
					PostID: sql.NullInt64{
						Int64: post.ID,
						Valid: true,
					},
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListPostComments(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(comments, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPostCommentsResponse(t, recorder.Body, expected)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1, // invalid page id
				pageSize: n,
			},
			body: gin.H{
				"post_id": post.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					ListPostComments(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageID:   1,
				pageSize: 100000, // invalid page size
			},
			body: gin.H{
				"post_id": post.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					ListPostComments(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidMerchantID",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			body: gin.H{
				"post_id": -1, // invalid merchant id
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					ListPostComments(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			body: gin.H{
				"post_id": post.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				// no authorization
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					ListPostComments(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			body: gin.H{
				"post_id": post.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					ListPostComments(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Comment{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			body: gin.H{
				"post_id": post.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					ListPostComments(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Comment{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			url := "/posts/comments"
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomPostComment(postID int64) db.Comment {
	return db.Comment{
		ID:          utils.RandomInt(1, 1000),
		CommentType: db.CommentTypePost,
		PostID: sql.NullInt64{
			Int64: postID,
			Valid: true,
		},
		Owner:   utils.RandomOwner(),
		Comment: utils.RandomString(14),
	}
}

func requireBodyMatchPostCommentsResponse(t *testing.T, body *bytes.Buffer, expected []commentResponse) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var actual []commentResponse
	err = json.Unmarshal(data, &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func requireBodyMatchPostCommentResponse(t *testing.T, body *bytes.Buffer, expected commentResponse) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var actual commentResponse
	err = json.Unmarshal(data, &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
