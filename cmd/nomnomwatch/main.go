package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ArunGautham-Soundarrajan/nomnomwatch/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel(),
	})))

	// In the cluster, config comes from the environment, so a missing .env
	// file is expected.
	if err := godotenv.Load(); err != nil {
		slog.Debug("no .env file, using environment only")
	}

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fatal("database connection error", err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(context.TODO()); err != nil {
		fatal("database ping failed", err)
	}

	s := server.New(dbpool)
	httpSrv := &http.Server{
		Addr:              ":8080",
		Handler:           s.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	slog.Info("listening", "addr", httpSrv.Addr)
	fatal("server stopped", httpSrv.ListenAndServe())
}

// logLevel reads LOG_LEVEL (debug, info, warn, error); the default is info.
func logLevel() slog.Level {
	var level slog.Level
	if err := level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))); err != nil {
		return slog.LevelInfo
	}
	return level
}

func fatal(msg string, err error) {
	slog.Error(msg, "err", err)
	os.Exit(1)
}
