package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/asdsec/thenut/db/sqlc"
	"github.com/gin-gonic/gin"
)

type commentResponse struct {
	ID         int64     `json:"id"`
	PostID     int64     `json:"post_id,omitempty"`
	MerchantID int64     `json:"merchant_id,omitempty"`
	Owner      string    `json:"owner"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
}

func newCommentResponse(c db.Comment) commentResponse {
	return commentResponse{
		ID:         c.ID,
		PostID:     c.PostID.Int64,
		MerchantID: c.MerchantID.Int64,
		Owner:      c.Owner,
		Comment:    c.Comment,
		CreatedAt:  c.CreatedAt,
	}
}

type listPostCommentsRequestUri struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type listPostCommentsRequest struct {
	PostID int64 `json:"post_id" binding:"required,min=1"`
}

func (server *Server) listPostComments(ctx *gin.Context) {
	var req listPostCommentsRequest
	var query listPostCommentsRequestUri
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListPostCommentsParams{
		PostID: sql.NullInt64{
			Int64: req.PostID,
			Valid: true,
		},
		Limit:  query.PageSize,
		Offset: (query.PageID - 1) * query.PageSize,
	}

	comments, err := server.store.ListPostComments(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := make([]commentResponse, len(comments))
	for i := range comments {
		rsp[i] = newCommentResponse(comments[i])
	}

	ctx.JSON(http.StatusOK, rsp)
}
