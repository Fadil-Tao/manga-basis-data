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

func (a *AuthorRepository)CreateAuthor(ctx context.Context, author *model.Author) error{
	query := `call add_author(?,?,?);`
	stmt, err := a.DB.Prepare(query)
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

func (a *AuthorRepository) GetAllAuthor()([]*model.Author, error){
	query := `call get_all_author;`

	stmt, err := a.DB.Prepare(query)
	if err != nil{
		slog.Error("error preparing statement", "error", err)
		return nil,err
	}
	defer stmt.Close()
	
	rows, err := stmt.Query()
	if err != nil{
		slog.Error("error", "message", err)
		return nil,err
	}
	defer rows.Close()
	
	var authors []*model.Author
	for rows.Next(){
		var author model.Author
		if err := rows.Scan(&author.Id, &author.Name,&author.Birthday); err != nil {
			return authors, nil
		}
		authors = append(authors, &author)
	}
	if err = rows.Err(); err != nil {
		return authors, err
	}
	return authors, nil
}

func (a *AuthorRepository) DeleteAuthorById(ctx context.Context, id int)error{
	query := `call delete_author(?);`

	stmt,err := a.DB.Prepare(query)
	if err != nil{
		slog.Error("error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()

	_ , err = stmt.ExecContext(ctx,id)
	if err != nil{
		slog.Error("error","message",err)
		return err
	}
	return nil
}

func (a *AuthorRepository) UpdateAuthor(ctx context.Context, id int,data *model.Author )error{
	query := `call update_author(?,?,?,?);`

	stmt,err := a.DB.Prepare(query)
	if err != nil{
		slog.Error("error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()
	_ , err = stmt.ExecContext(ctx,id,data.Name,data.Birthday,data.Biography)
	if err != nil{
		slog.Error("error","message",err)
		return err
	}
	return nil
}