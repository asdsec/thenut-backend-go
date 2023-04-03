package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/asdsec/thenut/db/sqlc"
	"github.com/asdsec/thenut/token"
	"github.com/asdsec/thenut/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type authResponse struct {
	AccessToken  string `json:"access_token"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Username     string `json:"username"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type registerUserRequest struct {
	Username    string    `json:"username" binding:"required,alphanum,min=6"`
	Password    string    `json:"password" binding:"required,min=6,max=72"`
	Email       string    `json:"email" binding:"required,email"`
	FullName    string    `json:"full_name" binding:"required"`
	PhoneNumber string    `json:"phone_number" binding:"required,min=11"`
	Gender      string    `json:"gender" binding:"required,min=1,max=1"`
	BirthDate   time.Time `json:"birth_date" binding:"required"`
}

func (server *Server) registerUser(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
		PhoneNumber:    req.PhoneNumber,
		Gender:         req.Gender,
		BirthDate:      req.BirthDate,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := authResponse{
		AccessToken:  accessToken,
		Email:        user.Email,
		FullName:     user.FullName,
		Username:     user.Username,
		RefreshToken: refreshToken,
		ExpiresIn:    int(server.config.AccessTokenDuration.Seconds()),
	}
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, ok := getUserFromStore(ctx, server, req.Username)
	if !ok {
		return
	}

	err := utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := authResponse{
		AccessToken:  accessToken,
		Email:        user.Email,
		FullName:     user.FullName,
		Username:     user.Username,
		RefreshToken: refreshToken,
		ExpiresIn:    int(server.config.AccessTokenDuration.Seconds()),
	}
	ctx.JSON(http.StatusOK, rsp)
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PhoneNumber       string    `json:"phone_number"`
	ImageUrl          string    `json:"image_url"`
	Gender            string    `json:"gender"`
	Disabled          bool      `json:"disabled"`
	BirthDate         time.Time `json:"birth_date"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		Email:             user.Email,
		FullName:          user.FullName,
		PhoneNumber:       user.PhoneNumber,
		Gender:            user.Gender,
		BirthDate:         user.BirthDate,
		ImageUrl:          user.ImageUrl,
		Disabled:          user.Disabled,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

type getUserRequest struct {
	Username string `uri:"username" binding:"required,min=6"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, ok := getUserFromStore(ctx, server, req.Username)
	if !ok {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.TokenPayload)
	if user.Username != authPayload.Username {
		err := errors.New("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type updateEmailRequest struct {
	Username string `json:"username" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

func (server *Server) updateEmail(ctx *gin.Context) {
	var req updateEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.TokenPayload)
	if req.Username != authPayload.Username {
		err := errors.New("username does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateEmailParams{
		Username: req.Username,
		Email:    req.Email,
	}

	user, err := server.store.UpdateEmail(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type updatePasswordRequest struct {
	Username    string `json:"username" binding:"required,min=6"`
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (server *Server) updatePassword(ctx *gin.Context) {
	var req updatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.TokenPayload)
	if req.Username != authPayload.Username {
		err := errors.New("username does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	user, ok := getUserFromStore(ctx, server, req.Username)
	if !ok {
		return
	}

	err := utils.CheckPassword(req.OldPassword, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdatePasswordParams{
		Username:          req.Username,
		HashedPassword:    hashedPassword,
		PasswordChangedAt: time.Now(),
	}

	user, err = server.store.UpdatePassword(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type updateUserRequest struct {
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	Gender      string    `json:"gender"`
	BirthDate   time.Time `json:"birth_date"`
	ImageUrl    string    `json:"image_url"`
	Username    string    `json:"username" binding:"required,min=6"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.TokenPayload)
	if req.Username != authPayload.Username {
		err := errors.New("username does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		Username: req.Username,
		FullName: sql.NullString{
			String: req.FullName,
			Valid:  len(req.FullName) > 0,
		},
		PhoneNumber: sql.NullString{
			String: req.PhoneNumber,
			Valid:  len(req.PhoneNumber) > 0,
		},
		Gender: sql.NullString{
			String: req.Gender,
			Valid:  len(req.Gender) > 0,
		},
		ImageUrl: sql.NullString{
			String: req.ImageUrl,
			Valid:  len(req.ImageUrl) > 0,
		},
		BirthDate: sql.NullTime{
			Time:  req.BirthDate,
			Valid: true,
		},
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

func getUserFromStore(ctx *gin.Context, server *Server, username string) (db.User, bool) {
	user, err := server.store.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return user, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return user, false
	}
	return user, true
}
