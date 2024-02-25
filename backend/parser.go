package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	add = iota
	sub
)

type datetime struct {
	parameter                     string
	dt                            time.Time
	ts                            int // no of seconds
	day, month, year              int
	hour, minute, second          int
	days, hours, minutes, seconds float32
}

func updateDT(dt *datetime) {
	s := float32(dt.second + dt.minute*60 + dt.hour*3600)

	dt.seconds = s
	dt.minutes = s / 60.0
	dt.hours = s / 3600.0
	dt.days = s / (3600.0 * 24.0)

	dt.ts = int(s)
	dt.dt = time.Date(dt.year, time.Month(dt.month), dt.day, dt.hour, dt.minute, dt.second, 0, time.UTC)
}

// Try to guess field format.
//
// Can be any of:
//
// Time component formats `<time>`:
//   - `<ss>`
//   - `<mm:ss>`
//   - `<hh:mm:ss>`
//
// Date component formats `<date>`:
//   - If configured `DD/MM/YYYY`
//   - `<DD>/<MM>`
//   - `<DD>/<MM>/<YYYY>`
//   - If configured `MM/DD/YYYY`
//   - `<MM>/<DD>`
//   - `<MM>/<DD>/<YYYY>`
//
// Compount duration component `<period>`:
//   - `<d>d<h>h<m>m<s>s`
//   - Any component can be ommited, e.g. `1d4h`
func parseField(f string, dt *datetime) error {
	parsers := []parser{
		//   - `<ss>`
		{
			regex:          `^([0-9]+)$`,
			noOfParameters: 1,
			parserFunc: func(match []string, dt *datetime) {
				dt.second = Atoi(match[1])
			},
		},
		//   - `<mm:ss>`
		{
			regex:          `^([0-9]+):([0-9]+)$`,
			noOfParameters: 2,
			parserFunc: func(match []string, dt *datetime) {
				dt.minute = Atoi(match[1])
				dt.second = Atoi(match[2])
			},
		},
		//   - `<hh:mm:ss>`
		{
			regex:          `^([0-9]+):([0-9]+):([0-9]+)$`,
			noOfParameters: 3,
			parserFunc: func(match []string, dt *datetime) {
				dt.hour = Atoi(match[1])
				dt.minute = Atoi(match[2])
				dt.second = Atoi(match[3])
			},
		},
	}

	for _, p := range parsers {
		re := regexp.MustCompile(p.regex)
		match := re.FindStringSubmatch(f)

		if len(match) == p.noOfParameters+1 {
			p.parserFunc(match, dt)
			updateDT(dt)
			return nil
		}
	}

	return errors.New("<none>")
}

func (r *datetime) calculateDT(dt1 datetime, dt2 datetime, operation int) {
	if operation == add {
		r.ts = dt2.ts + dt1.ts
	} else if operation == sub {
		r.ts = dt2.ts - dt1.ts
	}

	n := r.ts
	r.second = n % 60
	n /= 60
	r.minute = n % 60
	n /= 60
	r.hour = n % 60
	n /= 60
	r.day = n % 24
	n /= 24

	updateDT(r)

	r.dt = time.Now()
}

func Atof(f string) float32 {
	if s, err := strconv.ParseFloat(f, 32); err == nil {
		return float32(s)
	} else {
		return 0.0
	}
}

func Atoi(f string) int {
	if s, err := strconv.Atoi(f); err == nil {
		return s
	} else {
		return 0.0
	}
}

type parser struct {
	regex          string
	noOfParameters int
	parserFunc     func(match []string, dt *datetime)
}

func parse(p string) string {
	fields := strings.Split(p, " ")

	switch len(fields) {
	case 0:
		result := datetime{
			parameter: p,
		}
		return getMarshalledItems(result)
	case 1:
		result := datetime{
			parameter: p,
		}

		if err := parseField(fields[0], &result); err == nil {
			return getMarshalledItems(result)
		}
	case 2:
		dt1 := datetime{
			parameter: p,
		}

		if err := parseField(fields[0], &dt1); err != nil {
			return getMarshalledItems(dt1)
		}

		dt2 := datetime{
			parameter: p,
		}

		if err := parseField(fields[1], &dt2); err != nil {
			return getMarshalledItems(dt2)
		}

		result := datetime{
			parameter: p,
		}

		result.calculateDT(dt2, dt1, sub)
		return getMarshalledItems(result)
	case 3:
		dt1 := datetime{
			parameter: p,
		}

		if err := parseField(fields[0], &dt1); err != nil {
			return getMarshalledItems(dt1)
		}

		dt2 := datetime{
			parameter: p,
		}

		if err := parseField(fields[2], &dt2); err != nil {
			return getMarshalledItems(dt2)
		}

		result := datetime{
			parameter: p,
		}

		if fields[1] == "+" {
			result.calculateDT(dt2, dt1, add)
		} else if fields[1] == "-" {
			result.calculateDT(dt2, dt1, sub)
		}

		return getMarshalledItems(result)
	}

	result := datetime{
		parameter: p,
	}
	return getMarshalledItems(result)
}
