package apperror

import (
	"github.com/google/uuid"
)

func errID() string {
	return uuid.New().String()
}
