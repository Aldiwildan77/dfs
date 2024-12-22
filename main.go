package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339

	ctx := context.Background()

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	db, err := InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	discord := NewDiscord(cfg.DiscordToken)
	rotator := NewRotator(cfg, db, discord)
	handler := NewHandlers(&cfg, db, rotator)

	router := mux.NewRouter()

	router.Use(LogRequest)

	apiRouter := router.PathPrefix("/api").Subrouter()
	fileRouter := apiRouter.PathPrefix("/file").Subrouter()
	fileRouter.HandleFunc("/save", handler.SaveDiscordLink).Methods(http.MethodPost)
	fileRouter.HandleFunc("/get", handler.GetDiscordLink).Methods(http.MethodPost)

	go rotator.StartRotator(ctx, time.Duration(cfg.RotatorDelay)*time.Hour)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Info().
			Int("port", cfg.ServerPort).
			Msg("Starting server...")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	<-stop

	log.Info().Msg("Shutting down server...")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutDown); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
