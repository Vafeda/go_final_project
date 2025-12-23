package nextdate

import (
	"fmt"
	"strconv"
	"strings"
)

type RepeatInfo struct {
	data interface{}
}

func (r *RepeatInfo) Parse(repeat string) error {
	strs := strings.Fields(repeat)

	switch strs[0] {
	case "d":
		if len(strs) != 2 {
			return ErrInvalidArgsCount
		}

		d := Day{}
		err := d.parseInfo(strs[1])
		if err != nil {
			return fmt.Errorf("day parse info: %w", err)
		}
		r.data = d

	case "y":
		if len(strs) != 1 {
			return ErrInvalidArgsCount
		}
		r.data = Year{}

	case "w":
		if len(strs) != 2 {
			return ErrInvalidArgsCount
		}

		w := Week{}
		err := w.parseInfo(strs[1])
		if err != nil {
			return fmt.Errorf("week parse info %w", err)
		}
		r.data = w

	case "m":
		m := Month{}
		switch len(strs) {
		case 2:
			err := m.parseDaysInfo(strs[1])
			if err != nil {
				return fmt.Errorf("month parse info: %w", err)
			}
			m.monthsEmpty = true
			r.data = m
		case 3:
			err := m.parseDaysInfo(strs[1])
			if err != nil {
				return fmt.Errorf("month parse info: %w", err)
			}

			err = m.parseMonthInfo(strs[2])
			if err != nil {
				return fmt.Errorf("month parse info: %w", err)
			}

			r.data = m

		default:
			return ErrInvalidArgsCount
		}

	default:
		return ErrCommandNotFound
	}

	return nil
}

type Day struct {
	day int
}

func (d *Day) parseInfo(str string) error {
	day, err := strconv.Atoi(str)
	if err != nil {
		return err
	}

	if day < minDayCount || day > maxDayCount {
		return fmt.Errorf("%w: day must be between %d and %d",
			ErrArgumentOutOfRange, minDayCount, maxDayCount)
	}

	d.day = day

	return nil
}

type Week struct {
	weekDays [7]bool
}

func (w *Week) parseInfo(weekInfo string) error {
	weekDays := [7]bool{}
	for _, weekStr := range strings.Split(weekInfo, ",") {
		weekDay, err := strconv.Atoi(weekStr)
		if err != nil {
			return err
		}

		if weekDay < weekLowerBound || weekDay > weekUpperBound {
			return fmt.Errorf("%w: week day must be between %d and %d",
				ErrArgumentOutOfRange, weekLowerBound, weekUpperBound)
		}

		if weekDay == 7 {
			if weekDays[0] {
				return ErrDuplicateValue
			}
			weekDays[0] = true
			continue
		}

		if weekDays[weekDay] {
			return ErrDuplicateValue
		}
		weekDays[weekDay] = true
	}

	w.weekDays = weekDays

	return nil
}

type Month struct {
	days        [34]bool
	months      [12]bool
	monthsEmpty bool
}

func (m *Month) parseDaysInfo(daysInfo string) error {
	days := [34]bool{}
	for _, dayStr := range strings.Split(daysInfo, ",") {
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return err
		}

		if day == -1 {
			if days[32] {
				return ErrDuplicateValue
			}
			days[32] = true
			continue
		}
		if day == -2 {
			if days[33] {
				return ErrDuplicateValue
			}
			days[33] = true
			continue
		}

		if day < dayLowerBound || day > dayUpperBound {
			return fmt.Errorf("%w: day must be between %d and %d",
				ErrArgumentOutOfRange, dayLowerBound, dayUpperBound)
		}

		if days[day-1] {
			return ErrDuplicateValue
		}
		days[day-1] = true
	}

	m.days = days

	return nil
}

func (m *Month) parseMonthInfo(monthsInfo string) error {
	months := [12]bool{}
	for _, monthStr := range strings.Split(monthsInfo, ",") {
		month, err := strconv.Atoi(monthStr)
		if err != nil {
			return err
		}

		if month < monthLowerBound || month > monthUpperBound {
			return fmt.Errorf("%w: month must be between %d and %d",
				ErrArgumentOutOfRange, monthLowerBound, monthUpperBound)
		}

		if months[month-1] {
			return ErrDuplicateValue
		}
		months[month-1] = true
	}

	m.months = months

	return nil
}

type Year struct{}
