package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ArunGautham-Soundarrajan/nomnomwatch/internal/db"
	"github.com/ArunGautham-Soundarrajan/nomnomwatch/internal/templates"
	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5"
)

// Panel handlers return HTML fragments that htmx swaps into the dashboard.

func (s *Server) lastOrderPanel(w http.ResponseWriter, r *http.Request) {
	last, err := s.queries.GetLastOrderTime(r.Context())
	if errors.Is(err, pgx.ErrNoRows) {
		render(w, r, templates.LastOrderEmpty())
		return
	}
	if err != nil {
		panelError(w, r, "last order", err)
		return
	}

	render(w, r, templates.LastOrder(last, time.Now()))
}

func (s *Server) summaryPanel(w http.ResponseWriter, r *http.Request) {
	p, err := parsePeriod(r, time.Now())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render(w, r, templates.PanelError(err.Error()))
		return
	}

	sum, err := s.queries.GetPeriodSummary(r.Context(), db.GetPeriodSummaryParams{
		StartAt: p.Start,
		EndAt:   p.End,
	})
	if err != nil {
		panelError(w, r, "summary", err)
		return
	}

	render(w, r, templates.Summary(p, sum.Total, sum.OrderCount))
}

func (s *Server) topRestaurantsPanel(w http.ResponseWriter, r *http.Request) {
	p, err := parsePeriod(r, time.Now())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render(w, r, templates.PanelError(err.Error()))
		return
	}

	rows, err := s.queries.ListFavRestaurants(r.Context(), db.ListFavRestaurantsParams{
		StartAt: p.Start,
		EndAt:   p.End,
	})
	if err != nil {
		panelError(w, r, "top restaurants", err)
		return
	}

	render(w, r, templates.TopRestaurants(p, rows))
}

func (s *Server) monthlySpendPanel(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	p, err := parsePeriod(r, now)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render(w, r, templates.PanelError(err.Error()))
		return
	}

	// The chart covers the selected year, from the month of the first ever
	// order (earlier months have no data, not zero spending) up to the
	// current month.
	first, err := s.queries.GetFirstOrderTime(r.Context())
	if errors.Is(err, pgx.ErrNoRows) {
		render(w, r, templates.MonthlySpendChart(templates.MonthlySpend{Period: p, Partial: -1}))
		return
	}
	if err != nil {
		panelError(w, r, "monthly spend", err)
		return
	}

	start := time.Date(p.Year, time.January, 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(1, 0, 0)
	if firstMonth := monthStart(first.In(now.Location())); firstMonth.After(start) {
		start = firstMonth
	}
	if nextMonth := monthStart(now).AddDate(0, 1, 0); nextMonth.Before(end) {
		end = nextMonth
	}

	rows, err := s.queries.ListOrdersBetween(r.Context(), db.ListOrdersBetweenParams{
		StartAt: start,
		EndAt:   end,
	})
	if err != nil {
		panelError(w, r, "monthly spend", err)
		return
	}

	render(w, r, templates.MonthlySpendChart(monthlySpend(p, start, end, now, rows)))
}

// monthlySpend buckets orders into one point per month in [start, end).
// Grouping happens here rather than in SQL so month boundaries use the same
// timezone as parsePeriod.
func monthlySpend(p templates.Period, start, end, now time.Time, rows []db.ListOrdersBetweenRow) templates.MonthlySpend {
	m := templates.MonthlySpend{Period: p, Partial: -1}
	for month := start; month.Before(end); month = month.AddDate(0, 1, 0) {
		if month.Year() == now.Year() && month.Month() == now.Month() {
			m.Partial = len(m.Months)
		}
		m.Months = append(m.Months, templates.MonthPoint{Month: month})
	}
	for _, row := range rows {
		d := row.OrderDate.In(start.Location())
		i := (d.Year()-start.Year())*12 + int(d.Month()) - int(start.Month())
		if i >= 0 && i < len(m.Months) {
			m.Months[i].Total += row.Total
			m.Months[i].Orders++
		}
	}
	return m
}

func monthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	if err := c.Render(r.Context(), w); err != nil {
		log.Printf("render failed: %v", err)
	}
}

func panelError(w http.ResponseWriter, r *http.Request, panel string, err error) {
	log.Printf("%s panel failed: %v", panel, err)
	w.WriteHeader(http.StatusInternalServerError)
	render(w, r, templates.PanelError("Couldn't load this panel."))
}
