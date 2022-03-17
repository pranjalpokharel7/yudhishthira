package utility

import (
	"errors"
	"log"
	"net"
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

func getForwardSlashPosition(value string) int {
	for i, c := range value {
		if c == '/' {
			return i
		}
	}

	return -1
}

func GetNodeAddress() string {

	addresses, err := net.InterfaceAddrs()
	if err != nil {
		ErrThenLogPanic(err)
	}

	for _, addr := range addresses {
		addr_string := addr.String()
		position := getForwardSlashPosition(addr_string)

		if addr_string[:3] == "192" {
			return addr_string[:position]
		}
	}

	err = errors.New("Address not found")
	log.Panic(err)
	return ""
}
