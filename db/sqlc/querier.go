// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AddMerchantBalance(ctx context.Context, arg AddMerchantBalanceParams) (Merchant, error)
	CreateAppVersion(ctx context.Context, arg CreateAppVersionParams) (AppVersion, error)
	CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error)
	CreateConsultancy(ctx context.Context, arg CreateConsultancyParams) (Consultancy, error)
	CreateCustomer(ctx context.Context, owner string) (Customer, error)
	CreateMerchant(ctx context.Context, arg CreateMerchantParams) (Merchant, error)
	CreatePost(ctx context.Context, arg CreatePostParams) (Post, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteComment(ctx context.Context, id int64) error
	DeleteCustomer(ctx context.Context, id int64) error
	DeleteMerchant(ctx context.Context, id int64) error
	DeletePost(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, username string) error
	GetAppVersion(ctx context.Context, tag string) (AppVersion, error)
	GetComment(ctx context.Context, id int64) (Comment, error)
	GetConsultancy(ctx context.Context, id int64) (Consultancy, error)
	GetCustomer(ctx context.Context, id int64) (Customer, error)
	GetMerchant(ctx context.Context, id int64) (Merchant, error)
	GetPost(ctx context.Context, id int64) (Post, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, username string) (User, error)
	ListAppVersions(ctx context.Context) ([]AppVersion, error)
	ListConsultancies(ctx context.Context, arg ListConsultanciesParams) ([]Consultancy, error)
	ListMerchantComments(ctx context.Context, arg ListMerchantCommentsParams) ([]Comment, error)
	ListMerchantPosts(ctx context.Context, arg ListMerchantPostsParams) ([]Post, error)
	ListMerchants(ctx context.Context, arg ListMerchantsParams) ([]Merchant, error)
	ListPostComments(ctx context.Context, arg ListPostCommentsParams) ([]Comment, error)
	ListPosts(ctx context.Context, arg ListPostsParams) ([]Post, error)
	UpdateAppVersion(ctx context.Context, arg UpdateAppVersionParams) (AppVersion, error)
	UpdateCustomer(ctx context.Context, arg UpdateCustomerParams) (Customer, error)
	UpdateEmail(ctx context.Context, arg UpdateEmailParams) (User, error)
	UpdateMerchant(ctx context.Context, arg UpdateMerchantParams) (Merchant, error)
	UpdatePassword(ctx context.Context, arg UpdatePasswordParams) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
