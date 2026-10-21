package server

import (
	"net/http"

	"github.com/ArunGautham-Soundarrajan/nomnomwatch/internal/assets"
	"github.com/ArunGautham-Soundarrajan/nomnomwatch/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	db      *pgxpool.Pool
	queries *db.Queries
}

func New(pool *pgxpool.Pool) *Server {
	return &Server{
		db:      pool,
		queries: db.New(pool),
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.index)
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /orders", s.orders)
	mux.Handle("GET /static/", http.FileServerFS(assets.Static))
	mux.HandleFunc("GET /panels/summary", s.summaryPanel)
	mux.HandleFunc("GET /panels/monthly-spend", s.monthlySpendPanel)
	mux.HandleFunc("GET /panels/last-order", s.lastOrderPanel)
	mux.HandleFunc("GET /panels/top-restaurants", s.topRestaurantsPanel)
	return mux
}
