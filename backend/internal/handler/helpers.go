package handler

import (
	"github.com/google/uuid"
)

// parseUserID parses user ID from context
func parseUserID(userIDStr interface{}) (uuid.UUID, error) {
	if str, ok := userIDStr.(string); ok {
		return uuid.Parse(str)
	}
	return uuid.UUID{}, nil
}