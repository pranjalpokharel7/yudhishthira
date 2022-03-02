package utility

import (
	"log"
)

func ErrThenPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func ErrThenLogPanic(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ErrThenLogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
