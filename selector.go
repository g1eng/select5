package select5

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// RenderMenu draws the menu with the current selection (internal use)
func RenderMenu(list []string, selectedIndex int, prevIndex int) {

	// Position cursor at the top
	fmt.Print(ResetCursor)
	if selectedIndex == prevIndex && selectedIndex == 0 {
		for i, item := range list {
			if i == 0 {
				fmt.Printf("\033[%d;1H", i+1)
				fmt.Print(ClearLine)
				fmt.Print("> ")
				fmt.Print(item)
			} else {
				fmt.Printf("\033[%d;1H", i+1)
				fmt.Print(ClearLine)
				fmt.Print("  ")
				fmt.Print(item)
			}
		}
		return
	}

	for i, item := range list {
		if i == prevIndex {
			fmt.Printf("\033[%d;1H", i+1)
			fmt.Print(ClearLine)
			fmt.Print("  ")
			fmt.Print(item)
		}
		if i == selectedIndex {
			fmt.Printf("\033[%d;1H", i+1)
			fmt.Print(ClearLine)
			fmt.Print("> ")
			fmt.Print(item)
		}
	}
}

// SelectString presents a list of strings for selection and returns the selected string.
// It displays an interactive cursor that can be moved with arrow keys.
// Returns the selected string or an error if:
// - the provided slice is empty
// - the keyboard event channel closes
// - the user quits (q or Ctrl+C)
func SelectString(list []string) (string, error) {
	if len(list) == 0 {
		return "", fmt.Errorf("zero length list provided")
	}
	fmt.Print(ClearScreen)
	fmt.Print(HideCursor)

	keyEvents := CaptureKeyboardEvents(os.Stdin)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	selectedIndex := 0
	prevIndex := 0

	// Initial render of the menu
	RenderMenu(list, selectedIndex, prevIndex)

	for {
		prevIndex = selectedIndex
		select {
		case key, ok := <-keyEvents:
			if !ok {
				return "", fmt.Errorf("keyboard event channel closed")
			}

			if key.Special != "" {
				switch key.Special {
				case "UP":
					selectedIndex = (selectedIndex - 1 + len(list)) % len(list)
					RenderMenu(list, selectedIndex, prevIndex)
				case "DOWN":
					selectedIndex = (selectedIndex + 1) % len(list)
					RenderMenu(list, selectedIndex, prevIndex)
				case "ENTER":
					// Clear screen and show the selection
					fmt.Print(ClearScreen)
					fmt.Print(ResetCursor)
					fmt.Print(ShowCursor)
					return list[selectedIndex], nil
				}
			} else if key.Key == 'q' || (key.Ctrl && key.Key == 'c') {
				// Quit on 'q' or Ctrl+C
				fmt.Print(ClearScreen)
				fmt.Print(ResetCursor)
				fmt.Print('\n')
				return "", nil
			}

		case <-sigChan:
			fmt.Printf("\033[%d;1H", len(list)+1)
			return "", nil
		}
	}
}

// RenderTable draws the table with a row cursor. (internal use)
func RenderTable(list [][]any, selectedIndex int) error {
	var buf bytes.Buffer
	t := tablewriter.NewWriter(&buf)
	if selectedIndex < 0 {
		selectedIndex = 0
	}

	for _, row := range list {
		var newRow []string
		for _, r := range row {
			v, err := GetV(r)
			if err != nil {
				return err
			}
			newRow = append(newRow, v)
		}
		t.Append(newRow)
	}
	t.SetBorder(false)
	t.Render()

	data := buf.Bytes()
	if len(data) == 0 {
		return fmt.Errorf("no table data")
	}

	tableRowStringSlices := strings.Split(string(data), "\n")
	for i, row := range tableRowStringSlices {
		fmt.Printf(MoveTo, i+1, 1)
		if i == selectedIndex {
			fmt.Printf("\x1b[01;07m%s\x1b[01;00m\n", row)
		} else {
			fmt.Println(row)
		}
	}
	return nil
}

// SelectTableRow presents a table of mixed data types for selection and returns the selected row.
// Each row can contain different data types (string, int, float, bool, etc.).
// Returns the selected row as []any or an error if:
// - the provided slice is empty
// - the keyboard event channel closes
// - the user quits (q or Ctrl+C)
func SelectTableRow(list [][]any) ([]any, error) {
	if len(list) == 0 {
		return nil, fmt.Errorf("zero length list provided")
	}

	fmt.Print(ClearScreen)
	fmt.Print(ResetCursor)
	fmt.Print(HideCursor)

	keyEvents := CaptureKeyboardEvents(os.Stdin)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	selectedIndex := 0

	// Initial render of the menu
	RenderTable(list, selectedIndex)

	for {
		select {
		case key, ok := <-keyEvents:
			if !ok {
				return nil, fmt.Errorf("keyboard event channel closed")
			}

			if key.Special != "" {
				switch key.Special {
				case "UP":
					selectedIndex = (selectedIndex - 1 + len(list)) % len(list)
					// Clear and reposition cursor before redrawing
					fmt.Print(ClearScreen)
					fmt.Print(ResetCursor)
					RenderTable(list, selectedIndex)
				case "DOWN":
					selectedIndex = (selectedIndex + 1) % len(list)
					// Clear and reposition cursor before redrawing
					fmt.Print(ClearScreen)
					fmt.Print(ResetCursor)
					RenderTable(list, selectedIndex)
				case "ENTER":
					// Clear screen and show the selection
					fmt.Print(ClearScreen)
					fmt.Print(ResetCursor)
					fmt.Print(ShowCursor)
					return list[selectedIndex], nil
				}
			} else if key.Key == 'q' || (key.Ctrl && key.Key == 'c') {
				// Quit on q or Ctrl+C
				fmt.Print(ClearScreen)
				fmt.Print(ResetCursor)
				fmt.Print(ShowCursor)
				return nil, nil
			}

		case <-sigChan:
			fmt.Print(ClearScreen)
			fmt.Print(ResetCursor)
			fmt.Print(ShowCursor)
			return nil, nil
		}
	}
}
