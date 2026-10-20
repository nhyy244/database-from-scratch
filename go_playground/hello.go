package main

import (
	"fmt"
	"db_project"
	"rsc.io/quote"
)

func main() {
	fmt.Println("Hello World")
	fmt.Println(quote.Go())
	fmt.Println(quote.Glass())
	var x string = "hello"
	fmt.Println(x)
	{
		x := 10
		fmt.Println(x)
	}
	var strBytes []byte = []byte("k1")
	fmt.Println(strBytes)
	fmt.Println(string(strBytes))
}
