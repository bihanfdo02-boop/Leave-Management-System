package domain

import (
	"time"

	"github.com/google/uuid"
)

// LeaveType represents different types of leave available
type LeaveType struct {
	ID                 uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name               string    `json:"name" gorm:"uniqueIndex;type:varchar(50)"`
	Description        string    `json:"description" gorm:"type:text"`
	DefaultDaysPerYear int       `json:"default_days_per_year" gorm:"default:20"`
	IsPaid             bool      `json:"is_paid" gorm:"default:true"`
	RequiresApproval   bool      `json:"requires_approval" gorm:"default:true"`
	IsActive           bool      `json:"is_active" gorm:"default:true;index"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	LeaveBalances []LeaveBalance `json:"leave_balances,omitempty" gorm:"foreignKey:LeaveTypeID"`
	LeaveRequests []LeaveRequest `json:"leave_requests,omitempty" gorm:"foreignKey:LeaveTypeID"`
}

// TableName specifies the table name for LeaveType
func (LeaveType) TableName() string {
	return "leave_types"
}

// package domain

// import (
// 	"time"

// 	"github.com/google/uuid"
// )

// // LeaveType represents different types of leave available
// type LeaveType struct {
// 	ID                 uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
// 	Name               string    `json:"name" gorm:"uniqueIndex;type:varchar(50)"`
// 	Description        string    `json:"description" gorm:"type:text"`
// 	DefaultDaysPerYear int       `json:"default_days_per_year" gorm:"default:20"`
// 	IsPaid             bool      `json:"is_paid" gorm:"default:true"`
// 	RequiresApproval   bool      `json:"requires_approval" gorm:"default:true"`
// 	IsActive           bool      `json:"is_active" gorm:"default:true;index"`
// 	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
// 	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`

// 	// Relationships
// 	LeaveBalances []LeaveBalance `json:"leave_balances,omitempty" gorm:"foreignKey:LeaveTypeID"`
// 	LeaveRequests []LeaveRequest `json:"leave_requests,omitempty" gorm:"foreignKey:LeaveTypeID"`
// }

// // TableName specifies the table name for LeaveType
// func (LeaveType) TableName() string {
// 	return "leave_types"
// }