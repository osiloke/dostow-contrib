package api

import (
	"fmt"
	"strings"
)

// APIError ...
type APIError struct {
	Status  string `json:"code"`
	Message string `json:"msg"`
}

func (e APIError) Error() string {
	return strings.TrimSpace(fmt.Sprintf("%v %v", e.Status, e.Message))
}

// Empty ...
func (e APIError) Empty() bool {
	return false
}

// relevantError returns any non-nil http-related error (creating the request,
func relevantError(httpError error, apiError *APIError) error {
	if httpError != nil {
		return httpError
	} else if apiError.Message != "" {
		return apiError
	}
	return nil
}
