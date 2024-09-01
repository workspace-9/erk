// erk is a kit for errors
package erk

// Try assumes that err is nil and returns t, otherwise it panics.
func Try[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

// Must assumes that err is nil, otherwise it panics.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
