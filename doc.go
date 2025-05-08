package select5

// Package select5 provides interactive terminal-based selection utilities for Go applications.
// It allows you to select items from lists and tables using keyboard navigation.
//
// The package aims to provide simple and slim interactive layer for CLI written in Go.
//
// For more rich experience (but fat package), use [bubble](https://github.com/charmbracelet/bubbletea) or other solutions.
//
// # Overview
//
// select5 offers two primary selection modes:
// - Simple string selection from a list
// - Advanced table row selection with mixed data types
// - Console text editor with Emacs-like key binding
//
// Both modes support keyboard navigation with arrow keys and selection with Enter.
// The library handles terminal control sequences and cursor movement automatically.
//
// The package also implements an experimental console text editor with Emacs-like key binding.
// For more detail, see the [example](#text-editor).
//
// # Basic Utf8Char Selection
//
// Use SelectString to display a simple list of strings for selection:
//
// 	selected, err := select5.SelectString([]string{"Option A", "Option B", "Option C"})
// 	if err != nil {
// 	    log.Fatal(err)
// 	}
// 	fmt.Println("You selected:", selected)
//
// # Table Row Selection
//
// For more complex data, SelectTableRow supports tables with mixed data types (integers, floats, strings, booleans):
//
// 	data := [][]any{
// 	    {"a", "Apple Inc.", 178.72, true},
// 	    {"b", "Broadcom", 376.04, false},
// 	    {"c", "Cisco", 125.30, true},
// 	}
//
// 	selectedRow, err := select5.SelectTableRow(data)
// 	if err != nil {
// 	    log.Fatal(err)
// 	}
//
// 	// Access selected data
// 	code := selectedRow[0].(string)
// 	name := selectedRow[1].(string)
//
// 	fmt.Printf("Selected: %s - %s\n", code, name)
//
// # Generic entrypoint for data selector
//
// For more flexible implementation, you can use `Selector` struct to declare selectable data which may be list or table of primitives, or any type.
//
//	strPointed := "WORD"
//	otherStrPointed := "VERB"
//	list := select5.Selector{
//		Header: nil,
//		Data: [][]any{
//			{"AGI", 3, true, 3.58, nil},
//			{"BGM", 2, true, 0.9, nil},
//			{"CGC", 1, false, -93.20, &strPointed},
//			{"DPO", 1829, true, 3.58, &otherStrPointed},
//		},
//	}
// res, _ := list.Select()
// fmt.Printf("%v", res)
//
// You can apply additional type casting for the interface (a.k.a. `any` type) results.
//
// # Type Helpers for primitives
//
// The package includes helper functions to safely extract and convert values from the `any` type:
//
// 	// Extract string value
// 	code, err := select5.GetS(selectedRow[0])
//
// 	// Extract float value
// 	price, err := select5.GetF(selectedRow[2])
//
// 	// Extract bool value
// 	active, err := select5.GetB(selectedRow[3])
//
// # Builtin Type Detector for `Selector`
//
// The `Selector` implements type detector in `Type()`, which returns type information in byte expression.
//
// t := select5.Selector{
//    Data:   []any{"a","b","c"},
// }.Type()
// // 0x01 == IsList | IsString
//
// The type information is calculated by using following bitmask constants:
//
// 	IsString byte = 0x01
//	IsInt    byte = 0x02
//	IsInt8   byte = 0x02
//	IsInt16  byte = 0x02
//
//	IsInt32   byte = 0x02
//	IsInt64   byte = 0x04
//	IsFloat32 byte = 0x08
//	IsFloat64 byte = 0x10
//	IsBool    byte = 0x20
//	IsPointer byte = 0x40
//	IsAny     byte = 0x7f
//	IsList    byte = 0x00
//	IsTable   byte = 0x80
//
// For example, a float64 list will return `0x10` and a table of string and int values will return `0xff`.
// (If one or more types have been detected in the data, the result will be a logical sum by IsAny 0x7f).
// You can implement type switcher for various data structure, using this mechanism.
//
// # Terminal Control
//
// The package provides constants for terminal control operations:
//
// 	// Position cursor
// 	fmt.Printf(select5.MoveTo, line, column)
//
// 	// Clear screen
// 	fmt.Print(select5.ClearScreen)
//
// 	// Reset cursor to home position
// 	fmt.Print(select5.ResetCursor)
//
// # Keyboard Navigation
//
// - Up/Down arrows: Move selection
// - Enter: Confirm selection
// - 'q' or Ctrl+C: Quit without selection
//
// # Error Handling
//
// All selection functions return appropriate errors that should be checked:
// - Empty lists
// - Type conversion failures
// - Keyboard event channel closure
// - Interrupted selection
//
// # Text Editor
//
// select5 is not just a selection helper package, but also provides a tiny text editor.
// You can easily embed the editor in your application:
//
// import (
//	"fmt"
//	"github.com/g1eng/select5"
//	"os"
//	"strings"
//)
//
//func main() {
//	ed := select5.NewEditor()
//	res := ed.Edit()
//	t, err := os.CreateTemp(".", "result-*.txt")
//	defer t.Close()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf(select5.ClearScreen)
//	fmt.Printf(select5.MoveTo, 1, 0)
//	fmt.Print("[RESULT]")
//	for i, s := range strings.Split(res, "\n") {
//		if i != 0 {
//			fmt.Fprint(t, "\n")
//		}
//		fmt.Printf(select5.MoveTo, i+2, 0)
//		fmt.Print(s)
//		fmt.Fprint(t, s)
//	}
//	println()
//	fmt.Printf(select5.ClearScreenFromCursor)
//}
