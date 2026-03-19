package dto

import (
	"github.com/KalinduBihan/leave-management-api/internal/domain"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

// ErrorResponse represents an error in the response
type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// SuccessResponse creates a success response
func SuccessResponse(message string, data interface{}) *APIResponse {
	return &APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponseFromAppError creates an error response from AppError
func ErrorResponseFromAppError(err *domain.AppError) *APIResponse {
	return &APIResponse{
		Success: false,
		Message: err.Message,
		Error: &ErrorResponse{
			Code:    err.Code,
			Message: err.Message,
			Details: err.Details,
		},
	}
}