package server

import (
	"log"
	"net/http"
	"time"

	"github.com/ArunGautham-Soundarrajan/nomnomwatch/internal/templates"
)

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	p, err := parsePeriod(r, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	years, err := s.queries.ListOrderYears(r.Context())
	if err != nil {
		log.Printf("list order years failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if len(years) == 0 {
		years = []int32{int32(now.Year())}
	}

	if err := templates.Index(p, years).Render(r.Context(), w); err != nil {
		log.Printf("render index failed: %v", err)
	}
}
