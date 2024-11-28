package repository

import (
	"context"
	"database/sql"
	"log/slog"
)


type RatingRepo struct {
	DB *sql.DB
}

func NewRatingRepo(db *sql.DB)*RatingRepo{
	return &RatingRepo{
		DB: db,
	}
}


func (r *RatingRepo) RateManga(ctx context.Context, userId int, mangaId int, rating int)error{
	query := `call rate_manga(?,?,?)`

	stmt, err := r.DB.Prepare(query)
	if err != nil {
		slog.Error("error preparing statement","error", err)
		return err
	}
	defer stmt.Close()

	_ ,err = stmt.ExecContext(ctx,mangaId,userId,rating)
	if err != nil {
		slog.Error("error executing procedur","error", err)
		return handleSqlError(err)
	}
	return nil
}