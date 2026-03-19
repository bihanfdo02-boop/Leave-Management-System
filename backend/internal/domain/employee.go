package domain

import (
	"time"

	"github.com/google/uuid"
)

// Employee represents an employee in the organization
type Employee struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;uniqueIndex"`
	FirstName      string    `json:"first_name" gorm:"type:varchar(100)"`
	LastName       string    `json:"last_name" gorm:"type:varchar(100)"`
	Email          string    `json:"email" gorm:"type:varchar(255);index"`
	Phone          string    `json:"phone" gorm:"type:varchar(20)"`
	DepartmentID   *uuid.UUID `json:"department_id,omitempty" gorm:"type:uuid;index"`
	ManagerID      *uuid.UUID `json:"manager_id,omitempty" gorm:"type:uuid;index"`
	HireDate       time.Time `json:"hire_date" gorm:"type:date"`
	IsActive       bool      `json:"is_active" gorm:"default:true;index"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	User          *User              `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Department    *Department        `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	Manager       *Employee          `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
	LeaveRequests []LeaveRequest     `json:"leave_requests,omitempty" gorm:"foreignKey:EmployeeID"`
	LeaveBalances []LeaveBalance     `json:"leave_balances,omitempty" gorm:"foreignKey:EmployeeID"`
	Approvals     []LeaveRequest     `json:"approvals,omitempty" gorm:"foreignKey:ApprovedByID"`
}

// TableName specifies the table name for Employee
func (Employee) TableName() string {
	return "employees"
}

// GetFullName returns the full name of the employee
func (e *Employee) GetFullName() string {
	return e.FirstName + " " + e.LastName
}
// package domain

// import (
// 	"time"

// 	"github.com/google/uuid"
// )

// // Employee represents an employee in the organization
// type Employee struct {
// 	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
// 	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;uniqueIndex"`
// 	FirstName      string    `json:"first_name" gorm:"type:varchar(100)"`
// 	LastName       string    `json:"last_name" gorm:"type:varchar(100)"`
// 	Email          string    `json:"email" gorm:"type:varchar(255);index"`
// 	Phone          string    `json:"phone" gorm:"type:varchar(20)"`
// 	DepartmentID   *uuid.UUID `json:"department_id,omitempty" gorm:"type:uuid;index"`
// 	ManagerID      *uuid.UUID `json:"manager_id,omitempty" gorm:"type:uuid;index"`
// 	HireDate       time.Time `json:"hire_date" gorm:"type:date"`
// 	IsActive       bool      `json:"is_active" gorm:"default:true;index"`
// 	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
// 	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

// 	// Relationships
// 	User          *User              `json:"user,omitempty" gorm:"foreignKey:UserID"`
// 	Department    *Department        `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
// 	Manager       *Employee          `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
// 	LeaveRequests []LeaveRequest     `json:"leave_requests,omitempty" gorm:"foreignKey:EmployeeID"`
// 	LeaveBalances []LeaveBalance     `json:"leave_balances,omitempty" gorm:"foreignKey:EmployeeID"`
// 	Approvals     []LeaveRequest     `json:"approvals,omitempty" gorm:"foreignKey:ApprovedByID"`
// }

// // TableName specifies the table name for Employee
// func (Employee) TableName() string {
// 	return "employees"
// }

// // GetFullName returns the full name of the employee
// func (e *Employee) GetFullName() string {
// 	return e.FirstName + " " + e.LastName
// }

// // package domain

// // import (
// // 	"time"

// // 	"github.com/google/uuid"
// // )

// // // Employee represents an employee in the system
// // type Employee struct {
// // 	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
// // 	UserID       uuid.UUID  `json:"user_id" gorm:"type:uuid;uniqueIndex;not null"`
// // 	FirstName    string     `json:"first_name" gorm:"type:varchar(100)"`
// // 	LastName     string     `json:"last_name" gorm:"type:varchar(100)"`
// // 	EmployeeID   string     `json:"employee_id" gorm:"type:varchar(50);uniqueIndex"`
// // 	Department   string     `json:"department" gorm:"type:varchar(100)"`
// // 	Position     string     `json:"position" gorm:"type:varchar(100)"`
// // 	ReportingTo  *uuid.UUID `json:"reporting_to,omitempty" gorm:"type:uuid"`
// // 	JoinDate     time.Time  `json:"join_date"`
// // 	LeaveBalance int        `json:"leave_balance" gorm:"default:20"`
// // 	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
// // 	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

// // 	// Relationships
// // 	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
// // }

// // // TableName specifies the table name for Employee
// // func (Employee) TableName() string {
// // 	return "employees"
// // }
