package domain

import (
	"strconv"
	"strings"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

var (
	ErrInvalidHour = custom_errors.BadException("invalid hour")
	ErrInvalidMinute = custom_errors.BadException("invalid minute")
	ErrInvalidTimeFormat = custom_errors.BadException("invalid time format, expected HH:MM")
	ErrInvalidTimeConversion = func(message string) error {return custom_errors.BadException(message)}
)

type TimeOfDay struct {
	hour   int
	minute int
}

func NewTimeOfDay(hour int, minute int) (TimeOfDay, error) {
	if hour < 0 || hour > 23 {
		return TimeOfDay{}, ErrInvalidHour
	}
	if minute < 0 || minute > 59 {
		return TimeOfDay{}, ErrInvalidMinute
	}
	return TimeOfDay{
		hour:   hour,
		minute: minute,
	}, nil
}

func NewTimeOfDayFromString(timeStr string) (TimeOfDay, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 || len(parts[0]) != 2 || len(parts[1]) != 2 {
		return TimeOfDay{}, ErrInvalidTimeFormat
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return TimeOfDay{}, ErrInvalidTimeConversion(err.Error())
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return TimeOfDay{}, ErrInvalidTimeConversion(err.Error())
	}

	return NewTimeOfDay(hour, minute)
}

func (t TimeOfDay) MinutesSinceMidnight() int {
	return t.hour*60 + t.minute
}

func (t TimeOfDay) Before(other TimeOfDay) bool {
	return t.MinutesSinceMidnight() < other.MinutesSinceMidnight()
}

func (t TimeOfDay) After(other TimeOfDay) bool {
	return t.MinutesSinceMidnight() > other.MinutesSinceMidnight()
}

func (t TimeOfDay) Equal(other TimeOfDay) bool {
	return t.MinutesSinceMidnight() == other.MinutesSinceMidnight()
}

func (t TimeOfDay) Hour() int {
	return t.hour
}

func (t TimeOfDay) Minute() int {
	return t.minute
}