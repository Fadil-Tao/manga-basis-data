package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func handleSqlError(err error) error {
	if err == nil {
		return fmt.Errorf("empty")
	}
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		if mysqlErr.Number == 1644 {
			return fmt.Errorf("%s", trimErrorMessage( mysqlErr.Message))
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("not found")
	}

	return fmt.Errorf("error %s", err)
}


func trimErrorMessage(msg string)string{
	if strings.Contains(msg , ":") {
		result := strings.TrimLeft(msg, ":")
		return result
	}
	return msg
}