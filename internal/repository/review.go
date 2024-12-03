package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Fadil-Tao/manga-basis-data/internal/model"
)


type ReviewRepo struct {
	DB *sql.DB
}

func NewReviewRepo(db *sql.DB)*ReviewRepo{
	return &ReviewRepo{
		DB: db,
	}
}

func (r *ReviewRepo) CreateReview(ctx context.Context, review *model.NewReviewRequest)error{
	query := `call add_review(?,?,?,?)`
	
	stmt, err := r.DB.Prepare(query)
	if err != nil{
		slog.Error("error at preparing", "message",err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, review.Manga_id, review.User_id,review.Review_text,review.Tag) 
	if err != nil{
		slog.Error("error at executing proc","message",err )
		return handleSqlError(err)
	}

	return nil
}

func (r *ReviewRepo) GetReviewFromManga(ctx context.Context, mangaId int)([]*model.Review, error){
	query := `call get_review_from_manga(?)`

	stmt, err := r.DB.Prepare(query)
	if err != nil{
		slog.Error("error at preparing", "message",err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, mangaId)
	if err != nil {
		slog.Error("error not found", "error", err)
		return nil, handleSqlError(err)
	}
	defer rows.Close()
	
	var reviews []*model.Review
	for rows.Next(){
		var review model.Review
		if err := rows.Scan(&review.User_id, &review.Username, &review.Review_text, &review.Tag,&review.Created_at,&review.Total_Like); err != nil{
			return nil, handleSqlError(err)
		}
		reviews = append(reviews, &review)
	}
	if err = rows.Err(); err != nil{
		return nil, handleSqlError(err)
	}
	return reviews,nil	
}

func (r *ReviewRepo) DeleteReview(ctx context.Context, doerId int, mangaId int, reviewerId int)error{
	query := `call delete_a_review(?,?,?)`
	
	stmt,err := r.DB.Prepare(query)
	if err != nil {
		slog.Error("error at preparing statement", "error" ,err)
		return err
	}
	defer stmt.Close()

	_,err = stmt.ExecContext(ctx, doerId, mangaId,reviewerId)
	if err != nil {
		slog.Error("error at executing proc", "error" ,err)
		return handleSqlError(err)
	}
	return nil
}


func (r *ReviewRepo) UpdateReview(ctx context.Context, doerId int,review *model.UpdateReview) error{
	query := `call update_a_review(?,?,?,?,?)`

	stmt,err := r.DB.Prepare(query)
	if err != nil {
		slog.Error("error at preparing statement", "error" ,err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, doerId, review.User_id,review.Manga_id,review.Review_text,review.Tag)
	if err != nil {
		slog.Error("error at executing proc", "error" ,err)
		return handleSqlError(err)
	}
	return nil
}


func (r *ReviewRepo) GetAReviewById(ctx context.Context, mangaId int, userId int)(*model.Review, error){
	query := `call get_review_by_id(?,?)`

	stmt,err := r.DB.Prepare(query)
	if err != nil{
		slog.Error("error at preparing statement", "error" ,err)
		return nil,err
	}
	defer stmt.Close()
	
	row := stmt.QueryRowContext(ctx, mangaId, userId)

	var Review model.Review
	
	err = row.Scan(&Review.Manga_id, &Review.User_id, &Review.Username,&Review.Review_text, &Review.Tag, &Review.Created_at, &Review.Total_Like)
	if err != nil{
		slog.Error("error at scanning", "error", err)
		return nil, handleSqlError(err)
	}
	return &Review, nil
}

func (r *ReviewRepo)ToggleLikeReview(ctx context.Context, doer int, mangaId int, reviewerId int)error{
	query := `call like_unlike_a_review(?,?,?)`
	
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		slog.Error("error at preparing statement", "error" ,err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, doer, mangaId,reviewerId) 
	if err != nil{
		slog.Error("error at executing proc","message",err )
		return handleSqlError(err)
	}
	return nil
}

