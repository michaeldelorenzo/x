package datetime_test

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/michaeldelorenzo/x/v2/pkg/types/datetime"
	"github.com/stretchr/testify/require"
)

func TestDuration_JSON(t *testing.T) {
	t.Run("successfully marshals Duration into a slice of bytes", func(t *testing.T) {
		d := datetime.Duration("P1Y-10M3DT-30H20M4.4S")
		out, err := json.Marshal(d)

		require.NoError(t, err)
		require.EqualValues(t, string(out), strconv.Quote(string(d)))
	})

	t.Run("successfully unmarshals into the Duration type", func(t *testing.T) {
		var tests = []struct {
			desc  string
			input string
		}{
			{desc: "given a valid full ISO8601 duration string", input: "P1Y-10M3DT-30H20M4.4S"},
			{desc: "given a date-only ISO8601 duration string", input: "P1Y-10M3D"},
			{desc: "given a time-only ISO8601 duration string", input: "PT-4H30M"},
			{desc: "given the zero ISO8601 duration string", input: "PT0S"},
		}

		for _, test := range tests {
			t.Run(test.desc, func(t *testing.T) {
				d := new(datetime.Duration)
				err := json.Unmarshal([]byte(strconv.Quote(test.input)), &d)

				require.NoError(t, err)
				require.EqualValues(t, *d, test.input)
			})
		}
	})

	t.Run("returns an error when unmarshalling invalid duration strings", func(t *testing.T) {
		var tests = []struct {
			desc  string
			input string
		}{
			{desc: "given an invalid ISO8601 duration string", input: "P20"},
			{desc: "given just `P` with no values", input: "P"},
			{desc: "given just `PT` with no values", input: "PT"},
			{desc: "given an otherwise valid ISO8601 duration string that has a trailing T", input: "P32DT"},
		}

		for _, test := range tests {
			t.Run(test.desc, func(t *testing.T) {
				d := new(datetime.Duration)
				err := json.Unmarshal([]byte(strconv.Quote(test.input)), &d)
				require.Error(t, err)
			})
		}
	})

	t.Run("returns an error when unmarshalling non-string", func(t *testing.T) {
		d := new(datetime.Duration)
		err := json.Unmarshal([]byte(`42`), &d)
		require.Error(t, err)
	})
}

func TestDuration_Graphql(t *testing.T) {
	t.Run("successfully marshals Duration into a slice of bytes", func(t *testing.T) {
		testStringBuilder := &strings.Builder{}
		d := datetime.Duration("P1Y-10M3DT-30H20M4.4S")
		d.MarshalGQL(testStringBuilder)
		require.EqualValues(t, testStringBuilder.String(), strconv.Quote(string(d)))
	})

	t.Run("successfully unmarshals into the Duration type", func(t *testing.T) {
		var tests = []struct {
			desc  string
			input string
		}{
			{desc: "given a valid full ISO8601 duration string", input: "P1Y-10M3DT-30H20M4.4S"},
			{desc: "given a date-only ISO8601 duration string", input: "P1Y-10M3D"},
			{desc: "given a time-only ISO8601 duration string", input: "PT-4H30M"},
			{desc: "given the zero ISO8601 duration string", input: "PT0S"},
		}

		for _, test := range tests {
			t.Run(test.desc, func(t *testing.T) {
				d := new(datetime.Duration)
				err := d.UnmarshalGQL(test.input)
				require.NoError(t, err)
				require.EqualValues(t, *d, test.input)
			})
		}
	})

	t.Run("returns an error when unmarshalling invalid duration strings", func(t *testing.T) {
		var tests = []struct {
			desc  string
			input string
		}{
			{desc: "given an invalid ISO8601 duration string", input: "P20"},
			{desc: "given just `P` with no values", input: "P"},
			{desc: "given just `PT` with no values", input: "PT"},
			{desc: "given an otherwise valid ISO8601 duration string that has a trailing T", input: "P32DT"},
		}

		for _, test := range tests {
			t.Run(test.desc, func(t *testing.T) {
				d := new(datetime.Duration)
				err := d.UnmarshalGQL(test.input)
				require.Error(t, err)
			})
		}
	})

	t.Run("returns an error when unmarshalling non-string", func(t *testing.T) {
		d := new(datetime.Duration)
		err := d.UnmarshalGQL(32)
		require.Error(t, err)
	})
}

func TestDuration_ValidatePos(t *testing.T) {
	t.Run("returns true if duration string is positive only", func(t *testing.T) {
		d := "P1Y10M3DT30H20M4.4S"
		output := datetime.ValidatePositiveDuration(d)
		require.True(t, output)
	})

	t.Run("returns false if duration string has negative components", func(t *testing.T) {
		d := "P1Y-10M3DT-30H20M4.4S"
		output := datetime.ValidatePositiveDuration(d)
		require.False(t, output)
	})
}
