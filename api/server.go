package api

import (
	"fmt"

	db "github.com/asdsec/thenut/db/sqlc"
	"github.com/asdsec/thenut/token"
	"github.com/asdsec/thenut/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.TokenMaker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and routing
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("gender", validGender)
	}

	server.setupRouter()
	return server, nil
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/auth/register", server.registerUser)
	router.POST("/auth/login", server.loginUser)
	router.POST("/tokens/renew", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/:username", server.getUser)
	// todo: add more auth method for updating email, like email verification
	authRoutes.POST("/users/email", server.updateEmail)
	// todo: add more auth method for updating password, like email verification
	authRoutes.POST("/users/password", server.updatePassword)
	authRoutes.PATCH("/users", server.updateUser)

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
