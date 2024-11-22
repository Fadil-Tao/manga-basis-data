package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Fadil-Tao/manga-basis-data/internal/model"
)

type GenreRepo struct {
	DB *sql.DB
}

func NewGenrerepo(db *sql.DB) *GenreRepo {
	return &GenreRepo{
		DB: db,
	}
}

func (g *GenreRepo) CreateGenre(ctx context.Context, genre *model.Genre,userId int) error {
	query := `call add_genre(?,?,?);`

	stmt, err := g.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, genre.Name, genre.Description, userId)
	if err != nil {
		slog.Error("Error executing procedure", "Message", err)
		return handleSqlError(err)
	}
	return nil
}

func (g *GenreRepo) GetAllGenre() ([]*model.Genre, error) {
	query := `call get_all_genre();`

	stmt, err := g.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		slog.Error("error", "message", err)
		return nil, handleSqlError(err)
	}
	defer rows.Close()

	var genres []*model.Genre
	for rows.Next() {
		var genre model.Genre

		if err := rows.Scan(&genre.Id, &genre.Name, &genre.Description); err != nil {
			return genres, nil
		}
		genres = append(genres, &genre)
	}
	if err = rows.Err(); err != nil {
		return genres, handleSqlError(err )
	}
	return genres, nil
}

func (g *GenreRepo) DeleteGenreById(ctx context.Context, id int, userId int) error {
	query := `call delete_genre(?,?);`

	stmt, err := g.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id, userId)
	if err != nil {
		slog.Error("error", "message", err)
		return handleSqlError(err )
	}
	return nil
}

func (g *GenreRepo) UpdateGenre(ctx context.Context, id int, data *model.Genre, userId int) error {
	query := `call update_genre(?,?,?,?);`

	stmt, err := g.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, id, data.Name, data.Description, userId)
	if err != nil {
		slog.Error("error", "message", err)
		return handleSqlError(err)
	}
	return nil
}