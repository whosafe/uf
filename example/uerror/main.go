package main

import (
	"fmt"

	"github.com/whosafe/uf/uerror"
)

func testError() error {
	return uerror.New("test error1")
}

func recursiveCall(n int) {
	if n >= 40 {
		err := uerror.New("deep error")
		fmt.Printf("%+v", err)
		return
	}
	recursiveCall(n + 1)
}

func main() {

	// err := uerror.New("test error")
	// if err != nil {
	// 	fmt.Printf("%+v", err)
	// }

	fmt.Printf("%+v", testError())
	// recursiveCall(0)

}
