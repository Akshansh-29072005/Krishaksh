package providers

import (
	"errors"
	"fmt"
	"time"
)

type RetryableError struct {
	Provider   string
	StatusCode int
	RetryAfter time.Duration
	Message    string
}

func (e *RetryableError) Error() string {
	if e == nil {
		return ""
	}
	if e.RetryAfter > 0 {
		return fmt.Sprintf("%s transient failure (status=%d, retry_after=%s): %s", e.Provider, e.StatusCode, e.RetryAfter, e.Message)
	}
	return fmt.Sprintf("%s transient failure (status=%d): %s", e.Provider, e.StatusCode, e.Message)
}

func (e *RetryableError) Unwrap() error {
	return nil
}

func AsRetryableError(err error) (*RetryableError, bool) {
	var retryableErr *RetryableError
	if errors.As(err, &retryableErr) {
		return retryableErr, true
	}
	return nil, false
}
