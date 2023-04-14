package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/asdsec/thenut/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createPostRequest struct {
	MerchantID int64  `json:"merchant_id" binding:"required,min=1"`
	Title      string `json:"title"`
	// implement real image
	ImageUrl string `json:"image_url"`
}

type postResponse struct {
	ID         int64     `json:"id"`
	MerchantID int64     `json:"merchant_id"`
	Title      string    `json:"title,omitempty"`
	ImageUrl   string    `json:"image_url,omitempty"`
	Likes      int32     `json:"likes"`
	CreatedAt  time.Time `json:"created_at"`
}

func newPostResponse(post db.Post) postResponse {
	return postResponse{
		ID:         post.ID,
		MerchantID: post.MerchantID,
		Title:      post.Title.String,
		ImageUrl:   post.ImageUrl.String,
		Likes:      post.Likes,
		CreatedAt:  post.CreatedAt,
	}
}

func (server *Server) createPost(ctx *gin.Context) {
	var req createPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.ImageUrl == "" && req.Title == "" {
		err := errors.New("both title and image_url is null")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	merchant, err := server.store.GetMerchant(ctx, req.MerchantID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	authPayload := server.getAuthPayload(ctx)
	if merchant.Owner != authPayload.Username {
		err := errors.New("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.CreatePostParams{
		MerchantID: req.MerchantID,
		Title: sql.NullString{
			String: req.Title,
			Valid:  req.Title != "",
		},
		ImageUrl: sql.NullString{
			String: req.ImageUrl,
			Valid:  req.ImageUrl != "",
		},
	}

	post, err := server.store.CreatePost(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newPostResponse(post))
}

type deletePostRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deletePost(ctx *gin.Context) {
	var req deletePostRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	post, err := server.store.GetPost(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	merchant, err := server.store.GetMerchant(ctx, post.MerchantID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	authPayload := server.getAuthPayload(ctx)
	if merchant.Owner != authPayload.Username {
		err := errors.New("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeletePost(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type listMerchantPostsRequestQuery struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type listMerchantPostsRequest struct {
	MerchantID int64 `json:"merchant_id" binding:"required,min=1"`
}

func (server *Server) listMerchantPosts(ctx *gin.Context) {
	var req listMerchantPostsRequest
	var query listMerchantPostsRequestQuery
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListMerchantPostsParams{
		MerchantID: req.MerchantID,
		Limit:      query.PageSize,
		Offset:     (query.PageID - 1) * query.PageSize,
	}

	posts, err := server.store.ListMerchantPosts(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := make([]postResponse, len(posts))
	for i := range posts {
		rsp[i] = newPostResponse(posts[i])
	}

	ctx.JSON(http.StatusOK, rsp)
}
