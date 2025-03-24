// Package ptr provides generic utils for safely creating and dereferencing pointers
package ptr

// Addr returns a pointer to the provided value
//
// `myPtr := ptr.Addr("something you cant usually take the address of inline")`
func Addr[I any](v I) *I {
	return &v
}

// Deref coalesces a pointer, dereferencing or returning the zero value of the underlying type if the pointer is nil
func Deref[I any](p *I) I {
	if p == nil {
		var zero I
		return zero
	}

	return *p
}

// DerefWithFallback coalesces a pointer, dereferencing or returning the provided fallback value if the pointer is nil
func DerefWithFallback[I any](p *I, fallback I) I {
	if p == nil {
		return fallback
	}

	return *p
}
