package sequence_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/michaeldelorenzo/x/pkg/utils/ordered"
	"github.com/michaeldelorenzo/x/pkg/utils/sequence"
	"github.com/stretchr/testify/require"
)

func TestContains(t *testing.T) {
	// testCase provides a structured generic test case
	type testCase[T comparable] struct {
		name  string
		input T
		seq   []T
		want  bool
	}
	//: Testing against a collection of numbers
	testNum := []testCase[int]{
		{
			name:  "should return whether the slice contains the provided integer (true)",
			input: 50,
			seq:   []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			want:  true,
		},
		{
			name:  "should return whether the slice contains the provided integer (false)",
			input: 1,
			seq:   []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			want:  false,
		},
	}

	//: Testing against a collection of strings
	testStr := []testCase[string]{
		{
			name:  "should return whether the slice contains the provided string (true)",
			input: "i",
			seq:   []string{"a", "e", "i", "o", "u"},
			want:  true,
		},
		{
			name:  "should return whether the slice contains the provided string (false)",
			input: "y",
			seq:   []string{"a", "e", "i", "o", "u"},
			want:  false,
		},
	}

	//: Testing against a collection of structs
	testAnon := []testCase[struct {
		name string
		age  int
	}]{
		{
			name: "should return whether the slice contains the provided struct (true)",
			input: struct {
				name string
				age  int
			}{
				name: "John Doe",
				age:  34,
			},
			seq: []struct {
				name string
				age  int
			}{
				{
					name: "John Doe",
					age:  34,
				},
				{
					name: "Jane Doe",
					age:  40,
				},
			},
			want: true,
		},
		{
			name: "should return whether the slice contains the provided struct (false)",
			input: struct {
				name string
				age  int
			}{
				name: "John Doe",
				age:  50,
			},
			seq: []struct {
				name string
				age  int
			}{
				{
					name: "John Doe",
					age:  34,
				},
				{
					name: "Jane Doe",
					age:  40,
				},
			},
			want: false,
		},
	}

	//: Running our tests iterations
	for _, test := range testNum {
		t.Run(test.name, func(t *testing.T) {
			got := sequence.Contains(test.seq, test.input)

			if got != test.want {
				t.Errorf("expected %t, but got %t", test.want, got)
			}
		})
	}

	for _, test := range testStr {
		t.Run(test.name, func(t *testing.T) {
			got := sequence.Contains(test.seq, test.input)

			if got != test.want {
				t.Errorf("expected %t, but got %t", test.want, got)
			}
		})
	}

	for _, test := range testAnon {
		t.Run(test.name, func(t *testing.T) {
			got := sequence.Contains(test.seq, test.input)

			if got != test.want {
				t.Errorf("expected %t, but got %t", test.want, got)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	type testCase[T comparable] struct {
		name string
		seq  []T
		want []T
	}
	//: Testing against a collection of numbers
	testNum := []testCase[int]{
		{
			name: "should return only the unique elements in the sequence of ints",
			seq:  []int{10, 20, 30, 100, 40, 50, 80, 60, 70, 80, 25, 90, 100},
			want: []int{10, 20, 30, 100, 40, 50, 80, 60, 70, 25, 90},
		},
		{
			name: "should return the entire sequence of ints if every element is unique",
			seq:  []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			want: []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		},
	}

	//: Testing against a collection of strings
	testStr := []testCase[string]{
		{
			name: "should return only the unique elements in the sequence of strings",
			seq:  []string{"a", "e", "i", "o", "a", "i", "o", "u"},
			want: []string{"a", "e", "i", "o", "u"},
		},
		{
			name: "should return the entire sequence of strings if every element is unique",
			seq:  []string{"a", "e", "i", "o", "u"},
			want: []string{"a", "e", "i", "o", "u"},
		},
	}

	type el struct {
		name string
		age  int
	}

	//: Testing against a collection of structs
	testStruct := []testCase[el]{
		{
			name: "should return only the unique elements in the sequence of structs",
			seq: []el{
				{
					name: "John Doe",
					age:  34,
				},
				{
					name: "Jane Doe",
					age:  40,
				},
				{
					name: "John Doe",
					age:  34,
				},
				{
					name: "Jane Doe",
					age:  40,
				},
				{
					name: "Remmy",
					age:  40,
				},
			},
			want: []el{
				{
					name: "John Doe",
					age:  34,
				},
				{
					name: "Jane Doe",
					age:  40,
				},
				{
					name: "Remmy",
					age:  40,
				},
			},
		},
		{
			name: "should return the entire sequence of structs if every element is unique",
			seq: []el{
				{
					name: "John Doe",
					age:  34,
				},
				{
					name: "Jane Doe",
					age:  40,
				},
				{
					name: "Remmy",
					age:  40,
				},
			},
			want: []el{
				{
					name: "John Doe",
					age:  34,
				},
				{
					name: "Jane Doe",
					age:  40,
				},
				{
					name: "Remmy",
					age:  40,
				},
			},
		},
	}

	//: Running our tests iterations
	for _, test := range testNum {
		t.Run(test.name, func(t *testing.T) {
			unique := sequence.Unique(test.seq)

			require.EqualValues(t, test.want, unique)
		})
	}

	for _, test := range testStr {
		t.Run(test.name, func(t *testing.T) {
			unique := sequence.Unique(test.seq)

			require.EqualValues(t, test.want, unique)
		})
	}

	for _, test := range testStruct {
		t.Run(test.name, func(t *testing.T) {
			unique := sequence.Unique(test.seq)

			require.EqualValues(t, test.want, unique)
		})
	}
}

func TestFlatten(t *testing.T) {
	t.Run("flattens a slice of slices", func(t *testing.T) {
		// given
		expected := []int{1, 4, 12, 3, 2, 3, 5, 8, 1, 2}
		input := [][]int{{1, 4, 12, 3}, {2, 3, 5, 8}, {1, 2}, {}}

		// when
		output := sequence.Flatten(input)

		// then
		require.EqualValues(t, expected, output)
	})
}

var flattenRes []int

func BenchmarkFlatten(b *testing.B) {
	var out []int
	for n := 0; n < b.N; n++ {
		out = sequence.Flatten(intMatrix)
	}
	flattenRes = out
}

func TestFilter(t *testing.T) {
	// testCase provides a structured generic test case
	type testCase[T comparable] struct {
		name     string
		filterFn func(T) bool
		seq      []T
		want     []T
	}
	//: Testing against a collection of numbers
	testNum := []testCase[int]{
		{
			name:     "should return a properly filtered slice of integers",
			filterFn: func(i int) bool { return i > 50 && i <= 90 },
			seq:      []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			want:     []int{60, 70, 80, 90},
		},
		{
			name:     "should return an empty slice of integers if no value passes filterFn",
			filterFn: func(i int) bool { return i > 100 },
			seq:      []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			want:     []int{},
		},
		{
			name:     "should return a copy of the input slice of integers if every value passes filterFn",
			filterFn: func(i int) bool { return i <= 100 },
			seq:      []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			want:     []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		},
	}

	//: Testing against a collection of strings
	testStr := []testCase[string]{
		{
			name:     "should return a properly filtered slice of strings",
			filterFn: func(i string) bool { return i == "a" || i == "o" },
			seq:      []string{"a", "e", "i", "o", "u"},
			want:     []string{"a", "o"},
		},
		{
			name:     "should return an empty slice of strings if no value passes filterFn",
			filterFn: func(i string) bool { return len(i) > 1 },
			seq:      []string{"a", "e", "i", "o", "u"},
			want:     []string{},
		},
		{
			name:     "should return a copy of the input slice of strings if every value passes filterFn",
			filterFn: func(i string) bool { return len(i) > 0 },
			seq:      []string{"a", "e", "i", "o", "u"},
			want:     []string{"a", "e", "i", "o", "u"},
		},
	}

	//: Testing against a collection of structs
	testAnon := []testCase[struct {
		name string
		age  int
	}]{
		{
			name: "should return a properly filtered slice of structs",
			filterFn: func(s struct {
				name string
				age  int
			}) bool {
				return s.age > 35
			},
			seq: []struct {
				name string
				age  int
			}{
				{
					name: "John Doe",
					age:  34,
				},
				{
					name: "Jane Doe",
					age:  40,
				},
			},
			want: []struct {
				name string
				age  int
			}{
				{
					name: "Jane Doe",
					age:  40,
				},
			},
		},
		{
			name: "should return an empty slice of structs if no value passes filterFn",
			filterFn: func(s struct {
				name string
				age  int
			}) bool {
				return s.age > 45
			},
			seq: []struct {
				name string
				age  int
			}{
				{
					name: "John Doe",
					age:  34,
				},
				{
					name: "Jane Doe",
					age:  40,
				},
			},
			want: []struct {
				name string
				age  int
			}{},
		},
		{
			name: "should return a copy of the input slice of structs if every value passes filterFn",
			filterFn: func(s struct {
				name string
				age  int
			}) bool {
				return s.age > 30
			},
			seq: []struct {
				name string
				age  int
			}{
				{
					name: "John Doe",
					age:  34,
				},
				{
					name: "Jane Doe",
					age:  40,
				},
			},
			want: []struct {
				name string
				age  int
			}{
				{
					name: "John Doe",
					age:  34,
				},
				{
					name: "Jane Doe",
					age:  40,
				},
			},
		},
	}

	//: Running our tests iterations
	for _, test := range testNum {
		t.Run(test.name, func(t *testing.T) {
			got := sequence.Filter(test.seq, test.filterFn)
			require.Equal(t, test.want, got)
		})
	}

	for _, test := range testStr {
		t.Run(test.name, func(t *testing.T) {
			got := sequence.Filter(test.seq, test.filterFn)
			require.Equal(t, test.want, got)
		})
	}

	for _, test := range testAnon {
		t.Run(test.name, func(t *testing.T) {
			got := sequence.Filter(test.seq, test.filterFn)
			require.Equal(t, test.want, got)
		})
	}
}

func TestMap(t *testing.T) {
	t.Run("iterates through the original and outputs the transformed list", func(t *testing.T) {
		// given
		original := []string{"abc", "doerayme"}
		expected := []string{"ABC", "DOERAYME"}

		// when
		output := sequence.Map(original, func(str string) string {
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
		output := sequence.Map(original, func(str string) int {
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
		output := sequence.Map(original, func(str string) int {
			callbackHasBeenCalled = true
			return len(str)
		})

		// then
		require.False(t, callbackHasBeenCalled)
		require.EqualValues(t, expected, output)
	})
}

func TestMapIndex(t *testing.T) {
	t.Run("iterates through the original and outputs the transformed list", func(t *testing.T) {
		// given
		original := []string{"abc", "doerayme"}
		expected := []string{"ABC0", "DOERAYME1"}

		// when
		output := sequence.MapIndex(original, func(i int, str string) string {
			return fmt.Sprintf("%s%d", strings.ToUpper(str), i)
		})

		// then
		require.EqualValues(t, expected, output)
	})

	t.Run("iterates through the original and outputs a transformed list of a different type", func(t *testing.T) {
		// given
		original := []string{"abc", "doerayme"}
		expected := []int{3, 9}

		// when
		output := sequence.MapIndex(original, func(i int, str string) int {
			return len(str) + i
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
		output := sequence.MapIndex(original, func(i int, str string) int {
			callbackHasBeenCalled = true
			return len(str)
		})

		// then
		require.False(t, callbackHasBeenCalled)
		require.EqualValues(t, expected, output)
	})
}

func TestReduce(t *testing.T) {
	t.Run("reduces through the slice and outputs the maximum length found", func(t *testing.T) {
		// given
		original := []string{"abc", "abcd", "abcdef", "ab"}
		expected := 6

		// when
		output := sequence.Reduce(original, 0, func(acc int, el string) int {
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
		output := sequence.Reduce(original, map[string]int{}, func(acc map[string]int, el string) map[string]int {
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
		output := sequence.Reduce(original, map[string]int{}, func(acc map[string]int, el string) map[string]int {
			callbackHasBeenCalled = true
			acc[el] = len(el)
			return acc
		})

		// then
		require.False(t, callbackHasBeenCalled)
		require.EqualValues(t, expected, output)
	})
}
