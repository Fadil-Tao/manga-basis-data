package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Fadil-Tao/manga-basis-data/internal/model"
)

type ReadlistRepo struct {
	DB *sql.DB
}

func NewReadlistRepo(db *sql.DB)*ReadlistRepo{
	return &ReadlistRepo{
		DB: db,
	}
}

func (rl *ReadlistRepo) CreateReadlist(ctx context.Context, userId int,readlist *model.NewReadlistRequest)error{
	query := `call add_readlist(?,?,?)`

	stmt , err := rl.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()
	
	_ , err = stmt.ExecContext(ctx, userId, readlist.Name, readlist.Description)
	if err != nil {
		slog.Error("error executing procedure", "error", err)
		return handleSqlError(err)
	}

	return nil
}

func (rl *ReadlistRepo) DeleteReadlist(ctx context.Context, id int, userId int)error{
	query := `call delete_readlist(?,?)`
	
	stmt, err := rl.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()
	
	_ , err = stmt.ExecContext(ctx, id, userId)
	if err != nil{
		slog.Error("error executing procedure","error", err)
		return handleSqlError(err)
	}
	return nil
}


func (rl *ReadlistRepo) UpdateReadlist(ctx context.Context, id int, userId int,readlist *model.NewReadlistRequest) error{
	query := `call update_readlist(?,?,?,?)`

	stmt , err := rl.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()
	
	_ , err = stmt.ExecContext(ctx, readlist.Name, readlist.Description, id,userId)
	if err != nil {
		slog.Error("error executing procedure", "error", err)
		return handleSqlError(err)
	}

	return nil
}

func (rl *ReadlistRepo) SearchReadlist(ctx context.Context, name string) ([]*model.Readlist, error){
	query := `call search_readlist(?)`

	stmt, err := rl.DB.Prepare(query)
	if err != nil {
		slog.Error("error at preparing statement", "error", err)
	}
	defer stmt.Close()

	rows , err := stmt.QueryContext(ctx ,name)
	if err != nil{
		slog.Error("error executing", "error", err)
		return nil, handleSqlError(err)
	}
	defer rows.Close()

	var readlists []*model.Readlist
	for rows.Next(){
		var readlist model.Readlist
		if err := rows.Scan(&readlist.Id, &readlist.UserName, &readlist.Name, &readlist.Description, &readlist.Created_at, &readlist.Updated_at); err != nil{
			return nil, handleSqlError(err)
		}
		readlists= append(readlists, &readlist)
	}
	if err = rows.Err(); err != nil {
		return nil, handleSqlError(err)
	}
	return readlists, nil
}


func (rl *ReadlistRepo) GetUserReadlist(ctx context.Context, userId int) ([]*model.Readlist, error){
	query := `call get_readlists_from_user(?)`

	stmt, err := rl.DB.Prepare(query)
	if err != nil {
		slog.Error("error at preparing statement", "error", err)
	}
	defer stmt.Close()

	rows , err := stmt.QueryContext(ctx ,userId)
	if err != nil{
		slog.Error("error executing", "error", err)
		return nil, handleSqlError(err)
	}
	defer rows.Close()

	var readlists []*model.Readlist
	for rows.Next(){
		var readlist model.Readlist
		if err := rows.Scan(&readlist.Id, &readlist.UserName, &readlist.Name, &readlist.Description, &readlist.Created_at, &readlist.Updated_at); err != nil{
			return nil, handleSqlError(err)
		}
		readlists= append(readlists, &readlist)
	}
	if err = rows.Err(); err != nil {
		return nil, handleSqlError(err)
	}
	return readlists, nil
}

func (rl *ReadlistRepo) AddToReadlist(ctx context.Context, userId int,readlistItem *model.NewReadlistItem)error{
	query := `call add_to_readlist(?, ?, ? ,?)`
	
	stmt , err := rl.DB.Prepare(query)

	if err != nil {
		slog.Error("error prepareing statement", "error", err)
		return err
	}
	defer stmt.Close()
	_ , err = stmt.QueryContext(ctx,userId,readlistItem.MangaId, readlistItem.Status, readlistItem.ReadlistId)
	if err != nil {
		slog.Error("error executing procedure", "error", err)
		return handleSqlError(err)
	}
	return nil
}

func (rl *ReadlistRepo) GetReadlistItem(ctx context.Context, id int)([]*model.ReadlistItem, error){
	query := `call get_manga_list_from_readlist(?);`
	stmt, err := rl.DB.Prepare(query)
	if err != nil {
		slog.Error("error at preparign statement", "error" , err)
		return nil,err
	}
	defer stmt.Close()
	
	rows, err := stmt.QueryContext(ctx,id)
	if err != nil {
		slog.Error("error not found","error", err )
		return nil, handleSqlError(err)
	}
	defer rows.Close()

	var mangas []*model.ReadlistItem
	for rows.Next(){
		var manga model.ReadlistItem
		if err := rows.Scan(&manga.Id,&manga.MangaId, &manga.Title, &manga.Status, &manga.AddedAt); err != nil{
			return nil, handleSqlError(err)
		}
		mangas = append(mangas, &manga)
	}
	if len(mangas) == 0{
		return nil, handleSqlError(err)
	}
	if err = rows.Err(); err != nil {
		return nil, handleSqlError(err) 
	}
	return mangas, nil
} 

func (rl *ReadlistRepo) DeleteReadlistItem(ctx context.Context, userId int, readlistItemId int) error{
	query := `call delete_readlist_item(?,?)`

	stmt, err := rl.DB.Prepare(query)
	if err != nil {
		slog.Error("error prepareing statement", "error", err)
		return err
	}
	defer stmt.Close()

	_ , err = stmt.ExecContext(ctx, readlistItemId, userId)
	if err != nil {
		slog.Error("error at executing procedure", "error", err)
		return handleSqlError(err)
	}

	return nil
}


func (rl *ReadlistRepo) UpdateReadlistItemStatus(ctx context.Context, readlistId int,userId int,status string) error {
	query := `call update_readlist_item_status(?,?,?)`

	stmt, err := rl.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement","error", err)
		return err
	}
	defer stmt.Close()

	_ , err = stmt.ExecContext(ctx,readlistId,status , userId)
	if err != nil {
		slog.Error("error at executing procedure", "error", err)
		return handleSqlError(err)
	}
	return nil
}