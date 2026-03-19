package domain

import (
	"time"

	"github.com/google/uuid"
)

// LeaveBalance represents the leave balance for an employee for a specific leave type and year
type LeaveBalance struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	EmployeeID     uuid.UUID `json:"employee_id" gorm:"type:uuid;uniqueIndex:idx_leave_balance"`
	LeaveTypeID    uuid.UUID `json:"leave_type_id" gorm:"type:uuid;uniqueIndex:idx_leave_balance"`
	TotalDays      int       `json:"total_days"`
	UsedDays       int       `json:"used_days" gorm:"default:0"`
	RemainingDays  int       `json:"remaining_days"`
	Year           int       `json:"year" gorm:"uniqueIndex:idx_leave_balance"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Employee  *Employee  `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	LeaveType *LeaveType `json:"leave_type,omitempty" gorm:"foreignKey:LeaveTypeID"`
}

// TableName specifies the table name for LeaveBalance
func (LeaveBalance) TableName() string {
	return "leave_balances"
}
// package domain

// import (
// 	"time"

// 	"github.com/google/uuid"
// )

// // LeaveBalance represents the leave balance for an employee for a specific leave type
// type LeaveBalance struct {
// 	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
// 	EmployeeID   uuid.UUID `json:"employee_id" gorm:"type:uuid;index;not null"`
// 	LeaveType    LeaveType `json:"leave_type" gorm:"type:varchar(50);not null"`
// 	Year         int       `json:"year" gorm:"column:year;not null"`
// 	TotalDays    int       `json:"total_days" gorm:"not null"`
// 	UsedDays     int       `json:"used_days" gorm:"default:0"`
// 	PendingDays  int       `json:"pending_days" gorm:"default:0"`
// 	RemainingDays int      `json:"remaining_days" gorm:"not null"`
// 	CarriedOver  int       `json:"carried_over" gorm:"default:0"`
// 	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
// 	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`

// 	// Relationships
// 	Employee *Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
// }

// // TableName specifies the table name for LeaveBalance
// func (LeaveBalance) TableName() string {
// 	return "leave_balances"
// }

// // UpdateRemainingDays recalculates the remaining days based on used and pending days
// func (lb *LeaveBalance) UpdateRemainingDays() {
// 	lb.RemainingDays = lb.TotalDays + lb.CarriedOver - lb.UsedDays - lb.PendingDays
// }
