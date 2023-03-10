package db

import (
	"context"
	"testing"

	"github.com/sametdmr/thenut/utils"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, customer Customer) {
	arg := CreateEntryParams{
		CustomerID: customer.ID,
		Amount:     utils.RandomMoney(),
		Type:       utils.RandomString(3),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.Amount, entry.Amount)
	require.Equal(t, arg.CustomerID, entry.CustomerID)
	require.Equal(t, arg.Type, entry.Type)
	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.ID)
}
