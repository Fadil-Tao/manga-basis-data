package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Fadil-Tao/manga-basis-data/internal/model"
)


type AuthorRepository struct{
	DB *sql.DB
}

func NewAuthorRepo(db *sql.DB) *AuthorRepository{
	return &AuthorRepository{
		DB: db,
	}
}


func (u *AuthorRepository)  CreateAuthor(ctx context.Context, author *model.Author) error{
	query := `call add_author(?,?,?) `

	stmt, err := u.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return err
	}

	defer stmt.Close()


	_ , err = stmt.ExecContext(ctx,author.Name, author.Birthday, author.Biography)
	if err != nil {
		slog.Error("Error Executing procedure", "err", err)
		return err
	}
	return nil
}