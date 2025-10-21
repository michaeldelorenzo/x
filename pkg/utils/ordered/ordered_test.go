package ordered_test

import (
	"testing"

	"github.com/michaeldelorenzo/x/v3/pkg/utils/ordered"
	"github.com/stretchr/testify/require"
)

func TestOrdered_Min(t *testing.T) {
	t.Run("returns the min of two floats", func(t *testing.T) {
		// given
		low := 3.5
		high := 8.1

		// when
		output1 := ordered.Min(high, low)
		output2 := ordered.Min(low, high)

		// then
		require.Equal(t, low, output1)
		require.Equal(t, low, output2)
	})

	t.Run("returns the min of two ints", func(t *testing.T) {
		// given
		low := 3
		high := 8

		// when
		output1 := ordered.Min(high, low)
		output2 := ordered.Min(low, high)

		// then
		require.Equal(t, low, output1)
		require.Equal(t, low, output2)
	})

	t.Run("returns the min of two strings", func(t *testing.T) {
		// given
		low := "2022-05-23T12:00:00Z"
		high := "2022-06-22T04:00:00Z"

		// when
		output1 := ordered.Min(high, low)
		output2 := ordered.Min(low, high)

		// then
		require.Equal(t, low, output1)
		require.Equal(t, low, output2)
	})

}

func TestOrdered_Max(t *testing.T) {
	t.Run("returns the max of two floats", func(t *testing.T) {
		// given
		low := 3.5
		high := 8.1

		// when
		output1 := ordered.Max(high, low)
		output2 := ordered.Max(low, high)

		// then
		require.Equal(t, high, output1)
		require.Equal(t, high, output2)
	})

	t.Run("returns the max of two ints", func(t *testing.T) {
		// given
		low := 3
		high := 8

		// when
		output1 := ordered.Max(high, low)
		output2 := ordered.Max(low, high)

		// then
		require.Equal(t, high, output1)
		require.Equal(t, high, output2)
	})

	t.Run("returns the max of two strings", func(t *testing.T) {
		// given
		low := "2022-05-23T12:00:00Z"
		high := "2022-06-22T04:00:00Z"

		// when
		output1 := ordered.Max(high, low)
		output2 := ordered.Max(low, high)

		// then
		require.Equal(t, high, output1)
		require.Equal(t, high, output2)
	})

}
