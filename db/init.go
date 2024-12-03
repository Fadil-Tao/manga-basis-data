package db

import (
	"database/sql"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
)



func InitDB() *sql.DB {
	// just for testing
	dsn := "mbd:mbd_123@/manga_basis_data"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		slog.Error("error", "error", err.Error())
		return nil
	}
	slog.Info("Succesfully connected to database")
	return db
}

// kalender
// galeri hape
