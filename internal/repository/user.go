package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Fadil-Tao/manga-basis-data/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewuserRepo(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) RegisterUser(ctx context.Context, newUser *model.NewUserRequest) error {
	query := `CALL register_user(?,?,?);`

	stmt, err := u.db.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, newUser.Username, newUser.Email, newUser.Password)
	if err != nil {
		slog.Error("Error executing procedure", "err", err)
		return handleSqlError(err)
	}
	return nil
}

func (u *UserRepository) Login(ctx context.Context, user *model.UserLoginRequest) (*model.UserResponse, error) {
	query := `CALL login_user(?,?);`

	stmt, err := u.db.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, user.Email, user.Password)

	var User model.UserResponse

	err = row.Scan(&User.Id, &User.Username, &User.Email, &User.Is_admin, &User.Created_at)
	if err != nil {
		slog.Error("Error at scaning row")
		return nil, handleSqlError(err)
	}
	return &User, nil
}


func (u *UserRepository) SearchUser(ctx context.Context, username string)([]*model.UserResponse, error){
	query := `call get_all_user_with_search_by_username(?)`

	stmt, err := u.db.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return nil, err
	}
	defer stmt.Close()


	rows,err := stmt.QueryContext(ctx, username)
	if err != nil {
		return nil, handleSqlError(err)
	}
	defer rows.Close()
	var Users []*model.UserResponse

	for rows.Next(){
		var user model.UserResponse
		if err := rows.Scan(&user.Id, &user.Username,&user.Email , &user.Created_at); err != nil {
			return nil , handleSqlError(err)
		}
		Users = append(Users, &user)
	}
	if len(Users) == 0 {
		return nil,handleSqlError(err)
	}
	if err = rows.Err(); err != nil {
		return Users, handleSqlError(err) 
	}
	return Users, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, usernameTarget string,userid int,userData *model.NewUserRequest) error{
	query := `call update_user(?,?,?,?,?)`
	
	stmt, err := u.db.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,userid,usernameTarget,&userData.Username,&userData.Password, &userData.Email)
	if err != nil {
		slog.Error("error at executing procedure", "error", err)
		return handleSqlError(err)
	}
	return nil
}

func (u *UserRepository) GetUserDetailByUsername(ctx context.Context, username string)(*model.UserResponse, error){
	query := `call get_user_detail_by_username(?)`
	stmt, err := u.db.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return nil , err
	}
	defer stmt.Close()

	row:= stmt.QueryRowContext(ctx,username)

	var user model.UserResponse
	err = row.Scan(&user.Id, &user.Username,&user.Email,&user.Created_at)
	if err != nil {
		return nil, handleSqlError(err)
	}
	return &user, nil
}

func (u *UserRepository) GetUserRatedManga(ctx context.Context, username string)([]*model.UserRatedManga, error){
	query := `call get_user_rated_manga(?)`
	stmt, err := u.db.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement", "error", err)
		return nil , err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var mangas []*model.UserRatedManga 
	for rows.Next(){
		var manga model.UserRatedManga

		if err := rows.Scan(&manga.Id, &manga.Title, &manga.Manga_status,&manga.Published_at, &manga.Finished_at, &manga.Created_at, &manga.YourRating, &manga.Rating, &manga.TotalUserRated); err != nil{
			slog.Error("error scaning row", "error", err)
			return nil, handleSqlError(err)
		}
		mangas = append(mangas, &manga)
	}
	if len(mangas)== 0 {
		slog.Warn("err here at len 0")
		return nil,handleSqlError(err)
	}
	if err = rows.Err(); err != nil {
		slog.Warn("err here at")
		return nil, handleSqlError(err) 
	}
	return mangas, nil
}
func (u *UserRepository) GetUserLikedManga(ctx context.Context, username string)([]*model.UserLikedManga, error){
	query := `call get_user_liked_manga(?)`


	rows, err := u.db.QueryContext(ctx, query,username)
	if err != nil {
		return nil, handleSqlError(err)
	}
	defer rows.Close()
	var mangas []*model.UserLikedManga 
	for rows.Next(){
		var manga model.UserLikedManga
		if err := rows.Scan(&manga.Id, &manga.Title, &manga.Manga_status,&manga.Published_at, &manga.Finished_at, &manga.Created_at, &manga.TotalLikes); err != nil{
			slog.Error("error scaning row", "error", err)
			return nil, handleSqlError(err)
		}
		mangas = append(mangas, &manga)
	}
	if len(mangas)== 0 {
		return nil,handleSqlError(err)
	}
	if err = rows.Err(); err != nil {
		return nil, handleSqlError(err) 
	}
	return mangas, nil
}

func (u *UserRepository) GetUserReadlist(ctx context.Context, username string) ([]*model.Readlist, error){
	query := `call get_readlist_from_user(?)`


	rows , err := u.db.QueryContext(ctx,query, username)
	if err != nil{
		slog.Error("error executing", "error", err)
		return nil, handleSqlError(err)
	}
	defer rows.Close()
	var readlists []*model.Readlist
	for rows.Next(){
		var readlist model.Readlist
		if err := rows.Scan(&readlist.Id, &readlist.Name, &readlist.Description, &readlist.Created_at, &readlist.Updated_at); err != nil{
			return nil, handleSqlError(err)
		}
		readlists= append(readlists, &readlist)
	}
	if err = rows.Err(); err != nil {
		return nil, handleSqlError(err)
	}
	return readlists, nil
} 

func (u *UserRepository) DeleteUser(ctx context.Context, userId int,username string)error{
	query := `call delete_user(?,?);`

	_ ,err := u.db.ExecContext(ctx, query,username, userId)
	if err != nil {
		slog.Error("error executing procedure","error", err)
		return handleSqlError(err)
	}
	return nil
} 