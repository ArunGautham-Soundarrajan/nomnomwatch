package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ArunGautham-Soundarrajan/nomnomwatch/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Server struct {
	db      *pgxpool.Pool
	queries *db.Queries
}

func NewServer(pool *pgxpool.Pool) *Server {
	return &Server{
		db:      pool,
		queries: db.New(pool),
	}
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	err := s.db.Ping(r.Context())
	if err != nil {
		http.Error(w, "database not ready", http.StatusServiceUnavailable)
		return
	}

	fmt.Fprintf(w, "Healthy")
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.health)
	mux.HandleFunc("GET /orders", s.orders)
	return mux
}

func (s *Server) orders(w http.ResponseWriter, r *http.Request) {
	rows, err := s.queries.ListOrders(r.Context())
	if err != nil {
		log.Printf("list orders failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file")
	}

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(context.TODO()); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	s := NewServer(dbpool)
	httpSrv := &http.Server{
		Addr:              ":8080",
		Handler:           s.routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Println("listening on %s", httpSrv.Addr)
	log.Fatal(httpSrv.ListenAndServe())
}
