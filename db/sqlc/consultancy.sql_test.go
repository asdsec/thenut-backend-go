package db

import (
	"context"
	"testing"
	"time"

	"github.com/asdsec/thenut/utils"
	"github.com/stretchr/testify/require"
)

func createRandomConsultancy(t *testing.T, merchant Merchant, customer Customer) Consultancy {
	arg := CreateConsultancyParams{
		MerchantID: merchant.ID,
		CustomerID: customer.ID,
		Cost:       utils.RandomMoney(),
	}

	consultancy, err := testQueries.CreateConsultancy(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, consultancy)
	require.Equal(t, arg.Cost, consultancy.Cost)
	require.Equal(t, arg.CustomerID, consultancy.CustomerID)
	require.Equal(t, arg.MerchantID, consultancy.MerchantID)
	require.NotZero(t, consultancy.ID)
	require.NotZero(t, consultancy.CreatedAt)

	return consultancy
}

func TestCreateConsultancy(t *testing.T) {
	tMerchant := createRandomMerchant(t)
	tCustomer := createRandomCustomer(t)

	arg := CreateConsultancyParams{
		MerchantID: tMerchant.ID,
		CustomerID: tCustomer.ID,
		Cost:       utils.RandomMoney(),
	}

	consultancy, err := testQueries.CreateConsultancy(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, consultancy)
	require.Equal(t, arg.Cost, consultancy.Cost)
	require.Equal(t, arg.CustomerID, consultancy.CustomerID)
	require.Equal(t, arg.MerchantID, consultancy.MerchantID)
	require.NotZero(t, consultancy.ID)
	require.NotZero(t, consultancy.CreatedAt)
}

func TestGetConsultancy(t *testing.T) {
	tMerchant := createRandomMerchant(t)
	tCustomer := createRandomCustomer(t)
	tConsultancy := createRandomConsultancy(t, tMerchant, tCustomer)

	consultancy, err := testQueries.GetConsultancy(context.Background(), tConsultancy.ID)

	require.NoError(t, err)
	require.NotEmpty(t, consultancy)
	require.Equal(t, tMerchant.ID, consultancy.MerchantID)
	require.Equal(t, tCustomer.ID, consultancy.CustomerID)
	require.Equal(t, tConsultancy.ID, consultancy.ID)
	require.Equal(t, tConsultancy.Cost, consultancy.Cost)
	require.WithinDuration(t, tConsultancy.CreatedAt, consultancy.CreatedAt, time.Second)
}

func TestListConsultancies(t *testing.T) {
	var tLastConsultancy Consultancy
	for i := 0; i < 10; i++ {
		tMerchant := createRandomMerchant(t)
		tCustomer := createRandomCustomer(t)
		tLastConsultancy = createRandomConsultancy(t, tMerchant, tCustomer)
	}
	arg := ListConsultanciesParams{
		MerchantID: tLastConsultancy.MerchantID,
		CustomerID: tLastConsultancy.CustomerID,
		Limit:      5,
		Offset:     0,
	}

	consultancies, err := testQueries.ListConsultancies(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, consultancies)
	for _, consultancy := range consultancies {
		require.NotEmpty(t, consultancy)
		require.Equal(t, tLastConsultancy.Cost, consultancy.Cost)
		require.WithinDuration(t, tLastConsultancy.CreatedAt, consultancy.CreatedAt, time.Second)
		require.Equal(t, tLastConsultancy.CustomerID, consultancy.CustomerID)
		require.Equal(t, tLastConsultancy.ID, consultancy.ID)
		require.Equal(t, tLastConsultancy.MerchantID, consultancy.MerchantID)
	}
}
