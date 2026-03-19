package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleManager  UserRole = "manager"
	RoleEmployee UserRole = "employee"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;type:varchar(255)"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255)"`
	Role         UserRole  `json:"role" gorm:"type:varchar(50);default:'employee'"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Employee *Employee `json:"employee,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}
// package domain

// import (
// 	"time"

// 	"github.com/google/uuid"
// )

// // UserRole represents the role of a user in the system
// type UserRole string

// const (
// 	RoleAdmin    UserRole = "admin"
// 	RoleManager  UserRole = "manager"
// 	RoleEmployee UserRole = "employee"
// )

// // User represents a user in the system
// type User struct {
// 	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
// 	Email        string    `json:"email" gorm:"uniqueIndex;type:varchar(255)"`
// 	PasswordHash string    `json:"-" gorm:"type:varchar(255)"`
// 	Role         UserRole  `json:"role" gorm:"type:varchar(50);default:'employee'"`
// 	IsActive     bool      `json:"is_active" gorm:"default:true"`
// 	LastLogin    *time.Time `json:"last_login,omitempty"`
// 	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
// 	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`

// 	// Relationships
// 	Employee *Employee `json:"employee,omitempty" gorm:"foreignKey:UserID"`
// }

// // TableName specifies the table name for User
// func (User) TableName() string {
// 	return "users"
// }