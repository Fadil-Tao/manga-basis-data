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

func NewuserRepo(db *sql.DB) *UserRepository{
	return &UserRepository{
			db: db,
	}
}

func (u *UserRepository) RegisterUser(ctx context.Context,newUser *model.NewUserRequest ) error{
	query := `CALL register_user(?,?,?);`
	
	stmt ,err := u.db.Prepare(query)
	if err != nil{
		slog.Error("error preparing statement", "error", err)
		return err
	}

	defer stmt.Close()

	_ ,err = stmt.ExecContext(ctx,newUser.Username,newUser.Email, newUser.Password )
	if err != nil {
		slog.Error("Error executing procedure", "err", err)
		return err
	}
	return nil
}

func (u *UserRepository) Login(ctx  context.Context, user *model.UserLoginRequest) (*model.UserResponse, error){
	query := `CALL login_user(?,?);`

	stmt, err := u.db.Prepare(query)
	if err != nil{
		slog.Error("error preparing statement", "error", err)
		return nil,err
	}
	defer stmt.Close()
	
	row := stmt.QueryRowContext(ctx, user.Email, user.Password)

	var User model.UserResponse

	err = row.Scan(&User.Id,&User.Username,&User.Email,&User.Is_admin,&User.Created_at)
	if err != nil {
		slog.Error("Error at scaning row")
		return nil, err
	}
	return &User, nil 
}

// func (u *UserRepository)GetAllUser(ctx )