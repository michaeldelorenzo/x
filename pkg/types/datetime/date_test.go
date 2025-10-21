package datetime_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/michaeldelorenzo/x/v2/pkg/types/datetime"
	"github.com/stretchr/testify/require"
)

func TestDate_JSON(t *testing.T) {
	t.Run("successfully JSON marshals date into string", func(t *testing.T) {
		expected := "2022-03-22"
		d := datetime.AtoDate(expected)
		out, err := d.MarshalJSON()

		require.NoError(t, err)
		require.EqualValues(t, strconv.Quote(expected), out)
	})

	t.Run("successfully JSON unmarshals into the Date type", func(t *testing.T) {
		input := "2022-03-22"
		d := new(datetime.Date)
		err := d.UnmarshalJSON([]byte(strconv.Quote(input)))

		require.NoError(t, err)
		require.EqualValues(t, input, d.String())
	})

	t.Run("returns error without panicking when unmarshalling into nil Date", func(t *testing.T) {
		input := "2022-03-22"
		var d *datetime.Date
		err := d.UnmarshalJSON([]byte(strconv.Quote(input)))

		require.Error(t, err)
	})

	t.Run("returns an error when JSON unmarshalling invalid time string", func(t *testing.T) {
		var tests = []struct {
			desc  string
			input string
		}{
			{desc: "given an invalid ISO8601 time string", input: `"2022-99-01"`},
			{desc: "given a full ISO8601 date string", input: `"2022-04-20T00:00:00Z"`},
			{desc: "given a non string", input: `42`},
		}
		for _, test := range tests {
			t.Run(test.desc, func(t *testing.T) {
				d := new(datetime.Date)
				err := d.UnmarshalJSON([]byte(test.input))

				require.Error(t, err)
			})
		}
	})
}

func TestDate_GraphQL(t *testing.T) {

	t.Run("successfully GraphQL marshals date into string", func(t *testing.T) {
		expected := "2022-03-22"
		testStringBuilder := &strings.Builder{}
		d := datetime.AtoDate(expected)
		d.MarshalGQL(testStringBuilder)

		require.EqualValues(t, strconv.Quote(expected), testStringBuilder.String())
	})

	t.Run("successfully GraphQL unmarshals into the Date type", func(t *testing.T) {
		input := "2022-03-22"
		d := new(datetime.Date)
		err := d.UnmarshalGQL(input)

		require.NoError(t, err)
		require.EqualValues(t, input, d.String())
	})

	t.Run("returns error without panicking when unmarshalling into nil Date", func(t *testing.T) {
		input := "2022-03-22"
		var d *datetime.Date
		err := d.UnmarshalGQL(input)

		require.Error(t, err)
	})

	t.Run("returns an error when GraphQL unmarshalling invalid time string", func(t *testing.T) {
		var tests = []struct {
			desc  string
			input interface{}
		}{
			{desc: "given an invalid ISO8601 time string", input: "2022-99-01"},
			{desc: "given a full ISO8601 date string", input: "2022-04-20T00:00:00Z"},
			{desc: "given a non string", input: 42},
		}
		for _, test := range tests {
			t.Run(test.desc, func(t *testing.T) {
				d := new(datetime.Date)
				err := d.UnmarshalGQL(test.input)

				require.Error(t, err)
			})
		}
	})
}

func TestDate_SQL(t *testing.T) {

	t.Run("successfully implements driver.Valuer to return ISO 8601 date string", func(t *testing.T) {
		expected := "2022-03-22"
		d := datetime.AtoDate(expected)

		out, err := d.Value()

		require.NoError(t, err)
		require.EqualValues(t, expected, out)
	})

	t.Run("successfully implements sql.Scanner to scan time.Time into Date", func(t *testing.T) {
		expected := "2022-03-22"
		input := datetime.AtoDate(expected).Time
		output := new(datetime.Date)

		err := output.Scan(input)

		require.NoError(t, err)
		require.EqualValues(t, expected, output.String())
	})

	t.Run("returns error without panicking when scanning into nil Date", func(t *testing.T) {
		expected := "2022-03-22"
		input := datetime.AtoDate(expected).Time
		var output *datetime.Date

		err := output.Scan(input)

		require.Error(t, err)
	})

	t.Run("returns error without panicking when scanning unexpected type into Date", func(t *testing.T) {
		output := new(datetime.Date)

		err := output.Scan(42)

		require.Error(t, err)
	})
}

func TestDate_UnmarshalText(t *testing.T) {
	t.Run("successfully unmarhsals text from an ISO 8601 date string", func(t *testing.T) {
		expected := time.Date(2022, 6, 14, 0, 0, 0, 0, time.UTC)
		input := []byte("2022-06-14")

		output := datetime.Date{}
		err := output.UnmarshalText(input)

		require.NoError(t, err)
		require.Equal(t, expected, output.Time)
	})

	t.Run("returns error without panicking when scanning unexpected type into Date", func(t *testing.T) {
		input := []byte("2022-06-1n")
		output := datetime.Date{}
		err := output.UnmarshalText(input)

		require.Error(t, err)
	})
}

func TestDate_DaysSince_and_DaysUntil(t *testing.T) {
	nyTZ, err := time.LoadLocation("America/New_York")

	require.NoError(t, err, "failed to load America/New_York time zone for test setup")

	testCases := []struct {
		name     string
		start    datetime.Date
		end      datetime.Date
		expected int
	}{
		{
			name:     "computes days properly",
			start:    datetime.AtoDate("2022-03-22"),
			end:      datetime.AtoDate("2022-03-25"),
			expected: 3,
		},
		{
			name:     "returns 0 days when start equals end",
			start:    datetime.AtoDate("2022-03-22"),
			end:      datetime.AtoDate("2022-03-22"),
			expected: 0,
		},
		{
			name:     "computes days properly when month changes",
			start:    datetime.AtoDate("2022-03-22"),
			end:      datetime.AtoDate("2022-04-22"),
			expected: 31,
		},
		{
			name:     "returns negative days when start is before end",
			start:    datetime.AtoDate("2022-03-22"),
			end:      datetime.AtoDate("2022-03-18"),
			expected: -4,
		},
		{
			name:     "properly handles leap years",
			start:    datetime.AtoDate("2024-02-01"),
			end:      datetime.AtoDate("2024-03-01"),
			expected: 29,
		},
		{
			name: "properly ignores daylight savings",
			// purposefully manually create the dates below rather than using TtoDate to test time zone correction
			start: datetime.Date{Time: time.Date(2024, 3, 10, 0, 0, 0, 0, nyTZ)},
			end:   datetime.Date{Time: time.Date(2024, 3, 11, 0, 0, 0, 0, nyTZ)},
			// should still return 1 even though in NY this only contains 23 hours
			expected: 1,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s calculating days since", testCase.name), func(t *testing.T) {
			require.Equal(t, testCase.expected, testCase.end.DaysSince(testCase.start))
		})
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s calculating days until", testCase.name), func(t *testing.T) {
			require.Equal(t, -testCase.expected, testCase.end.DaysUntil(testCase.start))
		})
	}
}

func Test_TtoDate(t *testing.T) {
	expected := time.Date(2022, 6, 14, 0, 0, 0, 0, time.UTC)

	output := datetime.TtoDate(expected)

	require.Equal(t, expected, output.Time)
}

func Test_Today(t *testing.T) {
	expected := time.Now().UTC().Format(datetime.DateFormat)

	output := datetime.Today().String()

	require.Equal(t, expected, output)
}

func Test_AddDays_and_AddDate(t *testing.T) {
	nyTZ, err := time.LoadLocation("America/New_York")
	require.NoError(t, err, "failed to load America/New_York time zone for test setup")

	testCases := []struct {
		name     string
		origin   datetime.Date
		expected datetime.Date
		offset   int
	}{
		{
			name:     "offsets days properly",
			origin:   datetime.AtoDate("2022-03-22"),
			expected: datetime.AtoDate("2022-03-25"),
			offset:   3,
		},
		{
			name:     "returns same day when offset equals 0",
			origin:   datetime.AtoDate("2022-03-22"),
			expected: datetime.AtoDate("2022-03-22"),
			offset:   0,
		},
		{
			name:     "computes days properly when month changes",
			origin:   datetime.AtoDate("2022-03-22"),
			expected: datetime.AtoDate("2022-04-22"),
			offset:   31,
		},
		{
			name:     "performs inverted calc when days is negative",
			origin:   datetime.AtoDate("2022-03-22"),
			expected: datetime.AtoDate("2022-03-18"),
			offset:   -4,
		},
		{
			name:     "properly handles leap years",
			origin:   datetime.AtoDate("2024-02-01"),
			expected: datetime.AtoDate("2024-03-01"),
			offset:   29,
		},
		{
			name:     "properly ignores daylight savings",
			origin:   datetime.Date{Time: time.Date(2024, 3, 10, 0, 0, 0, 0, nyTZ)},
			expected: datetime.Date{Time: time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC)},
			// should still return 1 even though in NY this only contains 23 hours
			offset: 1,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s performing day addition", testCase.name), func(t *testing.T) {
			require.Equal(t, testCase.expected, testCase.origin.AddDays(testCase.offset))
		})
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s performing date addition", testCase.name), func(t *testing.T) {
			require.Equal(t, testCase.expected, testCase.origin.AddDate(0, 0, testCase.offset))
		})
	}
}
