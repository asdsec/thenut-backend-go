package api

import (
	"database/sql"
	"net/http"

	db "github.com/asdsec/thenut/db/sqlc"
	"github.com/gin-gonic/gin"
)

const latest = "latest"

type createAppVersionRequest struct {
	Version string `json:"version" binding:"required"`
}

func (server *Server) createAppVersion(ctx *gin.Context) {
	var req createAppVersionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// todo: handle security

	oldVrs, err := server.store.GetAppVersion(ctx, latest)
	if err == nil {
		updateArg := db.UpdateAppVersionParams{
			ID:      oldVrs.ID,
			Version: oldVrs.Version,
			Tag:     "",
		}
		_, err = server.store.UpdateAppVersion(ctx, updateArg)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	arg := db.CreateAppVersionParams{
		Tag:     latest,
		Version: req.Version,
	}
	vrs, err := server.store.CreateAppVersion(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, vrs)
}

type versionResponse struct {
	Version string `json:"version"`
}

func (server *Server) getVersion(ctx *gin.Context) {
	vrs, err := server.store.GetAppVersion(ctx, latest)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, versionResponse{Version: vrs.Version})
}
