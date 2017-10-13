package tutils

import (
	"fmt"
	"os"
)

func Check(err error) bool {
	if err == nil {
		return true
	}
	return false
}

func Log(msg interface{}) {
	fmt.Printf("[Info] %v \n", msg)
}

func CheckWarn(err error) {
	if Check(err) != true {
		fmt.Printf("[Warning] %s \n", err)
	}
}

func CheckExit(err error) {
	if Check(err) != true {
		fmt.Printf("[Fatal] %s \n", err)
		os.Exit(1)
	}
}
