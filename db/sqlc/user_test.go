package db

import (
	"context"
	"testing"
	"time"

	"github.com/sametdmr/thenut/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
		PhoneNumber:    utils.RandomPhoneNumber(),
		Gender:         utils.RandomGender(),
		BirthDate:      utils.RandomBirthDate(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.WithinDuration(t, arg.BirthDate, user.BirthDate, time.Second)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Gender, user.Gender)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.PhoneNumber, user.PhoneNumber)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	require.False(t, user.Disabled)
	require.NotZero(t, user.ImageUrl)

	return user
}
