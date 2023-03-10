package db

import (
	"context"
	"database/sql"
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

func TestCreateUser(t *testing.T) {
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
}

func TestDeleteUser(t *testing.T) {
	tUser := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), tUser.Username)

	require.NoError(t, err)
	user, err := testQueries.GetUser(context.Background(), tUser.Username)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user)
}

func TestGetUser(t *testing.T) {
	tUser := createRandomUser(t)

	user, err := testQueries.GetUser(context.Background(), tUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.WithinDuration(t, tUser.BirthDate, user.BirthDate, time.Second)
	require.WithinDuration(t, tUser.CreatedAt, user.CreatedAt, time.Second)
	require.Equal(t, tUser.Disabled, user.Disabled)
	require.Equal(t, tUser.Email, user.Email)
	require.Equal(t, tUser.FullName, user.FullName)
	require.Equal(t, tUser.Gender, user.Gender)
	require.Equal(t, tUser.HashedPassword, user.HashedPassword)
	require.Equal(t, tUser.ImageUrl, user.ImageUrl)
	require.WithinDuration(t, tUser.PasswordChangedAt, user.PasswordChangedAt, time.Second)
	require.Equal(t, tUser.PhoneNumber, user.PhoneNumber)
	require.Equal(t, tUser.Username, user.Username)
}

func TestUpdateEmail(t *testing.T) {
	tUser := createRandomUser(t)
	arg := UpdateEmailParams{
		Username: tUser.Username,
		Email:    utils.RandomEmail(),
	}

	user, err := testQueries.UpdateEmail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, tUser.Username, user.Username)
	require.WithinDuration(t, tUser.BirthDate, user.BirthDate, time.Second)
	require.WithinDuration(t, tUser.CreatedAt, user.CreatedAt, time.Second)
	require.Equal(t, tUser.Disabled, user.Disabled)
	require.Equal(t, tUser.FullName, user.FullName)
	require.Equal(t, tUser.Gender, user.Gender)
	require.Equal(t, tUser.HashedPassword, user.HashedPassword)
	require.Equal(t, tUser.ImageUrl, user.ImageUrl)
	require.WithinDuration(t, tUser.PasswordChangedAt, user.PasswordChangedAt, time.Second)
	require.Equal(t, tUser.PhoneNumber, user.PhoneNumber)
}

func TestUpdatePassword(t *testing.T) {
	tUser := createRandomUser(t)
	hashedPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	arg := UpdatePasswordParams{
		Username:       tUser.Username,
		HashedPassword: hashedPassword,
	}

	user, err := testQueries.UpdatePassword(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, tUser.Username, user.Username)
	require.Equal(t, tUser.Email, user.Email)
	require.WithinDuration(t, tUser.BirthDate, user.BirthDate, time.Second)
	require.WithinDuration(t, tUser.CreatedAt, user.CreatedAt, time.Second)
	require.Equal(t, tUser.Disabled, user.Disabled)
	require.Equal(t, tUser.FullName, user.FullName)
	require.Equal(t, tUser.Gender, user.Gender)
	require.Equal(t, tUser.ImageUrl, user.ImageUrl)
	require.WithinDuration(t, tUser.PasswordChangedAt, user.PasswordChangedAt, time.Second)
	require.Equal(t, tUser.PhoneNumber, user.PhoneNumber)
}

func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name          string
		arg           UpdateUserParams
		checkResponse func(t *testing.T, arg UpdateUserParams, tUser User, user User, err error)
	}{
		{
			name: "All Fields Update",
			arg: UpdateUserParams{
				FullName: sql.NullString{
					String: utils.RandomOwner(),
					Valid:  true,
				},
				PhoneNumber: sql.NullString{
					String: utils.RandomPhoneNumber(),
					Valid:  true,
				},
				Gender: sql.NullString{
					String: utils.RandomGender(),
					Valid:  true,
				},
				ImageUrl: sql.NullString{
					String: utils.RandomImageUrl(),
					Valid:  true,
				},
				BirthDate: sql.NullTime{
					Time:  utils.RandomBirthDate(),
					Valid: true,
				},
			},
			checkResponse: func(t *testing.T, arg UpdateUserParams, tUser User, user User, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, user)
				require.WithinDuration(t, arg.BirthDate.Time, user.BirthDate, time.Second)
				require.Equal(t, arg.FullName.String, user.FullName)
				require.Equal(t, arg.Gender.String, user.Gender)
				require.Equal(t, arg.PhoneNumber.String, user.PhoneNumber)
				require.Equal(t, arg.ImageUrl.String, user.ImageUrl)
				require.Equal(t, tUser.Username, user.Username)
				require.Equal(t, tUser.Disabled, user.Disabled)
				require.Equal(t, tUser.Email, user.Email)
				require.Equal(t, tUser.HashedPassword, user.HashedPassword)
				require.WithinDuration(t, tUser.PasswordChangedAt, user.PasswordChangedAt, time.Second)
				require.WithinDuration(t, tUser.CreatedAt, user.CreatedAt, time.Second)
			},
		},
		{
			name: "ImageUrl Update",
			arg: UpdateUserParams{
				ImageUrl: sql.NullString{
					String: utils.RandomImageUrl(),
					Valid:  true,
				},
			},
			checkResponse: func(t *testing.T, arg UpdateUserParams, tUser User, user User, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, user)
				require.Equal(t, arg.ImageUrl.String, user.ImageUrl)
				require.WithinDuration(t, tUser.BirthDate, user.BirthDate, time.Second)
				require.Equal(t, tUser.FullName, user.FullName)
				require.Equal(t, tUser.Gender, user.Gender)
				require.Equal(t, tUser.PhoneNumber, user.PhoneNumber)
				require.Equal(t, tUser.Username, user.Username)
				require.Equal(t, tUser.Disabled, user.Disabled)
				require.Equal(t, tUser.Email, user.Email)
				require.Equal(t, tUser.HashedPassword, user.HashedPassword)
				require.WithinDuration(t, tUser.PasswordChangedAt, user.PasswordChangedAt, time.Second)
				require.WithinDuration(t, tUser.CreatedAt, user.CreatedAt, time.Second)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tUser := createRandomUser(t)
			tc.arg.Username = tUser.Username
			user, err := testQueries.UpdateUser(context.Background(), tc.arg)
			tc.checkResponse(t, tc.arg, tUser, user, err)
		})
	}
}
