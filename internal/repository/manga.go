package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Fadil-Tao/manga-basis-data/internal/model"
)

type MangaRepo struct {
	DB *sql.DB
}

func NewMangaRepo(db *sql.DB) *MangaRepo {
	return &MangaRepo{
		DB: db,
	}
}

func (m *MangaRepo) CreateManga(ctx context.Context, manga *model.Manga, userId int) error {
	query := `call add_manga(?,?,?,?,?,?);`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, manga.Title, manga.Synopsys, manga.Manga_status, manga.Published_at, manga.Finished_at,userId)
	if err != nil {
		slog.Error("Error executing", "Err", err)
		return handleSqlError(err)
	}
	return nil
}

func (m *MangaRepo) ConnectMangaAuthor(ctx context.Context, obj *model.MangaAuthorPivot, userId int) error {
	query := `call connect_author_manga(?,?,?);`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "Error", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, obj.Manga_id, obj.Author_id, userId)
	if err != nil {
		slog.Error("Error in executing", "Err", err)
		return handleSqlError(err)
	}
	return nil
}

func (m *MangaRepo) ConnectMangaGenre(ctx context.Context, obj *model.MangaGenrePivot, userId int) error {
	query := `call connect_genre_manga(?,?,?);`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "Error", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, obj.Manga_id, obj.Genre_id, userId)
	if err != nil {
		slog.Error("Error in executing", "Err", err)
		return handleSqlError(err)
	}
	return nil
}

func (m *MangaRepo) GetMangaById(ctx context.Context, id string) (*model.Manga, error) {
	query := `call get_manga_detail(?);`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return nil, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, id)

	var Manga model.Manga

	err = result.Scan(&Manga.Id, &Manga.Title, &Manga.Synopsys, &Manga.Manga_status, &Manga.Published_at, &Manga.Finished_at)
	if err != nil {
		slog.Error("Error", "message", err)
		return nil, err
	}
	return &Manga, nil
}

func (m *MangaRepo) GetMangaAuthor(ctx context.Context, idManga string) ([]model.Author, error) {
	query := `call get_manga_author(?);`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, idManga)
	if err != nil {
		slog.Error("error not found", "error", err)
		return nil, err
	}
	defer rows.Close()

	var authors []model.Author
	for rows.Next() {
		var author model.Author

		if err := rows.Scan(&author.Id, &author.Name); err != nil {
			return authors, nil
		}
		authors = append(authors, author)
	}
	if err = rows.Err(); err != nil {
		return authors, err
	}
	return authors, nil
}

func (m *MangaRepo) GetMangaGenre(ctx context.Context, idManga string) ([]model.Genre, error) {
	query := `call get_manga_genre(?);`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, idManga)
	if err != nil {
		slog.Error("error not found", "error", err)
		return nil, err
	}
	defer rows.Close()

	var genres []model.Genre
	for rows.Next() {
		var genre model.Genre

		if err := rows.Scan(&genre.Id, &genre.Name); err != nil {
			return genres, nil
		}
		genres = append(genres, genre)
	}
	if err = rows.Err(); err != nil {
		return genres, err
	}
	return genres, nil
}

func (m *MangaRepo) GetAllMangaWithLimit(ctx context.Context, limit int) ([]*model.MangaList, error) {
	query := `call get_all_manga_with_limit(?);`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, limit)
	if err != nil {
		slog.Error("error not found", "error", err)
		return nil, err
	}
	defer rows.Close()

	var mangas []*model.MangaList
	for rows.Next() {
		var manga model.MangaList
		if err := rows.Scan(&manga.Id, &manga.Title, &manga.Synopsys, &manga.Manga_status, &manga.Published_at, &manga.Total_like, &manga.Rating); err != nil {
			return nil, err
		}
		mangas = append(mangas, &manga)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return mangas, nil
}

func (m *MangaRepo) SearchMangaByName(ctx context.Context, name string) ([]model.Manga, error) {
	query := `call get_manga_by_title(?);`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return nil, err 
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, name)
	if err != nil {
		slog.Error("error not found","error", err )
		return nil, err
	}
	defer rows.Close()
	var mangas []model.Manga
	for rows.Next(){
		var manga model.Manga
		if err := rows.Scan(&manga.Id, &manga.Title,  &manga.Synopsys, &manga.Manga_status,&manga.Published_at, &manga.Finished_at); err != nil{
			return nil,handleSqlError(err)
		}
		mangas = append(mangas, manga)
	}
	if len(mangas) == 0 {
		return nil,handleSqlError(err)
	}
	if err = rows.Err(); err != nil {
		return mangas, handleSqlError(err) 
	}
	return mangas, nil
}

func (m *MangaRepo) GetaAllMangaWithLimit(ctx context.Context, limit int) {

}

func (m *MangaRepo) DeleteMangaById(ctx context.Context, id string) error {
	return nil
}
