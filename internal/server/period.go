package server

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ArunGautham-Soundarrajan/nomnomwatch/internal/templates"
)

// parsePeriod reads ?year= and ?month= from the request. A missing year
// defaults to the current year; a missing month means the whole year.
func parsePeriod(r *http.Request, now time.Time) (templates.Period, error) {
	p := templates.Period{Year: now.Year()}

	if v := r.URL.Query().Get("year"); v != "" {
		year, err := strconv.Atoi(v)
		if err != nil || year < 2000 || year > 3000 {
			return p, errors.New("invalid year")
		}
		p.Year = year
	}

	if v := r.URL.Query().Get("month"); v != "" {
		month, err := strconv.Atoi(v)
		if err != nil || month < 1 || month > 12 {
			return p, errors.New("invalid month")
		}
		p.Month = month
	}

	if p.Month == 0 {
		p.Start = time.Date(p.Year, time.January, 1, 0, 0, 0, 0, now.Location())
		p.End = p.Start.AddDate(1, 0, 0)
	} else {
		p.Start = time.Date(p.Year, time.Month(p.Month), 1, 0, 0, 0, 0, now.Location())
		p.End = p.Start.AddDate(0, 1, 0)
	}
	return p, nil
}
