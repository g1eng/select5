package main

import (
	"fmt"
	"github.com/g1eng/select5"
	"os"
	"strings"
)

func main() {
	ed := select5.NewEditor()
	res := ed.Edit()
	t, err := os.CreateTemp(".", "result-*.txt")
	defer t.Close()
	if err != nil {
		panic(err)
	}
	fmt.Printf(select5.ClearScreen)
	fmt.Printf(select5.MoveTo, 1, 0)
	fmt.Print("[RESULT]")
	for i, s := range strings.Split(res, "\n") {
		if i != 0 {
			fmt.Fprint(t, "\n")
		}
		fmt.Printf(select5.MoveTo, i+2, 0)
		fmt.Print(s)
		fmt.Fprint(t, s)
	}
	println()
	fmt.Printf(select5.ClearScreenFromCursor)
}
