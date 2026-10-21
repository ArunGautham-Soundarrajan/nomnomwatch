package templates

import (
	"strconv"
	"time"
)

// Period is the date range the dashboard is showing: Start is inclusive,
// End is exclusive. Month is 0 when the whole year is selected.
type Period struct {
	Year  int
	Month int
	Start time.Time
	End   time.Time
}

func (p Period) Label() string {
	if p.Month == 0 {
		return strconv.Itoa(p.Year)
	}
	return time.Month(p.Month).String() + " " + strconv.Itoa(p.Year)
}
