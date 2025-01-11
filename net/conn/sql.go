package conn

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ngoctd314/common/env"
)

type sqlDriver string

const (
	MySQLDriver sqlDriver = "mysql"
)

// SQL creates a new sql.DB instance
// it already ping the database to make sure the connection can be established
func SQL(driver sqlDriver, dbName string) (*sql.DB, error) {
	db, err := sql.Open(string(driver), env.GetString(fmt.Sprintf("%s.%s.dsn", driver, dbName)))
	if err != nil {
		return nil, fmt.Errorf("%w, detail: %w", errInvalidDSN, err)
	}

	// if err := db.Ping(); err != nil {
	// 	return nil, err
	// }

	if err := setSQLConnectionPoll(db, driver, dbName); err != nil {
		return nil, err
	}

	return db, nil
}

var (
	errInvalidMaxOpenConns    = errors.New("invalid maxOpenConns")
	errInvalidMaxIdleConns    = errors.New("invalid maxIdleConns")
	errInvalidConnMaxLifetime = errors.New("invalid connMaxLifetime")
	errInvalidConnMaxIdleTime = errors.New("invalid connMaxIdleTime")
	errInvalidDSN             = errors.New("invalid DSN, want format user:passwd@tcp(ip:port)/dbName")
)

func setSQLConnectionPoll(db *sql.DB, driver sqlDriver, dbName string) error {
	query, _ := url.ParseQuery(env.GetString(fmt.Sprintf("%s.%s.cfg", driver, dbName)))

	var errGroup error
	maxOpenConns, err := strconv.Atoi(query.Get("maxOpenConns"))
	if err != nil {
		errGroup = errors.Join(errGroup, errInvalidMaxOpenConns)
	}
	db.SetMaxOpenConns(maxOpenConns)

	maxIdleConns, err := strconv.Atoi(query.Get("maxIdleConns"))
	if err != nil {
		errGroup = errors.Join(errGroup, errInvalidMaxIdleConns)
	}
	db.SetMaxIdleConns(maxIdleConns)

	connMaxLifetime, err := time.ParseDuration(query.Get("connMaxLifetime"))
	if err != nil {
		errGroup = errors.Join(errGroup, errInvalidConnMaxLifetime)
	}
	db.SetConnMaxLifetime(connMaxLifetime)

	connMaxIdleTime, err := time.ParseDuration(query.Get("connMaxIdleTime"))
	if err != nil {
		errGroup = errors.Join(errGroup, errInvalidConnMaxIdleTime)
	}
	db.SetConnMaxIdleTime(connMaxIdleTime)

	if errGroup != nil {
		return errGroup
	}

	return nil
}
