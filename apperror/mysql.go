package apperror

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// IsMySQLDuplicate check err is MySQL Duplicate or not
func IsMySQLDuplicate(err error) bool {
	var mysqlErr *mysql.MySQLError

	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}
	return false
}

// IsMySQLOutOfRange check err is MySQL out of range  or not
func IsMySQLOutOfRange(err error) bool {
	var mysqlErr *mysql.MySQLError

	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1690 {
		return true
	}
	return false
}

func GormTranslator(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}
	if IsMySQLDuplicate(err) {
		return Conflict
	}

	return err
}
