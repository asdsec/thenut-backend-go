package db

import (
	"context"
	"testing"
	"time"

	"github.com/sametdmr/thenut/utils"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, customer Customer) Entry {
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

	return entry
}

func TestCreateEntry(t *testing.T) {
	tCustomer := createRandomCustomer(t)
	arg := CreateEntryParams{
		CustomerID: tCustomer.ID,
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

func TestGetEntry(t *testing.T) {
	tCustomer := createRandomCustomer(t)
	tEntry := createRandomEntry(t, tCustomer)

	entry, err := testQueries.GetEntry(context.Background(), tEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, tEntry.Amount, entry.Amount)
	require.WithinDuration(t, tEntry.CreatedAt, entry.CreatedAt, time.Second)
	require.Equal(t, tEntry.CustomerID, entry.CustomerID)
	require.Equal(t, tEntry.ID, entry.ID)
	require.Equal(t, tEntry.Type, entry.Type)
}

func TestListEntries(t *testing.T) {
	var tLastEntry Entry
	for i := 0; i < 10; i++ {
		tCustomer := createRandomCustomer(t)
		tLastEntry = createRandomEntry(t, tCustomer)
	}
	arg := ListEntriesParams{
		CustomerID: tLastEntry.CustomerID,
		Limit:      5,
		Offset:     0,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entries)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, tLastEntry.Amount, entry.Amount)
		require.WithinDuration(t, tLastEntry.CreatedAt, entry.CreatedAt, time.Second)
		require.Equal(t, tLastEntry.CustomerID, entry.CustomerID)
		require.Equal(t, tLastEntry.ID, entry.ID)
		require.Equal(t, tLastEntry.Type, entry.Type)
	}
}
