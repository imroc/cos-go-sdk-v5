package cos

import "fmt"

type Error struct {
	StatusCode int    // HTTP status code (200, 403, ...)
	Code       string // COS error code ("UnsupportedOperation", ...)
	Message    string // The human-oriented error message
	Resource   string
	RequestId  string
	TraceId    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("COS API Error: RequestId: %s Status Code: %d Code: %s Message: %s", e.RequestId, e.StatusCode, e.Code, e.Message)
}
