package main

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/Fadil-Tao/manga-basis-data/utils/loggers"
	"github.com/joho/godotenv"
)


func main(){
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
	if err != nil{
		slog.Error(err.Error())
	}

	cfg := config.New()
	
	Conn := db.InitDB(&cfg.DB)
	defer Conn.Close()
	
	userRepo := repository.NewuserRepo(Conn)
	authorRepo := repository.NewAuthorRepo(Conn)


	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthcheck", handlers.CheckHealth)
	// api declaration in url route
	api := http.NewServeMux()
	api.Handle("/api/", http.StripPrefix("/api", mux))
	handlers.NewUserHandler(mux, userRepo)
	handlers.NewAuthorHandler(mux, authorRepo)
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      api,
		ReadTimeout:  cfg.Server.TimeoutRead,
		WriteTimeout: cfg.Server.TimeoutWrite,
		IdleTimeout:  cfg.Server.TimeoutIdle,
	}
	go func ()  {
		slog.Info("Server succesfully started started", "Port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server startup failed")
		}
		slog.Info("Stopped serving new connections...")
	}() 

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,syscall.SIGINT , syscall.SIGTERM) 
	<- sigChan 

	shutdownCtx , shutdownRelease := context.WithTimeout(context.Background(), 10 * time.Second)
	defer shutdownRelease()

	if err = server.Shutdown(shutdownCtx);err != nil {
		slog.Error("HTTP Shutdown error" , "error" ,err)
	}
	slog.Info("Graceful shutdown complete")
}