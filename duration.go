package durago

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	stateParsePeriod = iota
	stateParseTime

	nsPerMillisecond = 1000000
	nsPerSecond      = 1000 * nsPerMillisecond
	nsPerMinute      = nsPerSecond * 60
	nsPerHour        = nsPerMinute * 60

	periodDay   = nsPerHour * 24
	periodWeek  = periodDay * 7
	periodMonth = periodYear / 12
	periodYear  = periodDay * 365

	secondDesignator      = 'S'
	minuteMonthDesignator = 'M'
	hourDesignator        = 'H'
	timeDesignator        = 'T'
	dayDesignator         = 'D'
	weekDesignator        = 'W'
	yearDesignator        = 'Y'
	durationDesignator    = 'P'

	positiveSign    = '+'
	negativeSign    = '-'
	floatDesignator = '.'

	zeroDuration = "PT0S"
)

var (
	ErrInvalidFormat = errors.New("invalid format")
	ErrParse         = errors.New("parse failed")
)

type Duration struct {
	d        time.Duration
	negative bool

	years   int
	months  int
	weeks   int
	days    int
	hours   int
	minutes int
	seconds float64
}

// ParseDuration attempts to parse the given duration string into a *Duration,
// if parsing fails an error is returned instead.
func ParseDuration(d string) (*Duration, error) {
	state := stateParsePeriod

	duration := &Duration{}
	num := make([]rune, 0, 4)
	parsedParts := []bool{
		false, // sign
		false, // duration
		false, // year
		false, // month
		false, // week
		false, // day
		false, // time
		false, // hour
		false, // minute
		false, // second
	}

	for _, char := range d {
		switch char {
		case positiveSign:
			if state != stateParsePeriod || parsedParts[0] {
				return nil, fmt.Errorf("%w: unexpected positive sign", ErrInvalidFormat)
			}

			parsedParts[0] = true
		case negativeSign:
			if state != stateParsePeriod || parsedParts[0] {
				return nil, fmt.Errorf("%w: unexpected negative sign", ErrInvalidFormat)
			}

			parsedParts[0] = true
			duration.negative = true
		case durationDesignator:
			if state != stateParsePeriod || parsedParts[1] {
				return nil, fmt.Errorf("%w: unexpected duration designator", ErrInvalidFormat)
			}
			parsedParts[1] = true
		case yearDesignator:
			if state != stateParsePeriod || parsedParts[2] {
				return nil, fmt.Errorf("%w: unexpected year designator", ErrInvalidFormat)
			}

			years, err := strconv.ParseInt(string(num), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("year %w: %s", ErrParse, err.Error())
			}

			parsedParts[2] = true
			num = num[:0]
			duration.d += time.Duration(years * periodYear)
			duration.years = int(years)
		case minuteMonthDesignator:
			if state == stateParsePeriod {
				if parsedParts[3] {
					return nil, fmt.Errorf("%w: unexpected month designator", ErrInvalidFormat)
				}

				months, err := strconv.ParseInt(string(num), 10, 64)
				if err != nil {
					return nil, fmt.Errorf("month %w: %s", ErrParse, err.Error())
				}

				parsedParts[3] = true
				num = num[:0]
				duration.d += time.Duration(months * periodMonth)
				duration.months = int(months)
				continue
			}

			if parsedParts[8] {
				return nil, fmt.Errorf("%w: unexpected minute designator", ErrInvalidFormat)
			}

			minutes, err := strconv.ParseInt(string(num), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("month %w: %s", ErrParse, err.Error())
			}

			parsedParts[8] = true
			num = num[:0]
			duration.d += time.Duration(minutes * nsPerMinute)
			duration.minutes = int(minutes)
		case weekDesignator:
			if state != stateParsePeriod || parsedParts[4] {
				return nil, fmt.Errorf("%w: unexpected week designator", ErrInvalidFormat)
			}

			weeks, err := strconv.ParseInt(string(num), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("week %w: %s", ErrParse, err.Error())
			}

			parsedParts[4] = true
			num = num[:0]
			duration.d += time.Duration(weeks * periodWeek)
			duration.weeks = int(weeks)
		case dayDesignator:
			if state != stateParsePeriod || parsedParts[5] {
				return nil, fmt.Errorf("%w: unexpected day designator", ErrInvalidFormat)
			}

			days, err := strconv.ParseInt(string(num), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("day %w: %s", ErrParse, err.Error())
			}

			parsedParts[5] = true
			num = num[:0]
			duration.d += time.Duration(days * periodDay)
			duration.days = int(days)
		case timeDesignator:
			if state != stateParsePeriod || parsedParts[6] {
				return nil, fmt.Errorf("%w: unexpected time designator", ErrInvalidFormat)
			}

			parsedParts[6] = true
			state = stateParseTime
		case hourDesignator:
			if state != stateParseTime || parsedParts[7] {
				return nil, fmt.Errorf("%w: unexpected hour designator", ErrInvalidFormat)
			}

			hours, err := strconv.ParseInt(string(num), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("hour %w: %s", ErrParse, err.Error())
			}

			parsedParts[7] = true
			num = num[:0]
			duration.d += time.Duration(hours * nsPerHour)
			duration.hours = int(hours)
		case secondDesignator:
			if state != stateParseTime || parsedParts[9] {
				return nil, fmt.Errorf("%w: unexpected second designator", ErrInvalidFormat)
			}

			seconds, err := strconv.ParseFloat(string(num), 64)
			if err != nil {
				return nil, fmt.Errorf("second %w: %s", ErrParse, err.Error())
			}

			parsedParts[9] = true
			duration.d += time.Duration(seconds * nsPerSecond)
			duration.seconds = seconds

			return duration, nil
		default:
			if unicode.IsNumber(char) || char == floatDesignator {
				num = append(num, char)
				continue
			}

			return nil, fmt.Errorf("%w: unexpected value or designator", ErrInvalidFormat)
		}
	}

	if len(num) > 0 {
		return nil, fmt.Errorf("%w: missing designator", ErrInvalidFormat)
	}

	return duration, nil
}

// GetTimeDuration returns underlying tim.Duration with corresponding sign
func (d *Duration) GetTimeDuration() time.Duration {
	if d.negative {
		return -d.d
	}

	return d.d
}

// FromTimeDuration converts the given time.Duration into durago.Duration.
func FromTimeDuration(d time.Duration) *Duration {
	duration := &Duration{}

	if d == 0 {
		return duration
	}

	if d < 0 {
		duration.negative = true
		d = -d
	}

	duration.d = d

	for d >= periodYear {
		duration.years++
		d -= periodYear
	}

	for d >= periodMonth {
		duration.months++
		d -= periodMonth
	}

	for d >= periodWeek {
		duration.weeks++
		d -= periodWeek
	}

	for d >= periodDay {
		duration.days++
		d -= periodDay
	}

	for d >= nsPerHour {
		duration.hours++
		d -= nsPerHour
	}

	for d >= nsPerMinute {
		duration.minutes++
		d -= nsPerMinute
	}

	duration.seconds = d.Seconds()

	return duration
}

// String returns the ISO8601 duration string for the *Duration
func (d *Duration) String() string {
	if d.d == 0 {
		return zeroDuration
	}

	var (
		b       strings.Builder
		hasTime bool
	)

	b.Grow(20)

	if d.negative {
		b.WriteString(string(negativeSign))
	}

	b.WriteString(string(durationDesignator))

	if d.years != 0 {
		b.WriteString(strconv.Itoa(d.years))
		b.WriteString(string(yearDesignator))
	}

	if d.months != 0 {
		b.WriteString(strconv.Itoa(d.months))
		b.WriteString(string(minuteMonthDesignator))
	}

	if d.weeks != 0 {
		b.WriteString(strconv.Itoa(d.weeks))
		b.WriteString(string(weekDesignator))
	}

	if d.days != 0 {
		b.WriteString(strconv.Itoa(d.days))
		b.WriteString(string(dayDesignator))
	}

	if d.hours != 0 {
		b.WriteString(string(timeDesignator))
		b.WriteString(strconv.Itoa(d.hours))
		b.WriteString(string(hourDesignator))
		hasTime = true
	}

	if d.minutes != 0 {
		if !hasTime {
			b.WriteString(string(timeDesignator))
			hasTime = true
		}
		b.WriteString(strconv.Itoa(d.minutes))
		b.WriteString(string(minuteMonthDesignator))
	}

	if d.seconds != 0 {
		if !hasTime {
			b.WriteString(string(timeDesignator))
			hasTime = true
		}
		b.WriteString(strconv.FormatFloat(d.seconds, 'f', -1, 64))
		b.WriteString(string(secondDesignator))
	}

	return b.String()
}

// MarshalJSON satisfies the Marshaler interface by return a valid JSON string representation of the duration
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON satisfies the Unmarshaler interface by return a valid JSON string representation of the duration
func (d *Duration) UnmarshalJSON(source []byte) error {
	var duration string
	if err := json.Unmarshal(source, &duration); err != nil {
		return err
	}

	parsed, err := ParseDuration(duration)
	if err != nil {
		return err
	}

	*d = *parsed
	return nil
}
