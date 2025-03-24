// Package sequence contains many common operations that depend on sequential access to a sequence’s values.
package sequence

// Contains determines whether a provided sequence type contains the provided
// comparable element.
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// Unique takes an input slice and returns a new slice without duplicate values.
//
// Note that order is preserved. If multiple elements compare equal, then the index of the first equal element will be retained.
func Unique[T comparable](s []T) []T {
	indexByElement := make(map[T]struct{}, len(s))
	uniqueElements := make([]T, 0, len(s))

	for _, el := range s {
		if _, ok := indexByElement[el]; !ok {
			indexByElement[el] = struct{}{}
			uniqueElements = append(uniqueElements, el)
		}
	}

	return uniqueElements
}

// Flatten takes an input slice of slices and flattens them into a slice of elements
func Flatten[T any](s [][]T) []T {
	// iterate through the slice first to find the final the total number of elements to allocate
	totalNumberOfElements := Reduce(s, 0, func(acc int, el []T) int {
		return acc + len(el)
	})

	// allocate the output array
	output := make([]T, totalNumberOfElements)

	// nested iteration through slice of slices to fill in the flattened output
	i := 0
	for _, arr := range s {
		for _, el := range arr {
			output[i] = el
			i += 1
		}
	}

	return output
}

// Filter takes an input slice and calls filterFn on each element, returning a new slice containing only elements
// where filterFn returned `true`
func Filter[T any](s []T, filterFn func(el T) bool) []T {
	// allocate the output array, worst case will be equal to `s` in size, but best case every element will be filtered out
	output := make([]T, 0, len(s))

	for _, el := range s {
		shouldKeepInList := filterFn(el)
		if shouldKeepInList {
			output = append(output, el)
		}
	}

	return output
}

// Map returns a new slice populated with the result of calling the provided function
// on every element in the provided input slice.
func Map[I, O any](input []I, fn func(I) O) []O {
	output := make([]O, len(input))
	for i, v := range input {
		output[i] = fn(v)
	}
	return output
}

// MapIndex returns a new slice populated with the result of calling the provided function
// on every index and corresponding element in the provided input slice.
func MapIndex[I, O any](input []I, fn func(int, I) O) []O {
	output := make([]O, len(input))
	for i, v := range input {
		output[i] = fn(i, v)
	}
	return output
}

// Reduce iterates through every element calling the provided fn to build up a new output.
// Useful when needing to reduce a slice into a single value such as map, the min/max of a list, etc.
func Reduce[I, O any](input []I, initialValue O, fn func(acc O, elem I) O) O {
	accumulator := initialValue
	for _, v := range input {
		accumulator = fn(accumulator, v)
	}
	return accumulator
}
