package datetime_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/michaeldelorenzo/x/pkg/types/datetime"
	"github.com/stretchr/testify/require"
)

func TestWallTime_JSON(t *testing.T) {
	t.Run("successfully JSON marshals time into string", func(t *testing.T) {
		expected := "13:32:21"
		d := datetime.AtoWallTime(expected)
		out, err := d.MarshalJSON()

		require.NoError(t, err)
		require.EqualValues(t, strconv.Quote(expected), out)
	})

	t.Run("successfully JSON unmarshals into the WallTime type", func(t *testing.T) {
		input := "22:32:00"
		d := new(datetime.WallTime)
		err := d.UnmarshalJSON([]byte(strconv.Quote(input)))

		require.NoError(t, err)
		require.EqualValues(t, input, d.String())
	})

	t.Run("returns error without panicking when unmarshalling into nil WallTime", func(t *testing.T) {
		input := "22:32:00"
		var d *datetime.WallTime
		err := d.UnmarshalJSON([]byte(strconv.Quote(input)))

		require.Error(t, err)
	})

	t.Run("returns an error when JSON unmarshalling invalid time string", func(t *testing.T) {
		var tests = []struct {
			desc  string
			input string
		}{
			{desc: "given an invalid ISO8601 time string", input: `"48:99:00"`},
			{desc: "given a full ISO8601 date string", input: `"2022-04-20T00:00:00Z"`},
			{desc: "given a non string", input: `42`},
		}
		for _, test := range tests {
			t.Run(test.desc, func(t *testing.T) {
				d := new(datetime.WallTime)
				err := d.UnmarshalJSON([]byte(test.input))

				require.Error(t, err)
			})
		}
	})
}

func TestWallTime_GraphQL(t *testing.T) {

	t.Run("successfully GraphQL marshals time into string", func(t *testing.T) {
		expected := "09:32:21"
		testStringBuilder := &strings.Builder{}
		d := datetime.AtoWallTime(expected)
		d.MarshalGQL(testStringBuilder)

		require.EqualValues(t, strconv.Quote(expected), testStringBuilder.String())
	})

	t.Run("successfully GraphQL unmarshals into the WallTime type", func(t *testing.T) {
		input := "00:32:00"
		d := new(datetime.WallTime)
		err := d.UnmarshalGQL(input)

		require.NoError(t, err)
		require.EqualValues(t, input, d.String())
	})

	t.Run("returns error without panicking when unmarshalling into nil WallTime", func(t *testing.T) {
		input := "22:32:00"
		var d *datetime.WallTime
		err := d.UnmarshalGQL(input)

		require.Error(t, err)
	})

	t.Run("returns an error when GraphQL unmarshalling invalid time string", func(t *testing.T) {
		var tests = []struct {
			desc  string
			input interface{}
		}{
			{desc: "given an invalid ISO8601 time string", input: "48:99:00"},
			{desc: "given a full ISO8601 date string", input: "2022-04-20T00:00:00Z"},
			{desc: "given a non string", input: 42},
		}
		for _, test := range tests {
			t.Run(test.desc, func(t *testing.T) {
				d := new(datetime.WallTime)
				err := d.UnmarshalGQL(test.input)

				require.Error(t, err)
			})
		}
	})
}

func TestWallTime_SQL(t *testing.T) {

	t.Run("successfully implements driver.Valuer to return ISO 8601 time string", func(t *testing.T) {
		expected := "13:32:21"
		d := datetime.AtoWallTime(expected)

		out, err := d.Value()

		require.NoError(t, err)
		require.EqualValues(t, expected, out)
	})

	t.Run("successfully implements sql.Scanner to scan time.Time into WallTime", func(t *testing.T) {
		expected := "09:32:21"
		input, _ := time.Parse(datetime.TimeFormat, expected)
		output := new(datetime.WallTime)

		err := output.Scan(input)

		require.NoError(t, err)
		require.EqualValues(t, expected, output.String())
	})

	t.Run("returns error without panicking when scanning into nil WallTime", func(t *testing.T) {
		expected := "09:32:21"
		input, _ := time.Parse(datetime.TimeFormat, expected)
		var output *datetime.WallTime

		err := output.Scan(input)

		require.Error(t, err)
	})

	t.Run("returns error without panicking when scanning unexpected type into WallTime", func(t *testing.T) {
		output := new(datetime.WallTime)

		err := output.Scan(42)

		require.Error(t, err)
	})
}

func TestWallTime_UnmarshalText(t *testing.T) {
	expected := time.Date(0, 1, 1, 13, 47, 19, 0, time.UTC)
	input := []byte("13:47:19")

	output := datetime.WallTime{}
	err := output.UnmarshalText(input)

	require.NoError(t, err)
	require.Equal(t, expected, output.Time)
}
