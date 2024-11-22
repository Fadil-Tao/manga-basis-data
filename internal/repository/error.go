package repository

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func handleSqlError(err error) error {
	var mysqlErr *mysql.MySQLError

	if errors.As(err, &mysqlErr) {
		if mysqlErr.Number == 1644 {
			return fmt.Errorf("%s", mysqlErr.Message)
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("not found")
	}
	return fmt.Errorf("error : %s", err)
}
