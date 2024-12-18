package main

import "fmt"

type X int

func main() {
	var x X
	for i := 0; i < 10; i++ {
		x = i
		fmt.Println(x)
	}
}
