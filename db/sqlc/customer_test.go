package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/asdsec/thenut/utils"
	"github.com/stretchr/testify/require"
)

func createRandomCustomer(t *testing.T) Customer {
	tUser := createRandomUser(t)

	customer, err := testQueries.CreateCustomer(context.Background(), tUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, customer)
	require.Equal(t, tUser.Username, customer.Owner)
	require.NotZero(t, customer.CreatedAt)
	require.NotZero(t, customer.ID)
	require.NotZero(t, customer.ImageUrl)

	return customer
}

func TestCreateCustomer(t *testing.T) {
	tUser := createRandomUser(t)

	customer, err := testQueries.CreateCustomer(context.Background(), tUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, customer)
	require.Equal(t, tUser.Username, customer.Owner)
	require.NotZero(t, customer.CreatedAt)
	require.NotZero(t, customer.ID)
	require.NotZero(t, customer.ImageUrl)
}

func TestDeleteCustomer(t *testing.T) {
	tCustomer := createRandomCustomer(t)

	err := testQueries.DeleteCustomer(context.Background(), tCustomer.ID)

	require.NoError(t, err)
	customer, err := testQueries.GetCustomer(context.Background(), tCustomer.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, customer)
}

func TestGetCustomer(t *testing.T) {
	tCustomer := createRandomCustomer(t)

	customer, err := testQueries.GetCustomer(context.Background(), tCustomer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, customer)
	require.Equal(t, tCustomer.ID, customer.ID)
	require.Equal(t, tCustomer.ImageUrl, customer.ImageUrl)
	require.Equal(t, tCustomer.Owner, customer.Owner)
	require.WithinDuration(t, tCustomer.CreatedAt, customer.CreatedAt, time.Second)
}

func TestUpdateCustomer(t *testing.T) {
	tCustomer := createRandomCustomer(t)
	arg := UpdateCustomerParams{
		ID:       tCustomer.ID,
		ImageUrl: utils.RandomImageUrl(),
	}

	customer, err := testQueries.UpdateCustomer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, customer)
	require.Equal(t, arg.ImageUrl, customer.ImageUrl)
	require.Equal(t, tCustomer.ID, customer.ID)
	require.Equal(t, tCustomer.Owner, customer.Owner)
	require.WithinDuration(t, tCustomer.CreatedAt, customer.CreatedAt, time.Second)
}
