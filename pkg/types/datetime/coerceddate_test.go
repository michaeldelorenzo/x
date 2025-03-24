package datetime_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/koneksahealth/x/pkg/types/datetime"
	"github.com/stretchr/testify/require"
)

func TestCoercedDate_JSON(t *testing.T) {
	t.Run("successfully JSON marshals CoercedDate into string", func(t *testing.T) {
		expected := "2022-03-22"
		d := datetime.CoercedDate{Date: datetime.AtoDate(expected)}
		out, err := d.MarshalJSON()

		require.NoError(t, err)
		require.EqualValues(t, strconv.Quote(expected), out)
	})

	t.Run("successfully JSON unmarshals into the CoercedDate type", func(t *testing.T) {
		input := "2022-03-22"
		d := new(datetime.CoercedDate)
		err := d.UnmarshalJSON([]byte(strconv.Quote(input)))

		require.NoError(t, err)
		require.EqualValues(t, input, d.String())

		longInput := "2022-08-23T17:48:48.000Z"
		expected := "2022-08-23"
		err = d.UnmarshalJSON([]byte(strconv.Quote(longInput)))

		require.NoError(t, err)
		require.EqualValues(t, expected, d.String())
	})

	t.Run("returns error without panicking when unmarshalling into nil CoercedDate", func(t *testing.T) {
		input := "2022-03-22"
		var d *datetime.CoercedDate
		err := d.UnmarshalJSON([]byte(strconv.Quote(input)))

		require.Error(t, err)
	})

	t.Run("returns an error when JSON unmarshalling invalid time string", func(t *testing.T) {
		var tests = []struct {
			desc  string
			input string
		}{
			{desc: "given an invalid ISO8601 time string", input: `"2022-99-01"`},
			{desc: "given a date string with less than 10 characters", input: `"2022-4-2"`},
			{desc: "given a non string", input: `42`},
		}
		for _, test := range tests {
			t.Run(test.desc, func(t *testing.T) {
				d := new(datetime.CoercedDate)
				err := d.UnmarshalJSON([]byte(test.input))

				require.Error(t, err)
			})
		}
	})
}

func TestCoercedDate_GraphQL(t *testing.T) {

	t.Run("successfully GraphQL marshals CoercedDate into string", func(t *testing.T) {
		expected := "2022-03-22"
		testStringBuilder := &strings.Builder{}
		d := datetime.CoercedDate{Date: datetime.AtoDate(expected)}
		d.MarshalGQL(testStringBuilder)

		require.EqualValues(t, strconv.Quote(expected), testStringBuilder.String())
	})

	t.Run("successfully GraphQL unmarshals into the CoercedDate type", func(t *testing.T) {
		input := "2022-03-22"
		d := new(datetime.CoercedDate)
		err := d.UnmarshalGQL(input)

		require.NoError(t, err)
		require.EqualValues(t, input, d.String())

		longInput := "2022-08-23T17:48:48.000Z"
		expected := "2022-08-23"
		err = d.UnmarshalGQL(longInput)

		require.NoError(t, err)
		require.EqualValues(t, expected, d.String())
	})

	t.Run("returns error without panicking when unmarshalling into nil CoercedDate", func(t *testing.T) {
		input := "2022-03-22"
		var d *datetime.CoercedDate
		err := d.UnmarshalGQL(input)

		require.Error(t, err)
	})

	t.Run("returns an error when GraphQL unmarshalling invalid time string", func(t *testing.T) {
		var tests = []struct {
			desc  string
			input interface{}
		}{
			{desc: "given an invalid ISO8601 time string", input: "2022-99-01"},
			{desc: "given a date string with less than 10 characters", input: "2022-4-2"},
			{desc: "given a non string", input: 42},
		}
		for _, test := range tests {
			t.Run(test.desc, func(t *testing.T) {
				d := new(datetime.CoercedDate)
				err := d.UnmarshalGQL(test.input)

				require.Error(t, err)
			})
		}
	})
}

func TestCoercedDate_SQL(t *testing.T) {

	t.Run("successfully implements driver.Valuer to return ISO 8601 date string", func(t *testing.T) {
		expected := "2022-03-22"
		d := datetime.CoercedDate{Date: datetime.AtoDate(expected)}

		out, err := d.Value()

		require.NoError(t, err)
		require.EqualValues(t, expected, out)
	})

	t.Run("successfully implements sql.Scanner to scan time.Time into Date", func(t *testing.T) {
		expected := "2022-03-22"
		input := time.Date(2022, 3, 22, 0, 0, 0, 0, time.UTC)
		output := new(datetime.CoercedDate)

		err := output.Scan(input)

		require.NoError(t, err)
		require.EqualValues(t, expected, output.String())

		longInput := time.Date(2022, 8, 23, 17, 48, 48, 0, time.UTC)
		expected = "2022-08-23"

		err = output.Scan(longInput)

		require.NoError(t, err)
		require.EqualValues(t, expected, output.String())
	})

	t.Run("returns error without panicking when scanning into nil Date", func(t *testing.T) {
		expected := "2022-03-22"
		input := datetime.AtoDate(expected).Time
		var output *datetime.CoercedDate

		err := output.Scan(input)

		require.Error(t, err)
	})

	t.Run("returns error without panicking when scanning unexpected type into Date", func(t *testing.T) {
		output := new(datetime.CoercedDate)

		err := output.Scan(42)

		require.Error(t, err)
	})
}

func TestCoercedDate_UnmarshalText(t *testing.T) {
	expected := time.Date(2022, 6, 14, 0, 0, 0, 0, time.UTC)
	input := []byte("2022-06-14")

	output := datetime.CoercedDate{}
	err := output.UnmarshalText(input)

	require.NoError(t, err)
	require.Equal(t, expected, output.Time)

	expected = time.Date(2022, 8, 23, 0, 0, 0, 0, time.UTC)
	input = []byte("2022-08-23T17:48:48.000Z")

	err = output.UnmarshalText(input)

	require.NoError(t, err)
	require.Equal(t, expected, output.Time)
}
