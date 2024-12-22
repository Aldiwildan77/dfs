package main

import (
	"errors"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const (
	MysqlDuplicateKeyErrorNumber = 1062
)

func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == MysqlDuplicateKeyErrorNumber {
		return true
	}

	return false
}
