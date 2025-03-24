// Package iter contains helper functions for iterating over lists.
//
// Deprecated: The functionality of the github.com/koneksahealth/x/pkg/utils/iter package has been moved to
// github.com/koneksahealth/x/pkg/utils/sequence. Use github.com/koneksahealth/x/pkg/utils/sequence instead.
package iter

// Map returns a new slice populated with the result of calling the provided function
// on every element in the provided input slice.
//
// Deprecated: The functionality of the github.com/koneksahealth/x/pkg/utils/iter package has been moved to
// github.com/koneksahealth/x/pkg/utils/sequence. Use github.com/koneksahealth/x/pkg/utils/sequence#Map instead.
func Map[I, O any](input []I, fn func(I) O) []O {
	output := make([]O, len(input))
	for i, v := range input {
		output[i] = fn(v)
	}
	return output
}

// Reduce iterates through every element calling the provided fn to build up a new output.
// Useful when needing to reduce a slice into a single value such as map, the min/max of a list, etc.
//
// Deprecated: The functionality of the github.com/koneksahealth/x/pkg/utils/iter package has been moved to
// github.com/koneksahealth/x/pkg/utils/sequence. Use github.com/koneksahealth/x/pkg/utils/sequence#Reduce instead.
func Reduce[I, O any](input []I, initialValue O, fn func(acc O, elem I) O) O {
	accumulator := initialValue
	for _, v := range input {
		accumulator = fn(accumulator, v)
	}
	return accumulator
}
