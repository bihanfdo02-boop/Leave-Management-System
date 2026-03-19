package domain

import (
	"time"

	"github.com/google/uuid"
)

// LeaveStatus represents the status of a leave request
type LeaveStatus string

const (
	StatusPending  LeaveStatus = "pending"
	StatusApproved LeaveStatus = "approved"
	StatusRejected LeaveStatus = "rejected"
)

// LeaveRequest represents a request for leave by an employee
type LeaveRequest struct {
	ID               uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey"`
	EmployeeID       uuid.UUID   `json:"employee_id" gorm:"type:uuid;index"`
	LeaveTypeID      uuid.UUID   `json:"leave_type_id" gorm:"type:uuid"`
	StartDate        time.Time   `json:"start_date" gorm:"type:date;index"`
	EndDate          time.Time   `json:"end_date" gorm:"type:date;index"`
	NumberOfDays     float32     `json:"number_of_days" gorm:"type:decimal(5,2)"`
	Reason           string      `json:"reason" gorm:"type:text"`
	AttachmentURL    string      `json:"attachment_url" gorm:"type:varchar(500)"`
	Status           LeaveStatus `json:"status" gorm:"type:varchar(50);default:'pending';index"`
	ApprovedByID     *uuid.UUID  `json:"approved_by_id,omitempty" gorm:"type:uuid"`
	ApprovalDate     *time.Time  `json:"approval_date,omitempty"`
	RejectionReason  string      `json:"rejection_reason" gorm:"type:text"`
	CreatedAt        time.Time   `json:"created_at" gorm:"autoCreateTime;index"`
	UpdatedAt        time.Time   `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Employee   *Employee  `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	LeaveType  *LeaveType `json:"leave_type,omitempty" gorm:"foreignKey:LeaveTypeID"`
	ApprovedBy *Employee  `json:"approved_by,omitempty" gorm:"foreignKey:ApprovedByID"`
}

// TableName specifies the table name for LeaveRequest
func (LeaveRequest) TableName() string {
	return "leave_requests"
}

// IsApproved checks if the leave request is approved
func (lr *LeaveRequest) IsApproved() bool {
	return lr.Status == StatusApproved
}

// IsRejected checks if the leave request is rejected
func (lr *LeaveRequest) IsRejected() bool {
	return lr.Status == StatusRejected
}

// IsPending checks if the leave request is pending
func (lr *LeaveRequest) IsPending() bool {
	return lr.Status == StatusPending
}
// package domain

// import (
// 	"time"

// 	"github.com/google/uuid"
// )

// // LeaveStatus represents the status of a leave request
// type LeaveStatus string

// const (
// 	LeaveStatusPending  LeaveStatus = "pending"
// 	LeaveStatusApproved LeaveStatus = "approved"
// 	LeaveStatusRejected LeaveStatus = "rejected"
// 	LeaveStatusCanceled LeaveStatus = "canceled"
// )

// // LeaveType represents the type of leave
// type LeaveType string

// const (
// 	LeaveTypeAnnual     LeaveType = "annual"
// 	LeaveTypeSick       LeaveType = "sick"
// 	LeaveTypeCasual     LeaveType = "casual"
// 	LeaveTypeMaternity  LeaveType = "maternity"
// 	LeaveTypePaternity  LeaveType = "paternity"
// 	LeaveTypeUnpaid     LeaveType = "unpaid"
// 	LeaveTypeSpecial    LeaveType = "special"
// )

// // LeaveRequest represents a leave request from an employee
// type LeaveRequest struct {
// 	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
// 	EmployeeID    uuid.UUID  `json:"employee_id" gorm:"type:uuid;index;not null"`
// 	LeaveType     LeaveType  `json:"leave_type" gorm:"type:varchar(50);not null"`
// 	StartDate     time.Time  `json:"start_date" gorm:"type:date;not null"`
// 	EndDate       time.Time  `json:"end_date" gorm:"type:date;not null"`
// 	Reason        string     `json:"reason" gorm:"type:text"`
// 	Status        LeaveStatus `json:"status" gorm:"type:varchar(50);default:'pending';index"`
// 	ApprovedByID  *uuid.UUID `json:"approved_by_id,omitempty" gorm:"type:uuid"`
// 	ApprovalDate  *time.Time `json:"approval_date,omitempty"`
// 	RejectionNote string     `json:"rejection_note,omitempty" gorm:"type:text"`
// 	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
// 	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

// 	// Relationships
// 	Employee     *Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
// 	ApprovedBy   *Employee `json:"approved_by,omitempty" gorm:"foreignKey:ApprovedByID"`
// }

// // TableName specifies the table name for LeaveRequest
// func (LeaveRequest) TableName() string {
// 	return "leave_requests"
// }

// // GetDuration returns the number of days for the leave request
// func (lr *LeaveRequest) GetDuration() int {
// 	return int(lr.EndDate.Sub(lr.StartDate).Hours()/24) + 1
// }
