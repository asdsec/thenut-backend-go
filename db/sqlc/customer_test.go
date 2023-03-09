package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateCustomer(t *testing.T) {
	user := createRandomUser(t)

	customer, err := testQueries.CreateCustomer(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, customer)
	require.Equal(t, user.Username, customer.Owner)
	require.NotZero(t, customer.CreatedAt)
	require.NotZero(t, customer.ID)
	require.NotZero(t, customer.ImageUrl)
}
