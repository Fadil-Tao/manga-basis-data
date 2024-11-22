package services

import (
	"context"
	"log/slog"

	"github.com/Fadil-Tao/manga-basis-data/internal/model"
)

type MangaRepo interface {
	CreateManga(ctx context.Context, manga *model.Manga,userId int) error
	GetMangaById(ctx context.Context, id string) (*model.Manga, error)
	GetMangaAuthor(ctx context.Context, idManga string) ([]model.Author, error)
	GetMangaGenre(ctx context.Context, idManga string) ([]model.Genre, error)
	ConnectMangaAuthor(ctx context.Context, obj *model.MangaAuthorPivot, userId int) error
	ConnectMangaGenre(ctx context.Context, obj *model.MangaGenrePivot, userId int) error
	GetAllMangaWithLimit(ctx context.Context, limit int) ([]*model.MangaList, error)
	SearchMangaByName(ctx context.Context, name string) ([]model.Manga, error)
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
		Manga:  *mangaData,
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

func (m *MangaService) GetAllMangaWithLimit(ctx context.Context, limit int) ([]*model.MangaList, error) {
	mangaList, err := m.MangaRepo.GetAllMangaWithLimit(ctx, limit)
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