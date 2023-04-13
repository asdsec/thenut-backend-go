package db

import (
	"context"
	"database/sql"
	"sort"
	"testing"
	"time"

	"github.com/asdsec/thenut/utils"
	"github.com/stretchr/testify/require"
)

func createRandomPost(t *testing.T, merchant Merchant) Post {
	arg := CreatePostParams{
		MerchantID: merchant.ID,
		Title: sql.NullString{
			String: utils.RandomString(6),
			Valid:  true,
		},
		ImageUrl: sql.NullString{
			String: utils.RandomImageUrl(),
			Valid:  true,
		},
	}

	post, err := testQueries.CreatePost(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, post)
	require.Equal(t, arg.ImageUrl, post.ImageUrl)
	require.Equal(t, arg.MerchantID, post.MerchantID)
	require.Equal(t, arg.Title, post.Title)
	require.NotEmpty(t, post.CreatedAt)
	require.NotEmpty(t, post.ID)
	require.Equal(t, int32(0), post.Likes)

	return post
}

func TestCreatePost(t *testing.T) {
	merchant := createRandomMerchant(t)
	arg := CreatePostParams{
		MerchantID: merchant.ID,
		Title: sql.NullString{
			String: utils.RandomString(6),
			Valid:  true,
		},
		ImageUrl: sql.NullString{
			String: utils.RandomImageUrl(),
			Valid:  true,
		},
	}

	post, err := testQueries.CreatePost(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, post)
	require.Equal(t, arg.ImageUrl, post.ImageUrl)
	require.Equal(t, arg.MerchantID, post.MerchantID)
	require.Equal(t, arg.Title, post.Title)
	require.NotEmpty(t, post.CreatedAt)
	require.NotEmpty(t, post.ID)
	require.Equal(t, int32(0), post.Likes)
}

func TestDeletePost(t *testing.T) {
	merchant := createRandomMerchant(t)
	post := createRandomPost(t, merchant)

	err := testQueries.DeletePost(context.Background(), post.ID)
	require.NoError(t, err)
	post, err = testQueries.GetPost(context.Background(), post.ID)
	require.Error(t, sql.ErrNoRows, err)
	require.Empty(t, post)
}

func TestGetPost(t *testing.T) {
	merchant := createRandomMerchant(t)
	expected := createRandomPost(t, merchant)

	post, err := testQueries.GetPost(context.Background(), expected.ID)
	require.NoError(t, err)
	require.NotEmpty(t, post)
	require.Equal(t, expected, post)
}

func TestListMerchantPosts(t *testing.T) {
	var expected []Post
	merchant := createRandomMerchant(t)
	for i := 0; i < 10; i++ {
		post := createRandomPost(t, merchant)
		expected = append(expected, post)
	}
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].ID > expected[j].ID
	})

	arg := ListMerchantPostsParams{
		MerchantID: merchant.ID,
		Limit:      10,
		Offset:     0,
	}

	posts, err := testQueries.ListMerchantPosts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, posts)
	require.Equal(t, len(expected), len(posts))
	for ie, e := range expected {
		for ip, p := range posts {
			if ie == ip {
				require.Equal(t, e.ID, p.ID)
				require.WithinDuration(t, e.CreatedAt, p.CreatedAt, time.Second)
				require.Equal(t, e.ImageUrl, p.ImageUrl)
				require.Equal(t, e.Likes, p.Likes)
				require.Equal(t, e.MerchantID, p.MerchantID)
				require.Equal(t, e.Title, p.Title)
			}
		}
	}
}

func TestListPosts(t *testing.T) {
	var expected []Post
	merchant := createRandomMerchant(t)
	for i := 0; i < 10; i++ {
		post := createRandomPost(t, merchant)
		expected = append(expected, post)
	}
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].CreatedAt.After(expected[j].CreatedAt)
	})

	arg := ListPostsParams{
		Limit:  10,
		Offset: 0,
	}

	posts, err := testQueries.ListPosts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, posts)
	require.Equal(t, len(expected), len(posts))
	for ie, e := range expected {
		for ip, p := range posts {
			if ie == ip {
				require.Equal(t, e.ID, p.ID)
				require.WithinDuration(t, e.CreatedAt, p.CreatedAt, time.Second)
				require.Equal(t, e.ImageUrl, p.ImageUrl)
				require.Equal(t, e.Likes, p.Likes)
				require.Equal(t, e.MerchantID, p.MerchantID)
				require.Equal(t, e.Title, p.Title)
			}
		}
	}
}
