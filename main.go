package main

import (
	helper "drvcheck/drvcheck"
	"fmt"
	// "fmt"
	// "github.com/goccy/go-yaml"
)

func main() {

	defer func() {
		panic_err := recover()
		if panic_err != nil {
			fmt.Println("Panic", panic_err)
		}
	}()
	// todo: executable?
	// todo: http://golang.org.pl/lang/08_oop.html
	helper.Run()

}
