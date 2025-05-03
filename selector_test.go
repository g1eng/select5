package select5_test

import (
	"github.com/g1eng/select5"
	"os"
	"testing"
	"time"
)

func TestRenderMenu(t *testing.T) {
	select5.RenderMenu([]string{"a", "b", "c"}, 1, 0)
}

func TestRenderTable(t *testing.T) {
	select5.RenderTable([][]any{{"a", 1}, {"b", "kichi"}, {"c", 1000.0}}, 3)
	select5.RenderTable([][]any{{"a", 1}, {"b", "kichi"}, {"c", 1000.0, true}}, 0)
}

func TestSelectString(t *testing.T) {
	// Test data
	options := []string{"Option 1", "Option 2", "Option 3"}

	// Create a pipe to simulate keyboard input
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	// Save original stdin and replace with our pipe
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }() // Restore stdin when test completes

	// Set up a goroutine to run SelectString
	resultCh := make(chan string)
	go func() {
		// Passing our output buffer to capture what would normally go to stdout
		result, err := select5.SelectString(options)
		if err != nil {
			panic(err)
		}
		resultCh <- result
	}()

	// Wait a moment for the prompt to be displayed
	time.Sleep(100 * time.Millisecond)

	// Simulate pressing DOWN arrow and then ENTER to select the second option
	// Send key sequence: DOWN arrow followed by ENTER
	w.Write([]byte{0x1b, '[', 'B'}) // DOWN arrow
	time.Sleep(50 * time.Millisecond)
	w.Write([]byte{13}) // ENTER key

	// Wait for result with timeout
	select {
	case result := <-resultCh:
		if result != "Option 2" {
			t.Fatalf("Expected 'Option 2' to be selected, got '%s'", result)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Test timed out waiting for selection")
	}
}

func TestSelectStringWithNoOptions(t *testing.T) {
	_, err := select5.SelectString([]string{})
	if err == nil {
		t.Fatal("Expected error, got none")
	}
}

func TestSelectTableRow(t *testing.T) {
	rows := [][]any{
		{1, "Alice", "alice@example.com", false},
		{2, "Bob", "bob@example.com", true},
		{3, "Charlie", "charlie@example.com", "FIRED!"},
	}

	// Create a pipe to simulate keyboard input
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	// Save original stdin and replace with our pipe
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }() // Restore stdin when test completes

	// Set up a goroutine to run SelectTableRow
	resultCh := make(chan []any)
	go func() {
		// Passing our output buffer to capture what would normally go to stdout
		result, err := select5.SelectTableRow(rows)
		if err != nil {
			panic(err)
		}
		resultCh <- result
	}()

	time.Sleep(100 * time.Millisecond)

	// Simulate pressing DOWN arrow twice and then ENTER to select the third row
	w.Write([]byte{0x1b, '[', 'B'}) // DOWN arrow
	time.Sleep(50 * time.Millisecond)
	w.Write([]byte{0x1b, '[', 'B'}) // DOWN arrow again
	time.Sleep(50 * time.Millisecond)
	w.Write([]byte{13}) // ENTER key

	// Wait for result with timeout
	select {
	case result := <-resultCh:
		// Check if we got the correct row (Charlie's data)
		expectedRow := rows[2]
		if !rowsEqual(result, expectedRow) {
			t.Fatalf("Expected row %v to be selected, got %v", expectedRow, result)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Test timed out waiting for selection")
	}
}

func TestSelectTableRowWithEmptyTable(t *testing.T) {
	_, err := select5.SelectTableRow([][]any{})
	if err == nil {
		t.Fatal("Expected error, got none")
	}
}

//func contains(s, substr string) bool {
//	return bytes.Contains([]byte(s), []byte(substr))
//}

func rowsEqual(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
