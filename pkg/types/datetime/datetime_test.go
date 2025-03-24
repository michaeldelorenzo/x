package datetime_test

import (
	"strconv"
	"testing"

	"github.com/koneksahealth/x/pkg/types/datetime"
	"github.com/stretchr/testify/require"
)

func TestDateTime_JSON(t *testing.T) {
	t.Run("successfully JSON marshals DateTime into string", func(t *testing.T) {
		expected := "2022-03-22T12:03:37.456Z"
		d := datetime.AtoDateTime(expected)
		out, err := d.MarshalJSON()

		require.NoError(t, err)
		require.EqualValues(t, strconv.Quote(expected), out)
	})

	t.Run("successfully JSON unmarshals into the DateTime type", func(t *testing.T) {
		input := "2022-03-22T12:03:37.456Z"
		d := new(datetime.DateTime)
		err := d.UnmarshalJSON([]byte(strconv.Quote(input)))

		require.NoError(t, err)
		require.EqualValues(t, input, d.String())
	})

	t.Run("returns error without panicking when unmarshalling into nil DateTime", func(t *testing.T) {
		input := "2022-03-22T12:03:37.456Z"
		var d *datetime.DateTime
		err := d.UnmarshalJSON([]byte(strconv.Quote(input)))

		require.Error(t, err)
	})

	t.Run("returns an error when JSON unmarshalling invalid time string", func(t *testing.T) {
		var tests = []struct {
			desc  string
			input string
		}{
			{desc: "given an invalid ISO8601 time string", input: `"2022-04-01"`},
			{desc: "given a full ISO8601 DateTime string", input: `"2022-04-20T99:00:00Z"`},
			{desc: "given a non string", input: `42`},
		}
		for _, test := range tests {
			t.Run(test.desc, func(t *testing.T) {
				d := new(datetime.DateTime)
				err := d.UnmarshalJSON([]byte(test.input))

				require.Error(t, err)
			})
		}
	})
}

func TestDateTime_SQL(t *testing.T) {

	t.Run("successfully implements driver.Valuer to return time.Time", func(t *testing.T) {
		expected := "2022-03-22T12:03:37.456Z"
		d := datetime.AtoDateTime(expected)

		out, err := d.Value()

		require.NoError(t, err)
		require.EqualValues(t, d.Time, out)
	})

	t.Run("successfully implements sql.Scanner to scan time.Time into DateTime", func(t *testing.T) {
		expected := "2022-03-22T12:03:37.456Z"
		input := datetime.AtoDateTime(expected).Time
		output := new(datetime.DateTime)

		err := output.Scan(input)

		require.NoError(t, err)
		require.EqualValues(t, expected, output.String())
	})

	t.Run("returns error without panicking when scanning into nil DateTime", func(t *testing.T) {
		expected := "2022-03-22T12:03:37.456Z"
		input := datetime.AtoDateTime(expected).Time
		var output *datetime.DateTime

		err := output.Scan(input)

		require.Error(t, err)
	})

	t.Run("returns error without panicking when scanning unexpected type into DateTime", func(t *testing.T) {
		output := new(datetime.DateTime)

		err := output.Scan(42)

		require.Error(t, err)
	})
}
