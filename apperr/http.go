package apperr

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
	// example: invalid_user
	ErrType    string `json:"type,omitempty"`
	HTTPCode   int    `json:"-"`
	Validation any    `json:"validation,omitempty"`

	// Further error details
	// Details any `json:"details,omitempty"`

	// Debug information
	//
	// This field is often not exposed to protect against leaking
	// sensitive information.
	//
	// example: SQL field "foo" is not a bool.
	// debug any `json:"-"`
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
			ancestorErr: err,
		},
		HTTPCode: httpCode,
	}
}

func (e *HTTPError) SetRequestID(rid string) *HTTPError {
	e.RequestID = rid
	return e
}

func (e *HTTPError) SetErrType(errCode string) *HTTPError {
	e.ErrType = errCode
	return e
}

// SetDetails override detail information
// func (e *HTTPError) SetDetails(details any) *HTTPError {
// 	if details != nil {
// 		e.Details = details
// 	}
// 	return e
// }

// func (e *HTTPError) SetDebug(debugInfo any) *HTTPError {
// 	if debugInfo != nil {
// 		e.debug = debugInfo
// 	}
// 	return e
// }
