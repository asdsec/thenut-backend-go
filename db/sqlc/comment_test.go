package db

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/asdsec/thenut/utils"
	"github.com/stretchr/testify/require"
)

func createRandomComment(t *testing.T, post Post, merchant Merchant, user User) Comment {
	arg := CreateCommentParams{
		CommentType: CommentTypeMerchant,
		MerchantID: sql.NullInt64{
			Int64: merchant.ID,
			Valid: merchant.ID != 0,
		},
		PostID: sql.NullInt64{
			Int64: post.ID,
			Valid: post.ID != 0,
		},
		Owner:   user.Username,
		Comment: utils.RandomString(12),
	}

	comment, err := testQueries.CreateComment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, comment)
	require.Equal(t, arg.Comment, comment.Comment)
	require.Equal(t, arg.CommentType, comment.CommentType)
	require.Equal(t, arg.MerchantID, comment.MerchantID)
	require.Equal(t, arg.Owner, comment.Owner)
	require.Equal(t, arg.PostID, comment.PostID)

	return comment
}

func TestCreateComment(t *testing.T) {
	merchant := createRandomMerchant(t)
	user := createRandomUser(t)
	post := createRandomPost(t, merchant)

	arg := CreateCommentParams{
		CommentType: CommentTypeMerchant,
		PostID: sql.NullInt64{
			Int64: post.ID,
			Valid: true,
		},
		Owner:   user.Username,
		Comment: utils.RandomString(12),
	}

	comment, err := testQueries.CreateComment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, comment)
	require.Equal(t, arg.Comment, comment.Comment)
	require.Equal(t, arg.CommentType, comment.CommentType)
	require.Equal(t, arg.MerchantID, comment.MerchantID)
	require.Equal(t, arg.Owner, comment.Owner)
	require.Equal(t, arg.PostID, comment.PostID)
}

func TestDeleteComment(t *testing.T) {
	user := createRandomUser(t)
	merchant := createRandomMerchant(t)
	post := createRandomPost(t, merchant)
	comment := createRandomComment(t, post, Merchant{}, user)

	err := testQueries.DeleteComment(context.Background(), comment.ID)
	require.NoError(t, err)
	comment, err = testQueries.GetComment(context.Background(), comment.ID)
	require.Error(t, err)
	require.Empty(t, comment)
}

func TestGetComment(t *testing.T) {
	user := createRandomUser(t)
	merchant := createRandomMerchant(t)
	post := createRandomPost(t, merchant)
	expected := createRandomComment(t, post, Merchant{}, user)

	comment, err := testQueries.GetComment(context.Background(), expected.ID)
	require.NoError(t, err)
	require.NotEmpty(t, comment)
	require.Equal(t, expected, comment)
}

func TestListMerchantComments(t *testing.T) {
	var expected []Comment
	merchant := createRandomMerchant(t)
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		comment := createRandomComment(t, Post{}, merchant, user)
		expected = append(expected, comment)
	}

	arg := ListMerchantCommentsParams{
		MerchantID: sql.NullInt64{
			Int64: merchant.ID,
			Valid: true,
		},
		Limit:  10,
		Offset: 0,
	}

	comments, err := testQueries.ListMerchantComments(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, comments)
	require.Equal(t, len(expected), len(comments))
	require.True(t, reflect.DeepEqual(comments, expected))
}

func TestListPostComments(t *testing.T) {
	var expected []Comment
	merchant := createRandomMerchant(t)
	post := createRandomPost(t, merchant)
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		comment := createRandomComment(t, post, Merchant{}, user)
		expected = append(expected, comment)
	}

	arg := ListPostCommentsParams{
		PostID: sql.NullInt64{
			Int64: post.ID,
			Valid: true,
		},
		Limit:  10,
		Offset: 0,
	}

	comments, err := testQueries.ListPostComments(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, comments)
	require.Equal(t, len(expected), len(comments))
	require.True(t, reflect.DeepEqual(comments, expected))
}
