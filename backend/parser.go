package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// what operation we're doing?
const (
	add = iota
	sub
	mul
	div
)

// source
// when invoking updateDT, what is the "source of truth"
// to update other fields
const (
	ts     = iota // dt.ts
	ymdhms        // dt.year, month, day, hour, minute & second firld
	ds
)

type parser struct {
	regex          string
	noOfParameters int
	parserFunc     func(match []string, dt *datetime)
	parseNext      bool
}

// datetime kind
const (
	none      = 0
	timestamp = 1 << (iota - 1)
	duration
	number
)

type datetime struct {
	kind                          int
	parameter                     string
	dt                            time.Time
	ts                            int64 // no of seconds
	day, month, year              int64
	hour, minute, second          int64
	days, hours, minutes, seconds float32
}

// Helpers

func Atof(f string) float32 {
	if s, err := strconv.ParseFloat(f, 32); err == nil {
		return float32(s)
	} else {
		return 0.0
	}
}

func Atoi(f string) int64 {
	if s, err := strconv.ParseInt(f, 10, 64); err == nil {
		return s
	} else {
		return 0
	}
}

func (dt *datetime) updateDT(source int) {

	var s int64

	if source == ymdhms {
		s = int64(dt.second + dt.minute*60 + dt.hour*3600 + dt.day*24*3600)
		dt.ts = s
	} else if source == ts {
		s = dt.ts
	} else if source == ds {
		// do nothing
		return
	}

	// calculating durations in units as below
	dt.seconds = float32(s)
	dt.minutes = float32(s) / 60.0
	dt.hours = float32(s) / 3600.0
	dt.days = float32(s) / (3600.0 * 24.0)

	// Finally align and rollover if needed
	// i.e. minute = 65 will become hour = 1 and minute = 5

	n := dt.ts
	dt.second = n
	dt.day = n / (24 * 3600)
	n %= (24 * 3600)
	dt.hour = n / 3600
	n %= 3600
	dt.minute = n / 60
	dt.second %= 60

	dt.dt = time.Date(int(dt.year), time.Month(dt.month), int(dt.day), int(dt.hour), int(dt.minute), int(dt.second), 0, time.UTC)
}

func (dt *datetime) calculateDT(dt1 datetime, dt2 datetime, operation int) {
	if operation == add {

		if dt1.kind == dt2.kind {
			dt.kind = dt1.kind
			dt.ts = dt1.ts + dt2.ts
		} else if (dt1.kind & duration) == (dt2.kind & duration) {
			dt.kind = duration
			dt.ts = dt1.ts + dt2.ts
		} else if (dt1.kind&timestamp != 0) && (dt2.kind&duration != 0) {
			dt.kind = timestamp
			dt.dt = dt1.dt
			dt.dt = dt.dt.Add(time.Second * time.Duration(dt2.second))
			dt.dt = dt.dt.Add(time.Minute * time.Duration(dt2.minute))
			dt.dt = dt.dt.Add(time.Hour * time.Duration(dt2.hour))
			dt.dt = dt.dt.Add(time.Hour * time.Duration(dt2.day*24))
			return
		}
	} else if operation == sub {
		dt.ts = dt1.ts - dt2.ts
		if (dt1.kind & duration) == (dt2.kind & duration) {
			dt.kind = duration
		}
	} else if operation == mul {
		dt.ts = dt1.ts * dt2.ts
		if (dt1.kind & duration) == (dt2.kind & number) {
			dt.kind = duration
		}
	} else if operation == div {
		if dt1.ts != 0 {
			dt.ts = dt1.ts / dt2.ts
		} else {
			dt.ts = dt1.ts
		}

		// fmt.Printf("Kinds: dt1 %d, dt2 %d\n", dt1.kind, dt2.kind)

		// 60 / 15s ->  4 (number) ??
		// (number) / (duration) = (number)

		if (dt1.kind&number != 0) && (dt2.kind&number != 0) {
			// 60 / 15 ->  4 (number) ??
			// (number) / (number) = (number)
			dt.kind = number
			// fmt.Printf("(2) Kind set to : %d\n", dt.kind)
		} else if (dt1.kind&duration != 0) && (dt2.kind&number != 0) {
			// 1h / 4   -> 15m
			// (duration) / (number) = (duration)
			dt.kind = duration
			// fmt.Printf("(3) Kind set to : %d\n", dt.kind)
		} else if (dt1.kind&duration != 0) && (dt2.kind&duration != 0) {
			// 1h / 15m -> 4 (number)
			// (duration) / (duration) = (number)
			dt.kind = number
			// fmt.Printf("(1) Kind set to : %d\n", dt.kind)
		}

		// if (dt1.kind&number == 0) && (dt2.kind&duration != 0) {
		// 	dt.kind = duration
		// 	fmt.Printf("(1) Kind set to : %d\n", dt.kind)
		// } else if (dt1.kind&duration != 0) && (dt2.kind&number != 0) {
		// 	dt.kind = duration
		// 	fmt.Printf("(2) Kind set to : %d\n", dt.kind)
		// } else if (dt1.kind&duration != 0) && (dt2.kind&duration != 0) {
		// 	dt.kind = number
		// 	fmt.Printf("(3) Kind set to : %d\n", dt.kind)
		// }
	}

	dt.updateDT(ts)

	// dt.dt = time.Now()
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
		//   - `<ss+>`
		{
			regex:          `^([0-9]+)$`,
			noOfParameters: 1,
			parserFunc: func(match []string, dt *datetime) {
				dt.second = Atoi(match[1])
				dt.kind = number | duration

				// updateDT needs to calculate:
				// - ts, days, hours, minutes, seconds
				dt.updateDT(ymdhms)
			},
		},
		//   - `<mm:ss>`
		{
			regex:          `^([0-9]+):([0-9]+)$`,
			noOfParameters: 2,
			parserFunc: func(match []string, dt *datetime) {
				dt.minute = Atoi(match[1])
				dt.second = Atoi(match[2])
				dt.kind = duration

				dt.updateDT(ymdhms)
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
				dt.kind = duration

				dt.updateDT(ymdhms)
			},
		},
		{
			regex:          `^([0-9]+)u$`,
			noOfParameters: 1,
			parserFunc: func(match []string, dt *datetime) {
				i := Atoi(match[1])
				dt.dt = time.Unix(i, 0)
				dt.kind = timestamp
			},
		},
		{
			regex:          `now`,
			noOfParameters: 0,
			parserFunc: func(match []string, dt *datetime) {
				dt.dt = time.Now().Round(time.Second)
				dt.kind = timestamp
			},
		},
		// Passers below needs to be ad the end
		// to support fields like 1d1h1s
		//   - `<d+>d`
		{
			regex:          `([0-9]+)d`,
			noOfParameters: 1,
			parserFunc: func(match []string, dt *datetime) {
				dt.day += Atoi(match[1])
				dt.kind = duration

				dt.updateDT(ymdhms)
			},
			parseNext: true,
		},
		//   - `<h+>h`
		{
			regex:          `([0-9]+)h`,
			noOfParameters: 1,
			parserFunc: func(match []string, dt *datetime) {
				dt.hour += Atoi(match[1])
				dt.kind = duration

				dt.updateDT(ymdhms)
			},
			parseNext: true,
		},
		//   - `<m+>m`
		{
			regex:          `([0-9]+)m`,
			noOfParameters: 1,
			parserFunc: func(match []string, dt *datetime) {
				dt.minute += Atoi(match[1])
				dt.kind = duration

				dt.updateDT(ymdhms)
			},
			parseNext: true,
		},
		//   - `<s+>s`
		{
			regex:          `([0-9]+)s`,
			noOfParameters: 1,
			parserFunc: func(match []string, dt *datetime) {
				dt.second += Atoi(match[1])
				dt.kind = duration

				dt.updateDT(ymdhms)
			},
			parseNext: true,
		},
	}

	found := false

	for _, p := range parsers {
		re := regexp.MustCompile(p.regex)
		match := re.FindStringSubmatch(f)

		if len(match) == p.noOfParameters+1 {
			// parserFunc will get only:
			//  - day, hour, minute, second
			p.parserFunc(match, dt)

			if !p.parseNext {
				return nil
			}
			found = true
		}
	}
	// we're at the end of the list
	// Have we found something? If not - this seems like a nerror
	if found {
		return nil
	} else {
		return errors.New("<none>")
	}
}

func parse(p string) (datetime, error) {
	// As later input string will be split using space(s)
	// ensure there IS at least one space around the operator
	// So accept: 3+4 or 3 + 4
	for _, c := range []string{"+", "-", "*", "/"} {
		p = strings.ReplaceAll(p, c, " "+c+" ")
	}

	fields := strings.Fields(p)

	switch len(fields) {
	case 0:
		result := datetime{
			parameter: p,
		}
		return result, errors.New("nothing to calculate")
	case 1:
		result := datetime{
			kind:      none,
			parameter: p,
		}

		if err := parseField(fields[0], &result); err == nil {
			return result, nil
		}
	case 2:
		result := datetime{
			parameter: p,
		}
		return result, errors.New("missing operator")
	case 3:
		dt1 := datetime{
			parameter: p,
		}

		if err := parseField(fields[0], &dt1); err != nil {
			return dt1, nil
		}

		dt2 := datetime{
			parameter: p,
		}

		if err := parseField(fields[2], &dt2); err != nil {
			return dt2, nil
		}

		result := datetime{
			parameter: p,
		}

		if fields[1] == "+" {
			result.calculateDT(dt1, dt2, add)
		} else if fields[1] == "-" {
			result.calculateDT(dt1, dt2, sub)
		} else if fields[1] == "*" {
			result.calculateDT(dt1, dt2, mul)
		} else if fields[1] == "/" {
			result.calculateDT(dt1, dt2, div)
		} else {
			return result, errors.New("allowed format: <field> <op> <field> where op = +-")
		}

		return result, nil
	}

	result := datetime{
		parameter: p,
	}
	return result, errors.New("input formatted incorrectly")
}
