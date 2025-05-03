package main

import (
	"fmt"
	"github.com/g1eng/select5"
)

func main() {
	l := [][]any{{
		"a", "Arista Networks", 1.5, false,
	}, {
		"b", "Broadcom", 3.82, false,
	}, {
		"c", "Cisco", 10.51, true,
	}, {
		"d", "Docker", 3.14, false,
	}, {
		"e", "Equinix", 0.09, true,
	}, {
		"f", "Fortinet", 15.1, true,
	}, {
		"g", "Godaddy", -9.2, true,
	}, {
		"h", "HuggingFace", 39.15, false,
	}}
	s, err := select5.SelectTableRow(l)
	if err != nil {
		panic(err)
	}
	w, err := select5.GetV(s[0])
	if err != nil {
		panic(err)
	}
	println(w)
	fmt.Printf(select5.MoveTo, 2, 0)
}
