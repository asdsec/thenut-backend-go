package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/asdsec/thenut/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createMerchantRequest struct {
	Owner      string `json:"owner" binding:"required,alphanum,min=6"`
	Profession string `json:"profession" binding:"required"`
	Title      string `json:"title" binding:"required"`
	About      string `json:"about" binding:"required"`
}

func (server *Server) createMerchant(ctx *gin.Context) {
	var req createMerchantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := server.getAuthPayload(ctx)
	if req.Owner != authPayload.Username {
		err := errors.New("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.CreateMerchantParams{
		Owner: req.Owner,
		// todo: delete balance for creation of merchant
		Profession: req.Profession,
		Title:      req.Title,
		About:      req.About,
	}

	merchant, err := server.store.CreateMerchant(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, merchant)
}

type deleteMerchantRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteMerchant(ctx *gin.Context) {
	var req deleteMerchantRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	merchant, err := server.store.GetMerchant(ctx, req.ID)
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

	err = server.store.DeleteMerchant(ctx, req.ID)
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

type getMerchantRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getMerchant(ctx *gin.Context) {
	var req getMerchantRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	merchant, err := server.store.GetMerchant(ctx, req.ID)
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

	ctx.JSON(http.StatusOK, merchant)
}

type listMerchantsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listMerchants(ctx *gin.Context) {
	var req listMerchantsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := server.getAuthPayload(ctx)
	arg := db.ListMerchantsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	merchants, err := server.store.ListMerchants(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, merchants)
}

type updateMerchantRequest struct {
	Profession string `json:"profession"`
	Title      string `json:"title"`
	About      string `json:"about"`
	ImageUrl   string `json:"image_url"`
	Rating     int32  `json:"rating"`
	ID         int64  `json:"id" binding:"required,min=1"`
}

func (server *Server) updateMerchant(ctx *gin.Context) {
	var req updateMerchantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	merchant, err := server.store.GetMerchant(ctx, req.ID)
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
		err := errors.New("username does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateMerchantParams{
		ID: req.ID,
		Profession: sql.NullString{
			String: req.Profession,
			Valid:  len(req.Profession) > 0,
		},
		Title: sql.NullString{
			String: req.Title,
			Valid:  len(req.Title) > 0,
		},
		About: sql.NullString{
			String: req.Title,
			Valid:  len(req.Title) > 0,
		},
		ImageUrl: sql.NullString{
			String: req.Title,
			Valid:  len(req.Title) > 0,
		},
		// fixme: rating valid
		Rating: sql.NullFloat64{
			Float64: float64(req.Rating),
			Valid:   req.Rating != 0,
		},
	}

	merchant, err = server.store.UpdateMerchant(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, merchant)
}
