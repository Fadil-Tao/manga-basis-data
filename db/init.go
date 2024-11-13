package db

import (
	"database/sql"
	"log/slog"

	config "github.com/Fadil-Tao/manga-basis-data/configs"
	_ "github.com/go-sql-driver/mysql"
)



func InitDB(cfg *config.ConfDB) *sql.DB{
	dsn := cfg.Username + ":" + cfg.Password + "@" + "/" + cfg.DBName

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		slog.Error("error", "error",err.Error())
		return nil
	}

	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetMaxOpenConns(cfg.MaxOpen)
	slog.Info("Succesfully connected to database")
	return db
} 


// kalender
// galeri hape 