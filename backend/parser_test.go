package main

import (
	"testing"
	"time"
)

// Try to guess field format.
//
// Can be any of:
//
// Time component formats `<time>`:
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

func TestParseField(t *testing.T) {
	tests := []struct {
		input    string
		expected datetime
	}{
		{
			//   - `<ss>`
			input: "12",
			expected: datetime{
				ts:      0*24 + 0*3600 + 0*60 + 12,
				day:     0,
				month:   0,
				year:    0,
				hour:    0,
				minute:  0,
				second:  12,
				days:    12.0 / 3600.0 / 24.0,
				hours:   12.0 / 3600.0,
				minutes: 12.0 / 60.0,
				seconds: 12.0,
			},
		},
		{
			//   - `<ss>`
			input: "60",
			expected: datetime{
				ts:      0*24 + 0*3600 + 1*60 + 0,
				day:     0,
				month:   0,
				year:    0,
				hour:    0,
				minute:  1,
				second:  0,
				days:    60.0 / 3600.0 / 24.0,
				hours:   60.0 / 3600.0,
				minutes: 60.0 / 60.0,
				seconds: 60.0,
			},
		},
		{
			//   - `<ss+>`
			input: "123",
			expected: datetime{
				ts:      0*24 + 0*3600 + 0*60 + 123,
				day:     0,
				month:   0,
				year:    0,
				hour:    0,
				minute:  2,
				second:  3,
				days:    123.0 / 3600.0 / 24.0,
				hours:   123.0 / 3600.0,
				minutes: 123.0 / 60.0,
				seconds: 123.0,
			},
		},
		{
			//   - `<mm:ss>`
			input: "12:34",
			expected: datetime{
				ts:      0*24 + 0*3600 + 12*60 + 34,
				day:     0,
				month:   0,
				year:    0,
				hour:    0,
				minute:  12,
				second:  34,
				days:    (12*60.0 + 34.00) / 3600.0 / 24.0,
				hours:   (12*60.0 + 34.00) / 3600.0,
				minutes: (12*60.0 + 34.00) / 60.0,
				seconds: (12*60.0 + 34.00),
			},
		},
		{
			//   - `<hh:mm:ss>`
			input: "12:34:56",
			expected: datetime{
				ts:      0*24 + 12*3600 + 34*60 + 56,
				day:     0,
				month:   0,
				year:    0,
				hour:    12,
				minute:  34,
				second:  56,
				days:    (12*3600 + 34*60.0 + 56.00) / 3600.0 / 24.0,
				hours:   (12*3600 + 34*60.0 + 56.00) / 3600.0,
				minutes: (12*3600 + 34*60.0 + 56.00) / 60.0,
				seconds: (12*3600 + 34*60.0 + 56.00),
			},
		},
		{
			//   - `<hh:mm:ss>`
			input: "1d",
			expected: datetime{
				ts:      0*24 + 0*3600 + 0*60 + 0 + 1*24*3600,
				day:     1,
				month:   0,
				year:    0,
				hour:    0,
				minute:  0,
				second:  0,
				days:    (0*3600 + 0*60.0 + 0.00 + 1.0*24.0*3600.0) / 3600.0 / 24.0,
				hours:   (0*3600 + 0*60.0 + 0.00 + 1.0*24.0*3600.0) / 3600.0,
				minutes: (0*3600 + 0*60.0 + 0.00 + 1.0*24.0*3600.0) / 60.0,
				seconds: (0*3600 + 0*60.0 + 0.00 + 1.0*24.0*3600.0),
			},
		},
	}

	for _, ts := range tests {
		ts.expected.dt = time.Date(ts.expected.year,
			time.Month(ts.expected.month),
			ts.expected.day,
			ts.expected.hour,
			ts.expected.minute,
			ts.expected.second,
			0,
			time.UTC)

		result := datetime{}

		parseField(ts.input, &result)

		if result != ts.expected {
			t.Logf(">>> Result NOK for input: >%s<\n", ts.input)
			t.Error(">>> Input", ts.input)
			t.Errorf(">>> Expected %+v\n", ts.expected)
			t.Errorf(">>> Result %+v\n", result)
		}
	}

}

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected datetime
	}{
		{
			//   - `<ss>`
			input: "10 + 2",
			expected: datetime{
				ts:      0*24 + 0*3600 + 0*60 + 12,
				day:     0,
				month:   0,
				year:    0,
				hour:    0,
				minute:  0,
				second:  12,
				days:    12.0 / 3600.0 / 24.0,
				hours:   12.0 / 3600.0,
				minutes: 12.0 / 60.0,
				seconds: 12.0,
			},
		},
		{
			//   - `<ss+>`
			input: "200 77",
			expected: datetime{
				ts:      0*24 + 0*3600 + 0*60 + 123,
				day:     0,
				month:   0,
				year:    0,
				hour:    0,
				minute:  2,
				second:  3,
				days:    123.0 / 3600.0 / 24.0,
				hours:   123.0 / 3600.0,
				minutes: 123.0 / 60.0,
				seconds: 123.0,
			},
		},
		{
			//   - `<mm:ss>`
			input: "10:34 + 2:26",
			expected: datetime{
				ts:      0*24 + 0*3600 + 13*60 + 00,
				day:     0,
				month:   0,
				year:    0,
				hour:    0,
				minute:  13,
				second:  00,
				days:    (13*60.0 + 0.00) / 3600.0 / 24.0,
				hours:   (13*60.0 + 0.00) / 3600.0,
				minutes: (13*60.0 + 0.00) / 60.0,
				seconds: (13*60.0 + 0.00),
			},
		},
		{
			input: "12:34:56 - 12:34:56",
			expected: datetime{
				ts:      0*24 + 0*3600 + 0*60 + 0,
				day:     0,
				month:   0,
				year:    0,
				hour:    0,
				minute:  0,
				second:  0,
				days:    (0*3600 + 0*60.0 + 0.00) / 3600.0 / 24.0,
				hours:   (0*3600 + 0*60.0 + 0.00) / 3600.0,
				minutes: (0*3600 + 0*60.0 + 0.00) / 60.0,
				seconds: (0*3600 + 0*60.0 + 0.00),
			},
		},
		{
			input: "1d2h3m4s - 4s3m2h1d",
			expected: datetime{
				ts:      0*24 + 0*3600 + 0*60 + 0,
				day:     0,
				month:   0,
				year:    0,
				hour:    0,
				minute:  0,
				second:  0,
				days:    (0*3600 + 0*60.0 + 0.00) / 3600.0 / 24.0,
				hours:   (0*3600 + 0*60.0 + 0.00) / 3600.0,
				minutes: (0*3600 + 0*60.0 + 0.00) / 60.0,
				seconds: (0*3600 + 0*60.0 + 0.00),
			},
		},
		{
			input: "1d",
			expected: datetime{
				ts:      0*24 + 0*3600 + 0*60 + 0 + 1*24*3600,
				day:     1,
				month:   0,
				year:    0,
				hour:    0,
				minute:  0,
				second:  0,
				days:    (0*3600 + 0*60.0 + 0.00 + 1.0*24.0*3600.0) / 3600.0 / 24.0,
				hours:   (0*3600 + 0*60.0 + 0.00 + 1.0*24.0*3600.0) / 3600.0,
				minutes: (0*3600 + 0*60.0 + 0.00 + 1.0*24.0*3600.0) / 60.0,
				seconds: (0*3600 + 0*60.0 + 0.00 + 1.0*24.0*3600.0),
			},
		},
		{
			input: "1d + 12",
			expected: datetime{
				ts:      1*24*3600 + 0*3600 + 0*60 + 12,
				day:     1,
				month:   0,
				year:    0,
				hour:    0,
				minute:  0,
				second:  12,
				days:    86412.0 / 3600.0 / 24.0,
				hours:   86412.0 / 3600.0,
				minutes: 86412.0 / 60.0,
				seconds: 86412.0,
			},
		},
		{
			input: "1d + 12s",
			expected: datetime{
				ts:      1*24*3600 + 0*3600 + 0*60 + 12,
				day:     1,
				month:   0,
				year:    0,
				hour:    0,
				minute:  0,
				second:  12,
				days:    86412.0 / 3600.0 / 24.0,
				hours:   86412.0 / 3600.0,
				minutes: 86412.0 / 60.0,
				seconds: 86412.0,
			},
		},
	}

	for _, ts := range tests {
		ts.expected.dt = time.Date(ts.expected.year,
			time.Month(ts.expected.month),
			ts.expected.day,
			ts.expected.hour,
			ts.expected.minute,
			ts.expected.second,
			0,
			time.UTC)

		ts.expected.parameter = ts.input
		result, _ := parse(ts.input)

		// ignore parameter modifications
		result.parameter = ts.input

		if result != ts.expected {
			t.Logf(">>> Result NOK for input: >%s<\n", ts.input)
			t.Error(">>> Input", ts.input)
			t.Errorf(">>> Expected %+v\n", ts.expected)
			t.Errorf(">>> Result %+v\n", result)
		}
	}

}
