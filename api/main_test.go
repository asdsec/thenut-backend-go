package api

import (
	"os"
	"testing"
	"time"

	db "github.com/asdsec/thenut/db/sqlc"
	"github.com/asdsec/thenut/token"
	"github.com/asdsec/thenut/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var testConfig = utils.Config{
	TokenSymmetricKey:    utils.RandomString(32),
	AccessTokenDuration:  time.Minute,
	RefreshTokenDuration: time.Hour,
	ServerAddress:        "http://localhost:8080",
}

func newTestServer(t *testing.T, store db.Store, tokenMaker token.TokenMaker) *Server {
	server, err := NewServer(testConfig, store, tokenMaker)
	require.NoError(t, err)

	return server
}

func newTestTokenMaker(t *testing.T) token.TokenMaker {
	tokenMaker, err := token.NewPasetoMaker(testConfig.TokenSymmetricKey)
	require.NoError(t, err)

	return tokenMaker
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
