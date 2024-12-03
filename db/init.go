package db

import (
	"database/sql"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
)



func InitDB() *sql.DB {
	dbUser := map[string]string{
		"user": "mbd",
		"password": "mbd_123",
		"dbName" : "manga_basis_data",
	}
	// just for testing
	dsn :=  dbUser["user"] + ":" + dbUser["password"] + "@/" +dbUser["dbName"]
	// "mbd:mbd_123@/manga_basis_data"
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
