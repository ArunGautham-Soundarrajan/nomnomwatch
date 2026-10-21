package server

import (
	"testing"
	"time"

	"github.com/ArunGautham-Soundarrajan/nomnomwatch/internal/db"
	"github.com/ArunGautham-Soundarrajan/nomnomwatch/internal/templates"
)

func TestMonthlySpend(t *testing.T) {
	loc := time.UTC
	start := time.Date(2026, time.August, 1, 0, 0, 0, 0, loc)
	end := time.Date(2026, time.November, 1, 0, 0, 0, 0, loc)
	now := time.Date(2026, time.October, 21, 12, 0, 0, 0, loc)
	rows := []db.ListOrdersBetweenRow{
		{OrderDate: time.Date(2026, time.August, 3, 19, 0, 0, 0, loc), Total: 10},
		{OrderDate: time.Date(2026, time.August, 31, 23, 59, 0, 0, loc), Total: 5},
		{OrderDate: time.Date(2026, time.October, 1, 0, 0, 0, 0, loc), Total: 7.5},
	}

	m := monthlySpend(templates.Period{Year: 2026}, start, end, now, rows)

	want := []struct {
		month  time.Month
		total  float64
		orders int64
	}{
		{time.August, 15, 2},
		{time.September, 0, 0},
		{time.October, 7.5, 1},
	}
	if len(m.Months) != len(want) {
		t.Fatalf("got %d months, want %d", len(m.Months), len(want))
	}
	for i, w := range want {
		got := m.Months[i]
		if got.Month.Month() != w.month || got.Total != w.total || got.Orders != w.orders {
			t.Errorf("month %d = %v %.2f %d, want %v %.2f %d", i, got.Month.Month(), got.Total, got.Orders, w.month, w.total, w.orders)
		}
	}
	if m.Partial != 2 {
		t.Errorf("partial = %d, want 2", m.Partial)
	}
	if avg, ok := m.Average(); !ok || avg != 7.5 {
		t.Errorf("average = %.2f, %v; want 7.50 (Aug + Sep, excluding in-progress Oct)", avg, ok)
	}
}
