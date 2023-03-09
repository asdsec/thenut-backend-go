package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/sametdmr/thenut/utils"
	"github.com/stretchr/testify/require"
)

func createRandomMerchant(t *testing.T) Merchant {
	tUser := createRandomUser(t)
	arg := CreateMerchantParams{
		Owner:      tUser.Username,
		Balance:    utils.RandomMoney(),
		Profession: utils.RandomString(12),
		Title:      utils.RandomString(6),
		About:      utils.RandomString(10),
	}

	merchant, err := testQueries.CreateMerchant(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, merchant)
	require.Equal(t, arg.Owner, merchant.Owner)
	require.Equal(t, arg.Balance, merchant.Balance)
	require.Equal(t, arg.About, merchant.About)
	require.Equal(t, arg.Profession, merchant.Profession)
	require.Equal(t, arg.Title, merchant.Title)
	require.NotZero(t, merchant.CreatedAt)
	require.NotZero(t, merchant.ID)
	require.NotZero(t, merchant.ImageUrl)
	require.Zero(t, merchant.Rating)

	return merchant
}

func TestAddMerchantBalance(t *testing.T) {
	tMerchant := createRandomMerchant(t)
	arg := AddMerchantBalanceParams{
		ID:     tMerchant.ID,
		Amount: utils.RandomMoney(),
	}

	merchant, err := testQueries.AddMerchantBalance(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, merchant)
	require.Equal(t, arg.Amount+tMerchant.Balance, merchant.Balance)
	require.Equal(t, arg.ID, merchant.ID)
	require.Equal(t, tMerchant.About, merchant.About)
	require.Equal(t, tMerchant.ImageUrl, merchant.ImageUrl)
	require.Equal(t, tMerchant.Owner, merchant.Owner)
	require.Equal(t, tMerchant.Profession, merchant.Profession)
	require.Equal(t, tMerchant.Rating, merchant.Rating)
	require.Equal(t, tMerchant.Title, merchant.Title)
	require.WithinDuration(t, tMerchant.CreatedAt, merchant.CreatedAt, time.Second)
}

func TestCreateMerchant(t *testing.T) {
	tUser := createRandomUser(t)
	arg := CreateMerchantParams{
		Owner:      tUser.Username,
		Balance:    utils.RandomMoney(),
		Profession: utils.RandomString(12),
		Title:      utils.RandomString(6),
		About:      utils.RandomString(10),
	}

	merchant, err := testQueries.CreateMerchant(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, merchant)
	require.Equal(t, arg.Owner, merchant.Owner)
	require.Equal(t, arg.Balance, merchant.Balance)
	require.Equal(t, arg.About, merchant.About)
	require.Equal(t, arg.Profession, merchant.Profession)
	require.Equal(t, arg.Title, merchant.Title)
	require.NotZero(t, merchant.CreatedAt)
	require.NotZero(t, merchant.ID)
	require.NotZero(t, merchant.ImageUrl)
	require.Zero(t, merchant.Rating)
}

func TestDeleteMerchant(t *testing.T) {
	tMerchant := createRandomMerchant(t)

	err := testQueries.DeleteMerchant(context.Background(), tMerchant.ID)

	require.NoError(t, err)
	merchant, err := testQueries.GetMerchant(context.Background(), tMerchant.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, merchant)
}

func TestGetMerchant(t *testing.T) {
	tMerchant := createRandomMerchant(t)

	merchant, err := testQueries.GetMerchant(context.Background(), tMerchant.ID)

	require.NoError(t, err)
	require.NotEmpty(t, merchant)
	require.Equal(t, tMerchant.About, merchant.About)
	require.Equal(t, tMerchant.Balance, merchant.Balance)
	require.WithinDuration(t, tMerchant.CreatedAt, merchant.CreatedAt, time.Second)
	require.Equal(t, tMerchant.ID, merchant.ID)
	require.Equal(t, tMerchant.ImageUrl, merchant.ImageUrl)
	require.Equal(t, tMerchant.Owner, merchant.Owner)
	require.Equal(t, tMerchant.Profession, merchant.Profession)
	require.Equal(t, tMerchant.Rating, merchant.Rating)
	require.Equal(t, tMerchant.Title, merchant.Title)
}

func TestGetMerchantForUpdate(t *testing.T) {
	tMerchant := createRandomMerchant(t)

	merchant, err := testQueries.GetMerchantForUpdate(context.Background(), tMerchant.ID)

	require.NoError(t, err)
	require.NotEmpty(t, merchant)
	require.Equal(t, tMerchant.About, merchant.About)
	require.Equal(t, tMerchant.Balance, merchant.Balance)
	require.WithinDuration(t, tMerchant.CreatedAt, merchant.CreatedAt, time.Second)
	require.Equal(t, tMerchant.ID, merchant.ID)
	require.Equal(t, tMerchant.ImageUrl, merchant.ImageUrl)
	require.Equal(t, tMerchant.Owner, merchant.Owner)
	require.Equal(t, tMerchant.Profession, merchant.Profession)
	require.Equal(t, tMerchant.Rating, merchant.Rating)
	require.Equal(t, tMerchant.Title, merchant.Title)
}

func TestListMerchants(t *testing.T) {
	var tLastMerchant Merchant
	for i := 0; i < 10; i++ {
		tLastMerchant = createRandomMerchant(t)
	}
	arg := ListMerchantsParams{
		Owner:  tLastMerchant.Owner,
		Limit:  5,
		Offset: 0,
	}

	merchants, err := testQueries.ListMerchants(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, merchants)
	for _, merchant := range merchants {
		require.NotEmpty(t, merchant)
		require.Equal(t, tLastMerchant.Owner, merchant.Owner)
	}
}

func TestUpdateMerchant(t *testing.T) {
	tMerchant := createRandomMerchant(t)
	arg := UpdateMerchantParams{
		ID:         tMerchant.ID,
		Balance:    utils.RandomMoney(),
		Profession: utils.RandomString(6),
		Title:      utils.RandomString(6),
		About:      utils.RandomString(12),
		ImageUrl:   utils.RandomImageUrl(),
		Rating:     int32(utils.RandomInt(0, 10)),
	}

	merchant, err := testQueries.UpdateMerchant(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, merchant)
	require.Equal(t, arg.About, merchant.About)
	require.Equal(t, arg.Balance, merchant.Balance)
	require.Equal(t, arg.ID, merchant.ID)
	require.Equal(t, arg.ImageUrl, merchant.ImageUrl)
	require.Equal(t, arg.Profession, merchant.Profession)
	require.Equal(t, arg.Rating, merchant.Rating)
	require.Equal(t, arg.Title, merchant.Title)
	require.WithinDuration(t, tMerchant.CreatedAt, merchant.CreatedAt, time.Second)
	require.Equal(t, tMerchant.Owner, merchant.Owner)
}

// func TestUpdateMerchantWithOnlyAbout(t *testing.T) {
// 	tMerchant := createRandomMerchant(t)
// 	arg := UpdateMerchantParams{
// 		ID: tMerchant.ID,
// 	}

// 	merchant, err := testQueries.UpdateMerchant(context.Background(), arg)

// 	require.NoError(t, err)
// 	require.NotEmpty(t, merchant)
// 	require.Equal(t, arg.About, merchant.About)
// 	require.Equal(t, arg.ID, merchant.ID)
// 	require.Equal(t, tMerchant.Balance, merchant.Balance)
// 	require.Equal(t, tMerchant.ImageUrl, merchant.ImageUrl)
// 	require.Equal(t, tMerchant.Profession, merchant.Profession)
// 	require.Equal(t, tMerchant.Rating, merchant.Rating)
// 	require.Equal(t, tMerchant.Title, merchant.Title)
// 	require.WithinDuration(t, tMerchant.CreatedAt, merchant.CreatedAt, time.Second)
// 	require.Equal(t, tMerchant.Owner, merchant.Owner)
// }
