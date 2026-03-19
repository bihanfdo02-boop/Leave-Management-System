package domain

import "fmt"

// AppError represents application-level errors
type AppError struct {
	Code    string
	Message string
	Status  int
	Details interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Error codes
const (
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeInternalServer = "INTERNAL_SERVER_ERROR"
	ErrCodeBadRequest     = "BAD_REQUEST"
	ErrCodeDuplicate      = "DUPLICATE_ENTRY"
	ErrCodeInvalidInput   = "INVALID_INPUT"
)

// Common errors
var (
	ErrUserNotFound       = &AppError{Code: ErrCodeNotFound, Message: "user not found", Status: 404}
	ErrEmployeeNotFound   = &AppError{Code: ErrCodeNotFound, Message: "employee not found", Status: 404}
	ErrLeaveRequestNotFound = &AppError{Code: ErrCodeNotFound, Message: "leave request not found", Status: 404}
	ErrLeaveTypeNotFound  = &AppError{Code: ErrCodeNotFound, Message: "leave type not found", Status: 404}
	ErrDepartmentNotFound = &AppError{Code: ErrCodeNotFound, Message: "department not found", Status: 404}
	ErrLeaveBalanceNotFound = &AppError{Code: ErrCodeNotFound, Message: "leave balance not found", Status: 404}

	ErrUnauthorized      = &AppError{Code: ErrCodeUnauthorized, Message: "unauthorized", Status: 401}
	ErrForbidden         = &AppError{Code: ErrCodeForbidden, Message: "forbidden", Status: 403}
	ErrInvalidCredentials = &AppError{Code: ErrCodeUnauthorized, Message: "invalid email or password", Status: 401}

	ErrDuplicateEmail    = &AppError{Code: ErrCodeDuplicate, Message: "email already exists", Status: 409}
	ErrDepartmentExists  = &AppError{Code: ErrCodeConflict, Message: "department already exists", Status: 409}

	ErrInternalServer    = &AppError{Code: ErrCodeInternalServer, Message: "internal server error", Status: 500}
	ErrInvalidInput      = &AppError{Code: ErrCodeInvalidInput, Message: "invalid input", Status: 400}
	ErrInsufficientLeave = &AppError{Code: ErrCodeValidation, Message: "insufficient leave balance", Status: 400}
)

// NewAppError creates a new AppError
func NewAppError(code string, message string, status int, details interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Details: details,
	}
}