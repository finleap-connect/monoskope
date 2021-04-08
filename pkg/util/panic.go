package util

// PanicOnError panics if the given error is not nil
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
