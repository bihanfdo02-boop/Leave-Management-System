package utils

import (
	"strconv"
	"strings"
)

// GetPaginationParams extracts and validates pagination parameters
func GetPaginationParams(page, pageSize string) (int, int) {
	pageNum := DefaultPage
	pageSizeNum := DefaultPageSize

	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageNum = p
	}

	if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 {
		if ps > MaxPageSize {
			ps = MaxPageSize
		}
		pageSizeNum = ps
	}

	return pageNum, pageSizeNum
}

// GetOffset calculates database offset for pagination
func GetOffset(page, pageSize int) int {
	if page <= 0 {
		page = 1
	}
	return (page - 1) * pageSize
}

// CalculateTotalPages calculates total pages from total count
func CalculateTotalPages(total int64, pageSize int) int {
	pages := total / int64(pageSize)
	if total%int64(pageSize) > 0 {
		pages++
	}
	return int(pages)
}

// TrimString removes leading and trailing whitespace
func TrimString(s string) string {
	return strings.TrimSpace(s)
}

// NormalizeEmail normalizes email to lowercase
func NormalizeEmail(email string) string {
	return strings.ToLower(TrimString(email))
}

// IsEmptyString checks if string is empty after trimming
func IsEmptyString(s string) bool {
	return len(TrimString(s)) == 0
}