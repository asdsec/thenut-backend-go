// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CommentType string

const (
	CommentTypePost     CommentType = "post"
	CommentTypeMerchant CommentType = "merchant"
)

func (e *CommentType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = CommentType(s)
	case string:
		*e = CommentType(s)
	default:
		return fmt.Errorf("unsupported scan type for CommentType: %T", src)
	}
	return nil
}

type NullCommentType struct {
	CommentType CommentType
	Valid       bool // Valid is true if CommentType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCommentType) Scan(value interface{}) error {
	if value == nil {
		ns.CommentType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.CommentType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCommentType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.CommentType, nil
}

type AppVersion struct {
	ID        int64     `json:"id"`
	Tag       string    `json:"tag"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID          int64       `json:"id"`
	CommentType CommentType `json:"comment_type"`
	// cannot be null if comment_type is post
	PostID sql.NullInt64 `json:"post_id"`
	// cannot be null if comment_type is merchant
	MerchantID sql.NullInt64 `json:"merchant_id"`
	Owner      string        `json:"owner"`
	Comment    string        `json:"comment"`
	CreatedAt  time.Time     `json:"created_at"`
}

type Consultancy struct {
	ID         int64 `json:"id"`
	MerchantID int64 `json:"merchant_id"`
	CustomerID int64 `json:"customer_id"`
	// must be positive
	Cost      int64     `json:"cost"`
	CreatedAt time.Time `json:"created_at"`
}

type Customer struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

type Merchant struct {
	ID         int64     `json:"id"`
	Owner      string    `json:"owner"`
	Balance    int64     `json:"balance"`
	Profession string    `json:"profession"`
	Title      string    `json:"title"`
	About      string    `json:"about"`
	ImageUrl   string    `json:"image_url"`
	Rating     float64   `json:"rating"`
	CreatedAt  time.Time `json:"created_at"`
}

type Post struct {
	ID         int64 `json:"id"`
	MerchantID int64 `json:"merchant_id"`
	// can be null only if image_url is not null
	Title sql.NullString `json:"title"`
	// can be null only if title is not null
	ImageUrl  sql.NullString `json:"image_url"`
	Likes     int32          `json:"likes"`
	CreatedAt time.Time      `json:"created_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PhoneNumber       string    `json:"phone_number"`
	ImageUrl          string    `json:"image_url"`
	Gender            string    `json:"gender"`
	Disabled          bool      `json:"disabled"`
	BirthDate         time.Time `json:"birth_date"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}
