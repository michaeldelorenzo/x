package maputils_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/michaeldelorenzo/x/v2/pkg/types/jsonblob"
	"github.com/michaeldelorenzo/x/v2/pkg/utils/maputils"
	"github.com/stretchr/testify/require"
)

func TestKeys(t *testing.T) {
	t.Run("iterates through the map and outputs its keys", func(t *testing.T) {
		// given
		original := map[string]string{"abc": "22", "doerayme": "howdy"}
		expected := []string{"abc", "doerayme"}

		// when
		output := maputils.Keys(original)

		// then
		require.ElementsMatch(t, expected, output)
	})

	t.Run("returns empty slice for empty map", func(t *testing.T) {
		// given
		original := map[string]string{}
		expected := []string{}

		// when
		output := maputils.Keys(original)

		// then
		require.ElementsMatch(t, expected, output)
	})
}

func TestValues(t *testing.T) {
	t.Run("iterates through the map and outputs its values", func(t *testing.T) {
		// given
		original := map[string]string{"abc": "22", "doerayme": "howdy"}
		expected := []string{"22", "howdy"}

		// when
		output := maputils.Values(original)

		// then
		require.ElementsMatch(t, expected, output)
	})

	t.Run("returns empty slice for empty map", func(t *testing.T) {
		// given
		original := map[string]string{}
		expected := []string{}

		// when
		output := maputils.Values(original)

		// then
		require.EqualValues(t, expected, output)
	})
}

func TestMap(t *testing.T) {
	t.Run("iterates through the original map and outputs a transformed list", func(t *testing.T) {
		// given
		original := map[string]string{"abc": "22", "doerayme": "howdy"}
		expected := []string{"abc-22", "doerayme-howdy"}

		// when
		output := maputils.Map(original, func(k, v string) string {
			return fmt.Sprintf("%s-%s", k, v)
		})

		// then
		require.ElementsMatch(t, expected, output)
	})

	t.Run("iterates through the original and outputs a transformed list of a different type", func(t *testing.T) {
		// given
		original := map[string]string{"abc": "22", "doerayme": "howdy"}
		expected := []int{5, 13}

		// when
		output := maputils.Map(original, func(k, v string) int {
			return len(k) + len(v)
		})

		// then
		require.ElementsMatch(t, expected, output)
	})

	t.Run("returns an empty array of the proper output type without calling the callback when given an empty slice", func(t *testing.T) {
		// given
		original := map[string]string{}
		expected := []int{}
		callbackHasBeenCalled := false

		// when
		output := maputils.Map(original, func(k, v string) int {
			callbackHasBeenCalled = true
			return len(k) + len(v)
		})

		// then
		require.False(t, callbackHasBeenCalled)
		require.ElementsMatch(t, expected, output)
	})
}

func TestReduce(t *testing.T) {
	t.Run("reduces through the map to build a filtered slice", func(t *testing.T) {
		// given
		original := map[string]string{"abc": "22", "doerayme": "howdy", "anotherOne": "booyah"}
		expected := []string{"22", "booyah"}

		// when
		output := maputils.Reduce(original, []string{}, func(acc []string, k, v string) []string {
			if strings.HasPrefix(k, "a") {
				return append(acc, v)
			}

			return acc
		})

		// then
		require.ElementsMatch(t, expected, output)
	})

	t.Run("reduces through the map and outputs a newly created map", func(t *testing.T) {
		// given
		original := map[string]string{"abc": "22", "doerayme": "howdy", "anotherOne": "booyah"}
		expected := map[string]int{
			"abc":        2,
			"doerayme":   5,
			"anotherOne": 6,
		}

		// when
		output := maputils.Reduce(original, map[string]int{}, func(acc map[string]int, k, v string) map[string]int {
			acc[k] = len(v)
			return acc
		})

		// then
		require.EqualValues(t, expected, output)
	})

	t.Run("returns the initial value without calling the callback when given an empty slice", func(t *testing.T) {
		// given
		original := map[string]string{}
		expected := map[string]int{}
		callbackHasBeenCalled := false

		// when
		output := maputils.Reduce(original, map[string]int{}, func(acc map[string]int, k, v string) map[string]int {
			callbackHasBeenCalled = true
			acc[k] = len(v)
			return acc
		})

		// then
		require.False(t, callbackHasBeenCalled)
		require.EqualValues(t, expected, output)
	})
}

func TestGetAndCoalesce(t *testing.T) {
	t.Run("returns the value from the blob with the expected type", func(t *testing.T) {
		// given
		var expectedType float64
		blob := jsonblob.Blob{
			"c": "doerayme",
			"a": 45.67,
			"b": map[string]interface{}{"yes": "hey", "no": "way"},
			"1": []interface{}{
				map[string]interface{}{"z": "oh", "x": "yes"},
				map[string]interface{}{"y": true},
			},
		}

		// when
		value, ok := maputils.GetAndCoalesce[float64](blob, "a")

		// then
		require.NotEmpty(t, value)
		require.Equal(t, 45.67, value)
		require.IsType(t, expectedType, value)
		require.True(t, ok)
	})

	t.Run("returns the zero value when key is not found in map", func(t *testing.T) {
		// given
		var expectedType float64
		blob := jsonblob.Blob{
			"c": "doerayme",
			"a": 45.67,
			"b": map[string]interface{}{"yes": "hey", "no": "way"},
			"1": []interface{}{
				map[string]interface{}{"z": "oh", "x": "yes"},
				map[string]interface{}{"y": true},
			},
		}

		// when
		value, ok := maputils.GetAndCoalesce[float64](blob, "no-key")

		// then
		require.Empty(t, value)
		require.IsType(t, expectedType, value)
		require.False(t, ok)
	})

	t.Run("returns the zero value when key-value pair is of different type", func(t *testing.T) {
		// given
		var expectedType []string
		blob := jsonblob.Blob{
			"c": "doerayme",
			"a": 45.67,
			"b": map[string]interface{}{"yes": "hey", "no": "way"},
			"1": []interface{}{
				map[string]interface{}{"z": "oh", "x": "yes"},
				map[string]interface{}{"y": true},
			},
		}

		// when
		value, ok := maputils.GetAndCoalesce[[]string](blob, "1")

		// then
		require.Empty(t, value)
		require.IsType(t, expectedType, value)
		require.False(t, ok)
	})
}
