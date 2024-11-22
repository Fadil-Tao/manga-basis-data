package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Fadil-Tao/manga-basis-data/internal/model"
)

type AuthorRepository struct {
	DB *sql.DB
}

func NewAuthorRepo(db *sql.DB) *AuthorRepository {
	return &AuthorRepository{
		DB: db,
	}
}

func (a *AuthorRepository) CreateAuthor(ctx context.Context, author *model.Author, userId int) error {
	query := `call add_author(?,?,?,?);`
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, author.Name, author.Birthday, author.Biography, userId)
	if err != nil {
		slog.Error("Error Executing procedure", "err", err)
		return handleSqlError(err)
	}
	return nil
}

func (a *AuthorRepository) GetAllAuthor() ([]*model.Author, error) {
	query := `call get_all_author;`

	stmt, err := a.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		slog.Error("error", "message", err)
		return nil, err
	}
	defer rows.Close()

	var authors []*model.Author
	for rows.Next() {
		var author model.Author
		if err := rows.Scan(&author.Id, &author.Name, &author.Birthday); err != nil {
			return authors, nil
		}
		authors = append(authors, &author)
	}
	if err = rows.Err(); err != nil {
		return authors, handleSqlError(err)
	}
	return authors, nil
}

func (a *AuthorRepository) DeleteAuthorById(ctx context.Context, id int, userId int) error {
	query := `call delete_author(?,?);`

	stmt, err := a.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, id, userId)
	if err != nil {
		slog.Error("error", "message", err)
		return err
	}
	return nil
}

func (a *AuthorRepository) UpdateAuthor(ctx context.Context, id int, data *model.Author, userId int) error {
	query := `call update_author(?,?,?,?,?);`

	stmt, err := a.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, id, data.Name, data.Birthday, data.Biography, userId)
	if err != nil {
		slog.Error("error", "message", err)
		return err
	}
	return nil
}

func (a *AuthorRepository) GetAuthorById(ctx context.Context, id int) (*model.Author, error) {
	query := `call get_author_by_id(?);`

	stmt, err := a.DB.Prepare(query)
	if err != nil {
		slog.Error("Error Preparing Query", "message", err)
		return nil,err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx,id)
	var author model.Author

	err = row.Scan(&author.Id, &author.Name, &author.Birthday, &author.Biography)
	if err != nil{
		slog.Error("Error", "message" , err)
		return nil, handleSqlError(err)
	} 
	return &author, nil
}

func (a *AuthorRepository) SearchAuthor(ctx context.Context, input string) ([]model.Author,error) {
	query := `call get_author_by_name(?);`

	stmt, err := a.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return nil, err 
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, input)
	if err != nil {
		slog.Error("error not found","error", err )
		return nil, err
	}

	var authors []model.Author
	for rows.Next(){
		var author model.Author
		if err := rows.Scan(&author.Id,&author.Name,&author.Birthday, &author.Biography); err != nil{
			return authors, err 
		}
		authors = append(authors, author)
	}

	if len(authors) == 0{
		return nil,handleSqlError(sql.ErrNoRows)
	}
	if err = rows.Err(); err != nil{
		return authors, err
	}
	return authors, nil
}

func (a *AuthorRepository) GetAuthorManga(ctx context.Context, id int)([]model.Manga, error){
	query := `call get_author_manga(?);`

	stmt , err := a.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx,id)
	if err != nil{
		slog.Error("error", "error at querying rows", err)
		return nil, err 
	}
	defer rows.Close()
	var mangas []model.Manga
	for rows.Next(){
		var manga model.Manga
		if err := rows.Scan(&manga.Id,&manga.Title,&manga.Manga_status,&manga.Synopsys,&manga.Published_at); err != nil{
			return mangas, err
		}
		mangas = append(mangas, manga)
	}
	if err = rows.Err(); err != nil{
		return mangas, err 
	}
	return mangas, err
}