package nextdate

import (
	"fmt"
	"strings"
	"time"
)

type nextDate struct {
	date   string
	repeat string
	want   string
}

func Test() {
	now, _ := time.Parse("20060102", "20240126")
	tbl := []nextDate{
		{"20240126", "", ""},
		{"20240126", "k 34", ""},
		{"20240126", "ooops", ""},
		{"15000156", "y", ""},
		{"ooops", "y", ""},
		{"16890220", "y", `20240220`},
		{"20250701", "y", `20260701`},
		{"20240101", "y", `20250101`},
		{"20231231", "y", `20241231`},
		{"20240229", "y", `20250301`},
		{"20240301", "y", `20250301`},
		{"20240113", "d", ""},
		{"20240113", "d 7", `20240127`},
		{"20240120", "d 20", `20240209`},
		{"20240202", "d 30", `20240303`},
		{"20240320", "d 401", ""},
		{"20231225", "d 12", `20240130`},
		{"20240228", "d 1", "20240229"},
	}

	for _, v := range tbl {
		res, _ := NextDate(now, v.date, v.repeat)
		if res != v.want {
			fmt.Println(v.date, v.repeat, v.want, res)
		}
	}
	tbl = []nextDate{
		{"20231106", "m 13", "20240213"},
		{"20240120", "m 40,11,19", ""},
		{"20240116", "m 16,5", "20240205"},
		{"20240126", "m 25,26,7", "20240207"},
		{"20240409", "m 31", "20240531"},
		{"20240329", "m 10,17 12,8,1", "20240810"},
		{"20230311", "m 07,19 05,6", "20240507"},
		{"20230311", "m 1 1,2", "20240201"},
		{"20240127", "m -1", "20240131"},
		{"20240222", "m -2", "20240228"},
		{"20240222", "m -2,-3", ""},
		{"20240326", "m -1,-2", "20240330"},
		{"20240201", "m -1,18", "20240218"},
		{"20240125", "w 1,2,3", "20240129"},
		{"20240126", "w 7", "20240128"},
		{"20230126", "w 4,5", "20240201"},
		{"20230226", "w 8,4,5", ""},
	}

	for _, v := range tbl {
		res, _ := NextDate(now, v.date, v.repeat)
		if res != v.want {
			fmt.Println(v.date, v.repeat, v.want, res)
		}
	}
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", fmt.Errorf("empty repeat")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	r := RepeatInfo{}
	err = r.Parse(repeat)
	if err != nil {
		return "", fmt.Errorf("invalid repeat rule: %w", err)
	}

	date, err = calculateNextDate(now, date, r)
	if err != nil {
		return "", err
	}

	return date.Format(dateFormat), nil
}

func calculateNextDate(now, date time.Time, r RepeatInfo) (time.Time, error) {
	switch v := r.data.(type) {
	case Day:
		return calculateForDay(now, date, v), nil
	case Week:
		return calculateForWeek(now, date, v), nil
	case Month:
		return calculateForMonth(now, date, v), nil
	case Year:
		return calculateForYear(now, date), nil
	default:
		return time.Time{}, fmt.Errorf("unsupported repeat type")
	}
}

func calculateForDay(now time.Time, date time.Time, d Day) time.Time {
	date = date.AddDate(0, 0, d.day)

	for !afterNow(date, now) {
		date = date.AddDate(0, 0, d.day)
	}

	return date
}

func calculateForWeek(now time.Time, date time.Time, w Week) time.Time {
	if !afterNow(date, now) {
		date = now
	}

	var count int

	for i := date.Weekday() + 1; i <= time.Saturday; i++ {
		count++
		if w.weekDays[i] {
			date = now.AddDate(0, 0, count)
			return date
		}
	}

	for i := time.Sunday; i <= date.Weekday(); i++ {
		count++
		if w.weekDays[i] {
			date = now.AddDate(0, 0, count)
			break
		}
	}

	return date
}

func calculateForMonth(now time.Time, date time.Time, m Month) time.Time {
	if !afterNow(date, now) {
		date = now
	}

	currentYear := date.Year()
	currentMonth := date.Month()
	currentDay := date.Day() + 1

	for ; currentMonth <= time.December; currentMonth++ {
		if !m.monthsEmpty && !m.months[currentMonth-1] {
			if currentDay != 1 {
				currentDay = 1
			}
			continue
		}

		for ; currentDay <= daysIn(currentMonth, currentYear); currentDay++ {
			if m.days[currentDay-1] {
				return time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, time.UTC)

			}
		}

		if m.days[33] {
			return time.Date(currentYear, currentMonth, daysIn(currentMonth, currentYear)-1, 0, 0, 0, 0, time.UTC)
		}

		if m.days[32] {
			return time.Date(currentYear, currentMonth, daysIn(currentMonth, currentYear), 0, 0, 0, 0, time.UTC)
		}

		currentDay = 1
	}

	currentYear += 1
	currentMonth = time.January

	for ; currentMonth <= date.Month(); currentMonth++ {
		if !m.monthsEmpty && !m.months[currentMonth-1] {
			if currentDay != 1 {
				currentDay = 1
			}
			continue
		}

		for ; currentDay <= daysIn(currentMonth, currentYear); currentDay++ {
			if m.days[currentDay-1] {
				return time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, time.UTC)
			}
		}

		if m.days[33] {
			return time.Date(currentYear, currentMonth, daysIn(currentMonth, currentYear)-1, 0, 0, 0, 0, time.UTC)
		}

		if m.days[32] {
			return time.Date(currentYear, currentMonth, daysIn(currentMonth, currentYear), 0, 0, 0, 0, time.UTC)
		}

		currentDay = 1
	}

	return time.Time{}
}

func calculateForYear(now time.Time, date time.Time) time.Time {
	if date.Year() >= now.Year() {
		return date.AddDate(1, 0, 0)
	}

	return date.AddDate(now.Year()-date.Year(), 0, 0)
}

func afterNow(date time.Time, now time.Time) bool {
	if date.Year() > now.Year() {
		return true
	}
	if date.Year() < now.Year() {
		return false
	}

	if date.Month() > now.Month() {
		return true
	}
	if date.Month() < now.Month() {
		return false
	}

	return date.Day() > now.Day()
}

func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
