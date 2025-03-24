// Package maputils contains many common utilities for operating on and iterating over maps
package maputils

// Keys returns a new slice populated with the keys in a map. Note that there is no guarentee on order.
func Keys[K comparable, V any](input map[K]V) []K {
	keys := make([]K, len(input))
	i := 0
	for k := range input {
		keys[i] = k
		i += 1
	}

	return keys
}

// Values returns a new slice populated with the values in a map. Note that there is no guarentee on order.
func Values[K comparable, V any](input map[K]V) []V {
	values := make([]V, len(input))
	i := 0
	for _, v := range input {
		values[i] = v
		i += 1
	}

	return values
}

// Map returns a new slice populated with the result of calling the provided function
// on every key-value pair in the provided input map. Note that there is no guarentee on order.
func Map[K comparable, V, O any](input map[K]V, fn func(key K, val V) O) []O {
	output := make([]O, len(input))
	i := 0
	for k, v := range input {
		output[i] = fn(k, v)
		i += 1
	}
	return output
}

// Reduce iterates through every key-value pair calling the provided fn to build up a new output.
// Useful when needing to reduce a map into another single value such as map, the min/max of a list, etc.
// Note that there is no guarentee on order.
func Reduce[K comparable, V, O any](input map[K]V, initialValue O, fn func(acc O, key K, val V) O) O {
	accumulator := initialValue
	for k, v := range input {
		accumulator = fn(accumulator, k, v)
	}
	return accumulator
}

// GetAndCoalesce pulls a value from a map of interfaces and coalesces it into expected type T. If the value is not found in the map
// or the type assertion fails, the zero value of type T is returned instead.
func GetAndCoalesce[T any, K comparable](input map[K]interface{}, key K) (T, bool) {
	valInterface, ok := input[key]
	if !ok {
		var zero T
		return zero, false
	}

	value, ok := valInterface.(T)
	if !ok {
		var zero T
		return zero, false
	}

	return value, true
}
