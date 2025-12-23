package nextdate

import (
	"time"
)

type nextDate struct {
	date   string
	repeat string
	want   string
	error  interface{}
}

var (
	now = time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
)

//func TestNextDate(t *testing.T) {
//	dayTest := []nextDate{
//		{"20260101", "d", "", ErrInvalidArgsCount},
//		{"20260101", "d 1 1", "", ErrInvalidArgsCount},
//		{"20260101", "d 0", "", ErrArgumentOutOfRange},
//		{"20260101", "d 401", "", ErrArgumentOutOfRange},
//	}
//
//	for _, dt := range dayTest {
//		res, err := NextDate(now, dt.date, dt.repeat)
//		if errors.As(err, &dt.error) {
//
//		}
//		if err {
//		}
//	}
//}
