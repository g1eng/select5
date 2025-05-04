# SELECT-5

[![codecov](https://codecov.io/gh/g1eng/select5/branch/main/graph/badge.svg)](https://codecov.io/gh/g1eng/select5)

 `select5` provides interactive terminal-based selection utilities for Go applications.
 It allows users to select items from lists and tables using keyboard navigation.

The package aims to provide simple and slim interactive layer for CLI written in Go.

For more rich experience (but fat package), use [bubble](https://github.com/charmbracelet/bubbletea) or other solutions.


 # Overview

 select5 offers two primary selection modes:
 - Simple string selection from a list
 - Advanced table row selection with mixed data types

 Both modes support keyboard navigation with arrow keys and selection with Enter.
 The library handles terminal control sequences and cursor movement automatically.

 # Basic String Selection

 Use SelectString to display a simple list of strings for selection:

 	selected, err := select5.SelectString([]string{"Option A", "Option B", "Option C"})
 	if err != nil {
 	    log.Fatal(err)
 	}
 	fmt.Println("You selected:", selected)

 # Table Row Selection

 For more complex data, SelectTableRow supports tables with mixed data types (integers, floats, strings, booleans):

 	data := [][]any{
 	    {"a", "Apple Inc.", 178.72, true},
 	    {"b", "Broadcom", 376.04, false},
 	    {"c", "Cisco", 125.30, true},
 	}

 	selectedRow, err := select5.SelectTableRow(data)
 	if err != nil {
 	    log.Fatal(err)
 	}

 	// Access selected data
 	code := selectedRow[0].(string)
 	name := selectedRow[1].(string)

 	fmt.Printf("Selected: %s - %s\n", code, name)

 # Type Helpers

 The package includes helper functions to safely extract and convert values from the any type:

 	// Extract string value
 	code, err := select5.GetS(selectedRow[0])

 	// Extract float32 value
 	price, err := select5.GetF32(selectedRow[2])

 	// Extract bool value
 	active, err := select5.GetB(selectedRow[3])

 # Terminal Control

 The package provides constants for terminal control operations:

 	// Position cursor
 	fmt.Printf(select5.MoveTo, line, column)

 	// Clear screen
 	fmt.Print(select5.ClearScreen)

 	// Reset cursor to home position
 	fmt.Print(select5.ResetCursor)

 # Keyboard Navigation

 - Up/Down arrows: Move selection
 - Enter: Confirm selection
 - 'q' or Ctrl+C: Quit without selection

 # Error Handling

 All selection functions return appropriate errors that should be checked:
 - Empty lists
 - Type conversion failures
 - Keyboard event channel closure
 - Interrupted selection

# Author

Select5 contributors