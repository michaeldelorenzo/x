package iter_test

import (
	"strings"
	"testing"

	"github.com/michaeldelorenzo/x/pkg/utils/iter"
	"github.com/michaeldelorenzo/x/pkg/utils/ordered"
	"github.com/stretchr/testify/require"
)

func TestIter_Map(t *testing.T) {
	t.Run("iterates through the original and outputs the transformed list", func(t *testing.T) {
		// given
		original := []string{"abc", "doerayme"}
		expected := []string{"ABC", "DOERAYME"}

		// when
		output := iter.Map(original, func(str string) string {
			return strings.ToUpper(str)
		})

		// then
		require.EqualValues(t, expected, output)
	})

	t.Run("iterates through the original and outputs a transformed list of a different type", func(t *testing.T) {
		// given
		original := []string{"abc", "doerayme"}
		expected := []int{3, 8}

		// when
		output := iter.Map(original, func(str string) int {
			return len(str)
		})

		// then
		require.EqualValues(t, expected, output)
	})

	t.Run("returns an empty array of the proper output type without calling the callback when given an empty slice", func(t *testing.T) {
		// given
		original := []string{}
		expected := []int{}
		callbackHasBeenCalled := false

		// when
		output := iter.Map(original, func(str string) int {
			callbackHasBeenCalled = true
			return len(str)
		})

		// then
		require.False(t, callbackHasBeenCalled)
		require.EqualValues(t, expected, output)
	})
}

func TestIter_Reduce(t *testing.T) {
	t.Run("reduces through the slice and outputs the maximum length found", func(t *testing.T) {
		// given
		original := []string{"abc", "abcd", "abcdef", "ab"}
		expected := 6

		// when
		output := iter.Reduce(original, 0, func(acc int, el string) int {
			return ordered.Max(len(el), acc)
		})

		// then
		require.EqualValues(t, expected, output)
	})

	t.Run("reduces through the slice and outputs a newly created map", func(t *testing.T) {
		// given
		original := []string{"abc", "abcd", "abcdef", "ab"}
		expected := map[string]int{
			"abc":    3,
			"abcd":   4,
			"abcdef": 6,
			"ab":     2,
		}

		// when
		output := iter.Reduce(original, map[string]int{}, func(acc map[string]int, el string) map[string]int {
			acc[el] = len(el)
			return acc
		})

		// then
		require.EqualValues(t, expected, output)
	})

	t.Run("returns the initial value without calling the callback when given an empty slice", func(t *testing.T) {
		// given
		original := []string{}
		expected := map[string]int{}
		callbackHasBeenCalled := false

		// when
		output := iter.Reduce(original, map[string]int{}, func(acc map[string]int, el string) map[string]int {
			callbackHasBeenCalled = true
			acc[el] = len(el)
			return acc
		})

		// then
		require.False(t, callbackHasBeenCalled)
		require.EqualValues(t, expected, output)
	})
}
