package qb

import (
	"strings"

	"gorm.io/gorm"
)

type Builder interface {
	Build(tx *gorm.DB) *gorm.DB
}

func notEmptyString(s string) bool {
	return strings.TrimSpace(s) != ""
}

// Noop represent no-operation builder
// it does nothing
func Noop() Builder {
	return nil
}
