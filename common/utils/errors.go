package utils

import (
	"fmt"
	"os"
)

func Check(err error, msg string) {
	if err != nil {
		RecoverableCheck(err, msg)
		os.Exit(1)
	}
}

func RecoverableCheck(err error, msg string) {
	if err != nil {
		fmt.Println(msg + "; " + err.Error())
	}
}

func Assert(ok bool, msg string) {
	if !ok {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func CheckWithFunc(err error, fn func() string) {
	if err != nil {
		msg := fn()

		Check(err, msg)
	}
}
