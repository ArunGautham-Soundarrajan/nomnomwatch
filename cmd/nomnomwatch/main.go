package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ArunGautham-Soundarrajan/nomnomwatch/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

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

	s := server.New(dbpool)
	httpSrv := &http.Server{
		Addr:              ":8080",
		Handler:           s.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("listening on %s", httpSrv.Addr)
	log.Fatal(httpSrv.ListenAndServe())
}
