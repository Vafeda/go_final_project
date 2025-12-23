package nextdate

import "errors"

var (
	ErrCommandNotFound    = errors.New("command not found")
	ErrInvalidArgsCount   = errors.New("invalid number of arguments")
	ErrArgumentOutOfRange = errors.New("argument out of range")
	ErrDuplicateValue     = errors.New("duplicate value detected")
)

const (
	minDayCount = 1
	maxDayCount = 400

	dayLowerBound = 1
	dayUpperBound = 31

	weekLowerBound = 1
	weekUpperBound = 7

	monthLowerBound = 1
	monthUpperBound = 12

	dateFormat = "20060102"
)
