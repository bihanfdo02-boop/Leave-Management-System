package domain

import (
	"time"

	"github.com/google/uuid"
)

// Department represents a department in the organization
type Department struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;type:varchar(100)"`
	Description string    `json:"description" gorm:"type:text"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Employees []Employee `json:"employees,omitempty" gorm:"foreignKey:DepartmentID"`
}

// TableName specifies the table name for Department
func (Department) TableName() string {
	return "departments"
}
// package domain

// import (
// 	"time"

// 	"github.com/google/uuid"
// )

// // Department represents a department in the organization
// type Department struct {
// 	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
// 	Name        string    `json:"name" gorm:"uniqueIndex;type:varchar(100)"`
// 	Description string    `json:"description" gorm:"type:text"`
// 	IsActive    bool      `json:"is_active" gorm:"default:true"`
// 	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
// 	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

// 	// Relationships
// 	Employees []*Employee `json:"employees,omitempty" gorm:"foreignKey:DepartmentID"`
// }

// // TableName specifies the table name for Department
// func (Department) TableName() string {
// 	return "departments"
// }