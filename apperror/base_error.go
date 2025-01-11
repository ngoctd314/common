package apperror

type BaseError struct {
	ID string `json:"id"`
	// example: The resource could not be found
	// required: true
	message string `json:"-"`
	// ancestor error
	ancestorErr error
}

// New return BaseError instance
func New(message string) *BaseError {
	return &BaseError{
		ID:      errID(),
		message: message,
	}
}

func (e *BaseError) Error() string {
	if e != nil {
		return e.message
	}
	return "internal server error, please contact the system administrator"
}

func (e *BaseError) Ancestor() error {
	return e.ancestorErr
}

// SetErr for debug, logging
func (e *BaseError) SetAncestor(err error) {
	if err != nil {
		e.ancestorErr = err
	}
}
