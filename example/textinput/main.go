package main

import (
	"fmt"
	"os"
)

func main() {
	ed := NewEditor()
	res := ed.Edit()
	println("[RESULT]\n\n", res)
	t, err := os.CreateTemp(".", "result-*.txt")
	defer t.Close()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(t, "%s", res)
}
