package conn

import (
	"errors"
	"os"
	"testing"
)

func TestSQL(t *testing.T) {
	type args struct {
		driver sqlDriver
		dbName string
	}

	type hook struct {
		before func()
		after  func()
	}

	testCases := []struct {
		name    string
		args    args
		hook    hook
		wantErr error
	}{
		{
			name: "test SQL success",
			args: args{
				driver: MySQLDriver,
				dbName: "test",
			},
			hook: hook{
				before: func() {
					os.Setenv("MYSQL_TEST_DSN", "user:passwd@tcp(localhost:3306)/dbName?charset=utf8mb4&loc=Local&parseTime=True")
					os.Setenv("MYSQL_TEST_CFG", "maxOpenConns=10&maxIdleConns=10&connMaxLifetime=10m&connMaxIdleTime=5m")
				},
				after: func() {
					os.Setenv("MYSQL_TEST_DSN", "")
					os.Setenv("MYSQL_TEST_CFG", "")
				},
			},
		},
		{
			name: "test SQL invalid DSN",
			args: args{
				driver: MySQLDriver,
				dbName: "test",
			},
			hook: hook{
				before: func() {
					os.Setenv("MYSQL_TEST_DSN", "localhost:3306/dbName?charset=utf8mb4&loc=Local&parseTime=True")
				},
				after: func() {
					os.Setenv("MYSQL_TEST_DSN", "")
				},
			},
			wantErr: errInvalidDSN,
		},
		{
			name: "test SQL invalid maxOpenConns",
			args: args{
				driver: MySQLDriver,
				dbName: "test",
			},
			hook: hook{
				before: func() {
					os.Setenv("MYSQL_TEST_DSN", "user:passwd@tcp(localhost:3306)/dbName?charset=utf8mb4&loc=Local&parseTime=True")
					os.Setenv("MYSQL_TEST_CFG", "maxIdleConns=10&connMaxLifetime=10m&connMaxIdleTime=5m")
				},
				after: func() {
					os.Setenv("MYSQL_TEST_DSN", "")
					os.Setenv("MYSQL_TEST_CFG", "")
				},
			},
			wantErr: errInvalidMaxOpenConns,
		},
		{
			name: "test SQL invalid maxIdleConns",
			args: args{
				driver: MySQLDriver,
				dbName: "test",
			},
			hook: hook{
				before: func() {
					os.Setenv("MYSQL_TEST_DSN", "user:passwd@tcp(localhost:3306)/dbName?charset=utf8mb4&loc=Local&parseTime=True")
					os.Setenv("MYSQL_TEST_CFG", "maxOpenConns=10&connMaxLifetime=10m&connMaxIdleTime=5m")
				},
				after: func() {
					os.Setenv("MYSQL_TEST_DSN", "")
					os.Setenv("MYSQL_TEST_CFG", "")
				},
			},
			wantErr: errInvalidMaxIdleConns,
		},
		{
			name: "test SQL invalid connMaxLifetime",
			args: args{
				driver: MySQLDriver,
				dbName: "test",
			},
			hook: hook{
				before: func() {
					os.Setenv("MYSQL_TEST_DSN", "user:passwd@tcp(localhost:3306)/dbName?charset=utf8mb4&loc=Local&parseTime=True")
					os.Setenv("MYSQL_TEST_CFG", "maxOpenConns=10&maxIdleConns=10&connMaxIdleTime=5m")
				},
				after: func() {
					os.Setenv("MYSQL_TEST_DSN", "")
					os.Setenv("MYSQL_TEST_CFG", "")
				},
			},
			wantErr: errInvalidConnMaxLifetime,
		},
		{
			name: "test SQL invalid connMaxIdleTime",
			args: args{
				driver: MySQLDriver,
				dbName: "test",
			},
			hook: hook{
				before: func() {
					os.Setenv("MYSQL_TEST_DSN", "user:passwd@tcp(localhost:3306)/dbName?charset=utf8mb4&loc=Local&parseTime=True")
					os.Setenv("MYSQL_TEST_CFG", "maxOpenConns=10&maxIdleConns=10&connMaxLifetime=10m")
				},
				after: func() {
					os.Setenv("MYSQL_TEST_DSN", "")
					os.Setenv("MYSQL_TEST_CFG", "")
				},
			},
			wantErr: errInvalidConnMaxIdleTime,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.hook.before != nil {
				tc.hook.before()
			}
			if tc.hook.after != nil {
				defer tc.hook.after()
			}

			_, err := SQL(tc.args.driver, tc.args.dbName)
			if (err != nil || tc.wantErr != nil) && !errors.Is(err, tc.wantErr) {
				t.Fatalf("SQL() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
