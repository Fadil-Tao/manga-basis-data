package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/Fadil-Tao/manga-basis-data/configs"
	"github.com/Fadil-Tao/manga-basis-data/db"
	"github.com/Fadil-Tao/manga-basis-data/internal/handlers"
	"github.com/Fadil-Tao/manga-basis-data/internal/repository"
	"github.com/Fadil-Tao/manga-basis-data/internal/services"
	"github.com/Fadil-Tao/manga-basis-data/utils/loggers"
	"github.com/joho/godotenv"
)

func main() {
	handlerOpts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key != slog.TimeKey {
				return a
			}
			t := a.Value.Time()

			a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05.000 -07"))
			return a
		},
		Level: slog.LevelDebug,
	}
	consoleLogger := loggers.NewHandler(handlerOpts)

	logger := slog.New(consoleLogger)
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
	}

	cfg := config.New()

	Conn := db.InitDB()
	defer Conn.Close()

	// repository inject
	userRepo := repository.NewuserRepo(Conn)
	authorRepo := repository.NewAuthorRepo(Conn)
	genreRepo := repository.NewGenrerepo(Conn)
	mangaRepo := repository.NewMangaRepo(Conn)
	reviewRepo := repository.NewReviewRepo(Conn)
	ratingRepo := repository.NewRatingRepo(Conn)
	mangaService := services.NewMangaService(mangaRepo)
	readlistRepo := repository.NewReadlistRepo(Conn)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthcheck", handlers.CheckHealth)
	api := http.NewServeMux()
	api.Handle("/api/", http.StripPrefix("/api", mux))

	handlers.NewUserHandler(mux, userRepo)
	handlers.NewAuthorHandler(mux, authorRepo)
	handlers.NewGenreHandler(mux, genreRepo)
	handlers.NewMangaHandler(mux, mangaService)
	handlers.NewReviewHandler(mux, reviewRepo)
	handlers.NewRatingHandler(mux, ratingRepo)
	handlers.NewReadlistHandler(mux, readlistRepo)

	server := http.Server{
		Addr:         ":8080",
		Handler:      api,
	}
	go func() {
		slog.Info("Server succesfully started started", "Port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server startup failed", "error", err)
		}
		slog.Info("Stopped serving new connections...")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err = server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP Shutdown error", "error", err)
	}
	slog.Info("Graceful shutdown complete")
}
