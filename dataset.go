package select5

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
	"os"
	"strings"
)

const (
	IsString byte = 0x01
	IsInt    byte = 0x02
	IsInt8   byte = 0x02
	IsInt16  byte = 0x02

	IsInt32   byte = 0x02
	IsInt64   byte = 0x04
	IsFloat32 byte = 0x08
	IsFloat64 byte = 0x10
	IsBool    byte = 0x20
	IsPointer byte = 0x40
	IsAny     byte = 0x7f
	IsList    byte = 0x00
	IsTable   byte = 0x80
)

// Selector is the alias of Dataset (for backword compatibility)
type Selector Dataset

// Dataset represents a selectable dataset with an optional header
type Dataset struct {
	Header []string // selection header
	Data   any
}

// NewDatasetFrom creates a new Dataset from a slice of any type
func NewDatasetFrom(p []any) *Dataset {
	var a []any
	for _, s := range p {
		a = append(a, s)
	}
	return &Dataset{
		Header: []string{},
		Data:   a,
	}
}

// NewSelectorFrom is the alias of NewDatasetFrom
func NewSelectorFrom(p []any) *Dataset {
	return NewDatasetFrom(p)
}

// Type determines the type of data in the selector
// Returns a byte value representing the data type(s).
func (d *Dataset) Type() byte {

	var elementType = IsList
	switch d.Data.(type) {
	case []string:
		elementType |= IsString
	case []any:
		for _, r := range d.Data.([]any) {
			elementType |= CheckPrimitive(r)
		}
	case [][]any:
		for i, _ := range d.Data.([][]any) {
			elementType |= IsTable
			for j, _ := range d.Data.([][]any)[i] {
				elementType |= CheckPrimitive(d.Data.([][]any)[i][j])
				if elementType == IsAny|IsTable {
					break
				}
			}
		}
	default:
		return elementType //non-list
	}
	typeMask := IsAny & elementType
	if typeMask&(typeMask-1) != 0 { // two or more types detected
		elementType |= IsAny
	}
	return elementType
}

// IsList presents whether the dataset is a simple list, not a two-dimensional data
func (d *Dataset) IsList() bool {
	return d.Type()&IsTable == 0
}

// IsTable presents whether the dataset is a two-dimensional data or not
func (d *Dataset) IsTable() bool {
	return d.Type()&IsTable != 0
}

// TypeList presents the list of byte representations of the list.
func (d *Dataset) TypeList() ([]byte, error) {
	var res []byte
	if d.IsTable() {
		return nil, fmt.Errorf("dataset is not a list, but a table")
	}
	for _, a := range d.Data.([]any) {
		res = append(res, CheckPrimitive(a))
	}
	return res, nil
}

// TypeTable presents the table of byte representations of the data.
func (d *Dataset) TypeTable() ([][]byte, error) {
	var res [][]byte
	if d.IsList() {
		return nil, fmt.Errorf("dataset is not a list, but a table")
	}
	for _, ao := range d.Data.([][]any) {
		var b []byte
		for _, a := range ao {
			b = append(b, CheckPrimitive(a))
		}
		res = append(res, b)
	}
	return res, nil
}

// RenderMenu draws the menu with the current selection (internal use)
func RenderMenu(list []string, selectedIndex int, prevIndex int) {

	// Position cursor at the top
	fmt.Print(ResetCursor)
	if selectedIndex == prevIndex && selectedIndex == 0 {
		for i, item := range list {
			if i == 0 {
				fmt.Printf(MoveTo, i+1, 1)
				fmt.Print(ClearLine)
				fmt.Print("> ")
				fmt.Print(item)
			} else {
				fmt.Printf(MoveTo, i+1, 1)
				fmt.Print(ClearLine)
				fmt.Print("  ")
				fmt.Print(item)
			}
		}
		return
	}

	for i, item := range list {
		if i == prevIndex {
			fmt.Printf(MoveTo, i+1, 1)
			fmt.Print(ClearLine)
			fmt.Print("  ")
			fmt.Print(item)
		}
		if i == selectedIndex {
			fmt.Printf(MoveTo, i+1, 1)
			fmt.Print(ClearLine)
			fmt.Print("> ")
			fmt.Print(item)
		}
	}
}

// Select performs the selection based on the data type.
// Returns the selected item or an error if selection is not supported
func (d *Dataset) Select() (any, error) {
	if d.Type()&IsTable == IsTable {
		return SelectTableRow(d.Data.([][]any))
	} else if d.Type()&IsAny == IsString {
		return SelectString(d.Data.([]string))
	} else {
		return nil, fmt.Errorf("selection not supported for the type %d %T", d.Type(), d.Data)
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

	var oldState *term.State
	isTerm := term.IsTerminal(int(os.Stdin.Fd()))
	if isTerm {
		var err error
		oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)
		defer term.Restore(int(os.Stdout.Fd()), oldState)
	}

	fmt.Print(ClearScreen)
	fmt.Print(HideCursor)
	defer fmt.Print(ShowCursor)

	keyEvents, sigChan := CaptureKeyboardEvents()

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

			if key.Special != 0 {
				switch key.Special {
				case UP:
					selectedIndex = (selectedIndex - 1 + len(list)) % len(list)
					RenderMenu(list, selectedIndex, prevIndex)
				case DOWN:
					selectedIndex = (selectedIndex + 1) % len(list)
					RenderMenu(list, selectedIndex, prevIndex)
				case ENTER:
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
			fmt.Printf(MoveTo, len(list)+1, 1)
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
	var oldState *term.State
	isTerm := term.IsTerminal(int(os.Stdin.Fd()))
	if isTerm {
		var err error
		oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return nil, err
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)
		defer term.Restore(int(os.Stdout.Fd()), oldState)
	}

	fmt.Print(ClearScreen)
	fmt.Print(ResetCursor)
	fmt.Print(HideCursor)

	keyEvents, sigChan := CaptureKeyboardEvents()

	selectedIndex := 0

	// Initial render of the menu
	RenderTable(list, selectedIndex)

	for {
		select {
		case key, ok := <-keyEvents:
			if !ok {
				return nil, fmt.Errorf("keyboard event channel closed")
			}

			if key.Special != 0 {
				switch key.Special {
				case UP:
					selectedIndex = (selectedIndex - 1 + len(list)) % len(list)
					// Clear and reposition cursor before redrawing
					fmt.Print(ClearScreen)
					fmt.Print(ResetCursor)
					RenderTable(list, selectedIndex)
				case DOWN:
					selectedIndex = (selectedIndex + 1) % len(list)
					// Clear and reposition cursor before redrawing
					fmt.Print(ClearScreen)
					fmt.Print(ResetCursor)
					RenderTable(list, selectedIndex)
				case ENTER:
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
