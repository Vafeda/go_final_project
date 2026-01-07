package nextdate

import (
	"fmt"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", fmt.Errorf("empty repeat")
	}

	date, err := time.Parse(DateFormat, dstart)
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

	return date.Format(DateFormat), nil
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
	for {
		date = date.AddDate(0, 0, d.day)
		if afterNow(date, now) {
			break
		}
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
	date = date.AddDate(1, 0, 0)

	for !afterNow(date, now) {
		date = date.AddDate(1, 0, 0) // Нужно присваивать результат!
	}
	return date
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
