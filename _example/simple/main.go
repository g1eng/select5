package main

import (
	"fmt"
	"github.com/g1eng/select5"
)

func main() {
	s, err := select5.SelectString([]string{"a", "b", "V", "O", "p"})
	if err != nil {
		panic(err)
	}
	println(s)
	fmt.Printf(select5.MoveTo, 2, 0)
}
