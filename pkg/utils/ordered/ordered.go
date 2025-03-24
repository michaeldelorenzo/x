// Package ordered provides generic utils for manipulating and comparing ordered inputs
package ordered

import "golang.org/x/exp/constraints"

// Min find the minimum of two ordered types
func Min[I constraints.Ordered](x, y I) I {
	if x < y {
		return x
	}
	return y
}

// Max finds the maximum of two ordered types
func Max[I constraints.Ordered](x, y I) I {
	if x > y {
		return x
	}
	return y
}
