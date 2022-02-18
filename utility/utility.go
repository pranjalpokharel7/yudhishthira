package utility

func ErrThenPanic(err error) {
	if err != nil {
		panic(err)
	}
}
