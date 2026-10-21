package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	err := s.db.Ping(r.Context())
	if err != nil {
		http.Error(w, "database not ready", http.StatusServiceUnavailable)
		return
	}

	fmt.Fprintf(w, "Healthy")
}

func (s *Server) orders(w http.ResponseWriter, r *http.Request) {
	rows, err := s.queries.ListOrdersMonthly(r.Context(), 2026)
	if err != nil {
		log.Printf("list orders failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}
