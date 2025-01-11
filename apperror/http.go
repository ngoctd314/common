package apperror

import "errors"

type HTTPError struct {
	// embed base error
	BaseError
	// The request ID
	//
	// The request ID is often exposed internally in order to trace
	// errors across service architectures. This is often a UUID.
	//
	// example: d7ef54b1-ec15-46e6-bccb-524b82c035e6
	RequestID string `json:"request_id,omitempty"`
	// Custom error code
	//
	// example: invalid_user
	ErrCode string `json:"code,omitempty"`
	// HTTP Code
	//
	// example: 400
	HTTPCode int `json:"-"`
	/// validator, cast
	Validation any `json:"validation,omitempty"`

	// Further error details
	Details any `json:"details,omitempty"`

	// Debug information
	//
	// This field is often not exposed to protect against leaking
	// sensitive information.
	//
	// example: SQL field "foo" is not a bool.
	debug any `json:"-"`
}

func NewHTTPError(err error, httpCode int) *HTTPError {
	if httpCode <= 0 && httpCode >= 600 {
		httpCode = 400
	}

	var baseErr *BaseError
	if errors.As(err, &baseErr) {
		return &HTTPError{
			BaseError: *baseErr,
			HTTPCode:  httpCode,
		}
	}

	return &HTTPError{
		BaseError: BaseError{
			ID:          errID(),
			message:     err.Error(),
			ancestorErr: err,
		},
		HTTPCode: httpCode,
	}
}

func (e *HTTPError) SetRequestID(rid string) *HTTPError {
	e.RequestID = rid
	return e
}

func (e *HTTPError) SetErrCode(errCode string) *HTTPError {
	e.ErrCode = errCode
	return e
}

// SetDetails override detail information
func (e *HTTPError) SetDetails(details any) *HTTPError {
	if details != nil {
		e.Details = details
	}
	return e
}

func (e *HTTPError) SetDebug(debugInfo any) *HTTPError {
	if debugInfo != nil {
		e.debug = debugInfo
	}
	return e
}
