package fabman

import (
	"fmt"
)

//Error contain error information returned by FabMab API
type Error struct {
	StatusCode int    `json:"statusCode"`
	ErrorCode  string `json:"error"`
	Message    string `json:"message"`
}

// Error returns string representation of FabMan Error
func (e Error) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.StatusCode, e.ErrorCode, e.Message)
}
