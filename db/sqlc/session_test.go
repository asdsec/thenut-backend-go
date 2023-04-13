package db

import (
	"context"
	"testing"

	"github.com/asdsec/thenut/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomSession(t *testing.T, user User) Session {
	arg := CreateSessionParams{
		ID:           uuid.New(),
		Username:     user.Username,
		RefreshToken: utils.RandomString(12),
		UserAgent:    utils.RandomString(12),
		ClientIp:     utils.RandomString(12),
		IsBlocked:    false,
	}

	err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	session, err := testQueries.GetSession(context.Background(), arg.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session)
	require.Equal(t, arg.ClientIp, session.ClientIp)
	require.Equal(t, arg.ExpiresAt, session.ExpiresAt)
	require.Equal(t, arg.ID, session.ID)
	require.Equal(t, arg.IsBlocked, session.IsBlocked)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.Username, session.Username)

	return session
}

func TestCreateSession(t *testing.T) {
	user := createRandomUser(t)
	arg := CreateSessionParams{
		ID:           uuid.New(),
		Username:     user.Username,
		RefreshToken: utils.RandomString(12),
		UserAgent:    utils.RandomString(12),
		ClientIp:     utils.RandomString(12),
		IsBlocked:    false,
	}

	err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	session, err := testQueries.GetSession(context.Background(), arg.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session)
	require.Equal(t, arg.ClientIp, session.ClientIp)
	require.Equal(t, arg.ExpiresAt, session.ExpiresAt)
	require.Equal(t, arg.ID, session.ID)
	require.Equal(t, arg.IsBlocked, session.IsBlocked)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.Username, session.Username)
}

func TestGetSession(t *testing.T) {
	user := createRandomUser(t)
	expected := createRandomSession(t, user)

	session, err := testQueries.GetSession(context.Background(), expected.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session)
	require.Equal(t, expected, session)
}
