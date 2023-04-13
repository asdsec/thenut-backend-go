package db

import (
	"context"
	"reflect"
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
	merchant := createRandomMerchant(t)
	customer := createRandomCustomer(t)

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
}

func TestGetConsultancy(t *testing.T) {
	merchant := createRandomMerchant(t)
	customer := createRandomCustomer(t)
	consultancy := createRandomConsultancy(t, merchant, customer)

	consultancy, err := testQueries.GetConsultancy(context.Background(), consultancy.ID)

	require.NoError(t, err)
	require.NotEmpty(t, consultancy)
	require.Equal(t, merchant.ID, consultancy.MerchantID)
	require.Equal(t, customer.ID, consultancy.CustomerID)
	require.Equal(t, consultancy.ID, consultancy.ID)
	require.Equal(t, consultancy.Cost, consultancy.Cost)
	require.WithinDuration(t, consultancy.CreatedAt, consultancy.CreatedAt, time.Second)
}

func TestListConsultancies(t *testing.T) {
	var expected []Consultancy
	merchant := createRandomMerchant(t)
	customer := createRandomCustomer(t)
	for i := 0; i < 10; i++ {
		consultancy := createRandomConsultancy(t, merchant, customer)
		expected = append(expected, consultancy)
	}

	arg := ListConsultanciesParams{
		MerchantID: merchant.ID,
		CustomerID: customer.ID,
		Limit:      10,
		Offset:     0,
	}

	consultancies, err := testQueries.ListConsultancies(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, consultancies)
	require.Equal(t, len(expected), len(consultancies))
	require.True(t, reflect.DeepEqual(consultancies, expected))
}
