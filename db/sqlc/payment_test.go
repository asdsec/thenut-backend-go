package db

import (
	"context"
	"testing"
	"time"

	"github.com/sametdmr/thenut/utils"
	"github.com/stretchr/testify/require"
)

func createRandomPayment(t *testing.T, merchant Merchant, customer Customer) Payment {
	arg := CreatePaymentParams{
		MerchantID: merchant.ID,
		CustomerID: customer.ID,
		Amount:     utils.RandomMoney(),
	}

	payment, err := testQueries.CreatePayment(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, payment)
	require.Equal(t, arg.Amount, payment.Amount)
	require.Equal(t, arg.CustomerID, payment.CustomerID)
	require.Equal(t, arg.MerchantID, payment.MerchantID)
	require.NotZero(t, payment.CreatedAt)
	require.NotZero(t, payment.ID)

	return payment
}

func TestCreatePayment(t *testing.T) {
	tMerchant := createRandomMerchant(t)
	tCustomer := createRandomCustomer(t)
	arg := CreatePaymentParams{
		MerchantID: tMerchant.ID,
		CustomerID: tCustomer.ID,
		Amount:     utils.RandomMoney(),
	}

	payment, err := testQueries.CreatePayment(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, payment)
	require.Equal(t, arg.Amount, payment.Amount)
	require.Equal(t, arg.CustomerID, payment.CustomerID)
	require.Equal(t, arg.MerchantID, payment.MerchantID)
	require.NotZero(t, payment.CreatedAt)
	require.NotZero(t, payment.ID)
}

func TestGetPayment(t *testing.T) {
	tMerchant := createRandomMerchant(t)
	tCustomer := createRandomCustomer(t)
	tPayment := createRandomPayment(t, tMerchant, tCustomer)

	payment, err := testQueries.GetPayment(context.Background(), tPayment.ID)

	require.NoError(t, err)
	require.NotEmpty(t, payment)
	require.Equal(t, tPayment.Amount, payment.Amount)
	require.Equal(t, tPayment.CustomerID, payment.CustomerID)
	require.Equal(t, tPayment.MerchantID, payment.MerchantID)
	require.WithinDuration(t, tPayment.CreatedAt, payment.CreatedAt, time.Second)
	require.Equal(t, tPayment.ID, payment.ID)
}

func TestListPayments(t *testing.T) {
	var tLastMerchant Merchant
	var tLastCustomer Customer
	var tLastPayment Payment
	for i := 0; i < 10; i++ {
		tLastMerchant = createRandomMerchant(t)
		tLastCustomer = createRandomCustomer(t)
		tLastPayment = createRandomPayment(t, tLastMerchant, tLastCustomer)
	}
	arg := ListPaymentsParams{
		MerchantID: tLastMerchant.ID,
		CustomerID: tLastCustomer.ID,
		Limit:      5,
		Offset:     0,
	}

	payments, err := testQueries.ListPayments(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, payments)
	for _, payment := range payments {
		require.Equal(t, tLastPayment.Amount, payment.Amount)
		require.Equal(t, tLastPayment.CustomerID, payment.CustomerID)
		require.Equal(t, tLastPayment.MerchantID, payment.MerchantID)
		require.WithinDuration(t, tLastPayment.CreatedAt, payment.CreatedAt, time.Second)
		require.Equal(t, tLastPayment.ID, payment.ID)
	}
}
