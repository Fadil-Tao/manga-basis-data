package services

import (
	"context"
	"log/slog"

	"github.com/Fadil-Tao/manga-basis-data/internal/model"
)

type MangaRepo interface {
	CreateManga(ctx context.Context, manga *model.Manga,userId int) error
	GetMangaById(ctx context.Context, id string) (*model.MangaList, error)
	GetMangaAuthor(ctx context.Context, idManga string) ([]model.Author, error)
	GetMangaGenre(ctx context.Context, idManga string) ([]model.Genre, error)
	ConnectMangaAuthor(ctx context.Context, obj *model.MangaAuthorPivot, userId int) error
	ConnectMangaGenre(ctx context.Context, obj *model.MangaGenrePivot, userId int) error
	GetAllManga(ctx context.Context, title string) ([]*model.MangaList, error)
	SearchMangaByName(ctx context.Context, name string) ([]model.Manga, error)
	DeleteMangaById(ctx context.Context, id int, userId int) error 
	UpdateManga(ctx context.Context,id int,manga *model.Manga, userId int)error
	GetMangaRankingList(ctx context.Context, period string)([]*model.MangaList, error)
	ToggleLikeManga(ctx context.Context, userId int , mangaId int)error
	DeleteMangaAuthorConnection(ctx context.Context, userId int, mangaId int, authorId int)error
	DeleteMangaGenreConnection(ctx context.Context, userId int, mangaId int, genreId int)error
}

type MangaService struct {
	MangaRepo MangaRepo
}

func NewMangaService(mangaRepo MangaRepo) *MangaService {
	return &MangaService{
		MangaRepo: mangaRepo,
	}
}

func (m *MangaService) CreateManga(ctx context.Context, manga *model.Manga, userId int) error {
	if *manga.Manga_status == "in_progress"{
		manga.Finished_at = nil
	}
	err := m.MangaRepo.CreateManga(ctx, manga,userId)
	if err != nil {
		slog.Error("Error", "message", err)
		return err
	}
	return nil
}

func (m *MangaService) GetMangaById(ctx context.Context, id string) (*model.MangaResponse, error) {
	mangaData, err := m.MangaRepo.GetMangaById(ctx, id)
	if err != nil {
		slog.Error("Error at getting manga data", "message", err)
		return nil, err
	}

	authors, err := m.MangaRepo.GetMangaAuthor(ctx, id)
	if err != nil {
		slog.Error("Error at getting author data", "message", err)
		return nil, err
	}

	genres, err := m.MangaRepo.GetMangaGenre(ctx, id)
	if err != nil {
		slog.Error("Error at getting author data", "message", err)
		return nil, err
	}

	return &model.MangaResponse{
		MangaList:  *mangaData,
		Genres: genres,
		Author: authors,
	}, nil
}

func (m *MangaService) ConnectMangaAuthor(ctx context.Context, ma *model.MangaAuthorPivot, userId int) error {
	err := m.MangaRepo.ConnectMangaAuthor(ctx, ma, userId)
	if err != nil {
		slog.Error("error at connecting to repo", "message", err)
		return err
	}
	return nil
}

func (m *MangaService) ConnectMangaGenre(ctx context.Context, mg *model.MangaGenrePivot, userId int) error {
	err := m.MangaRepo.ConnectMangaGenre(ctx, mg, userId)
	if err != nil {
		slog.Error("error at connecting to repo", "message", err)
		return err
	}
	return nil
}

func (m *MangaService) GetAllManga(ctx context.Context,  name string) ([]*model.MangaList, error) {
	mangaList, err := m.MangaRepo.GetAllManga(ctx,  name)
	if err != nil {
		slog.Error("error at connecting to repo", "message", err)
		return nil, err
	}
	return mangaList, nil
}

func (m *MangaService ) SearchMangaByName(ctx context.Context, name string) ([]model.Manga, error){
	mangaList, err := m.MangaRepo.SearchMangaByName(ctx,name)
	if err != nil{
		slog.Error("error at connecting to repo", "message", err)
		return nil, err
	}
	return mangaList, nil
}

func (m *MangaService) DeleteMangaById(ctx context.Context, id int, userId int) error {
	err := m.MangaRepo.DeleteMangaById(ctx,id, userId)
	if err != nil {
		slog.Error("error at deleting", "message", err)
		return err
	}
	return nil
}

func (m *MangaService) UpdateManga(ctx context.Context, id int , manga *model.Manga, userId int)error{
	if *manga.Manga_status == "in_progress"{
		manga.Finished_at = nil
	}
	err := m.MangaRepo .UpdateManga(ctx, id, manga, userId )
	if err != nil{
		slog.Error("error at updating", "message", err)
		return err
	}
	return nil
}

func (m *MangaService) GetMangaRankingList(ctx context.Context, period string)([]*model.MangaList, error){
	mangaList, err := m.MangaRepo.GetMangaRankingList(ctx,period)
	if err != nil{
		slog.Error("error at connecting to repo", "message", err)
		return nil, err
	}
	return mangaList, nil
}

func (m *MangaService) ToggleLikeManga(ctx context.Context, userId int , mangaId int)error{
 	err := m.MangaRepo.ToggleLikeManga(ctx, userId,mangaId);if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func (m *MangaService) DeleteMangaAuthorConnection(ctx context.Context, userId int, mangaId int, authorId int)error{
	err := m.MangaRepo.DeleteMangaAuthorConnection(ctx, userId, mangaId, authorId)
	if err !=nil{
		slog.Error(err.Error())
		return err
	}
	return nil 
}


func (m *MangaService) DeleteMangaGenreConnection(ctx context.Context, userId int, mangaId int, genreId int)error{
	err := m.MangaRepo.DeleteMangaGenreConnection(ctx, userId, mangaId, genreId)
	if err !=nil{
		slog.Error(err.Error())
		return err
	}
	return nil 
}

