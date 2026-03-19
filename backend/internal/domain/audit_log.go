package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AuditLog represents an audit log entry for tracking changes
type AuditLog struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey"`
	UserID    *uuid.UUID      `json:"user_id,omitempty" gorm:"type:uuid;index"`
	Action    string          `json:"action" gorm:"type:varchar(100)"`
	EntityType string          `json:"entity_type" gorm:"type:varchar(50);index"`
	EntityID  *uuid.UUID      `json:"entity_id,omitempty" gorm:"type:uuid"`
	OldValues json.RawMessage `json:"old_values" gorm:"type:jsonb"`
	NewValues json.RawMessage `json:"new_values" gorm:"type:jsonb"`
	IPAddress string          `json:"ip_address" gorm:"type:varchar(45)"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime;index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for AuditLog
func (AuditLog) TableName() string {
	return "audit_logs"
}

// Scan implements sql.Scanner interface
func (al *AuditLog) Scan(value interface{}) error {
	bytes, _ := value.([]byte)
	return json.Unmarshal(bytes, &al)
}

// Value implements driver.Valuer interface
func (al AuditLog) Value() (driver.Value, error) {
	return json.Marshal(al)
}