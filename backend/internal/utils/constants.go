package utils

import "time"

// Pagination constants
const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// Time constants
const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02T15:04:05Z07:00"
)

// JWT constants
const (
	ClaimsUserID   = "user_id"
	ClaimsRole     = "role"
	ClaimsEmail    = "email"
)

// Leave constants
const (
	HolidaysPerYear = 365
)

// Validation constants
const (
	MinPasswordLength = 8
	MaxPasswordLength = 256
	MaxNameLength     = 100
	MaxPhoneLength    = 20
	MaxReasonLength   = 5000
)

// Token expiration
var (
	TokenExpiration = 24 * time.Hour
	RefreshTokenExpiration = 7 * 24 * time.Hour
)