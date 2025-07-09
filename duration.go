package durago

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"unicode"
)

const (
	stateParsePeriod = iota
	stateParseTime

	nsPerSecond = 1000000000
	nsPerMinute = nsPerSecond * 60
	nsPerHour   = nsPerMinute * 60

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
)

var (
	ErrInvalidFormat = errors.New("invalid format")
	ErrParse         = errors.New("parse failed")
)

type Duration struct {
	d        time.Duration
	negative bool
}

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

			return duration, nil
		default:
			if unicode.IsNumber(char) || char == floatDesignator {
				num = append(num, char)
				continue
			}

			return nil, fmt.Errorf("%w: unexpected value or designator", ErrInvalidFormat)
		}
	}

	return duration, nil
}

func (d *Duration) GetTimeDuration() time.Duration {
	if d.negative {
		return -d.d
	}

	return d.d
}
