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
	w.Write([]byte{0x0a}) // ENTER key

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

func TestSelectStringWithBlankList(t *testing.T) {
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
	w.Write([]byte{0x0a}) // ENTER key

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

func TestSelector_Select_String(t *testing.T) {
	s := []string{"AGI", "BMI", "CES", "DEI", "ERR"}
	list := select5.Dataset{
		Header: nil,
		Data:   s,
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
	resultCh := make(chan any)
	errCh := make(chan error)

	go func() {
		res, err := list.Select()
		if err != nil {
			errCh <- err
		}
		resultCh <- res
	}()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	w.Write([]byte{0x1b, '[', 'B'}) // DOWN arrow
	time.Sleep(50 * time.Millisecond)
	w.Write([]byte{0x1b, '[', 'B'}) // DOWN arrow again
	time.Sleep(50 * time.Millisecond)
	w.Write([]byte{0x0a}) // ENTER key
	select {
	case result := <-resultCh:
		// Check if we got the correct row (Charlie's data)
		switch result.(type) {
		case string:
			if result.(string) != s[2] {
				t.Fatalf("Expected '%v', got '%v'", s[2], result.(string))
			}
		default:
			t.Fatalf("Expected result to be a string, got %T", result)
		}
	case e := <-errCh:
		t.Fatalf("Expected error channel to be nil, got %v", e)
	case <-time.After(3 * time.Second):
		t.Fatal("Test timed out waiting for selection")
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

//func TestSelector_Select_MixedList(t *testing.T) {
//
//	list := select5.Dataset{
//		Header: nil,
//		Dataset:   []any{"AGI", 2, true, 3.58, nil},
//	}
//	// Create a pipe to simulate keyboard input
//	r, w, err := os.Pipe()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer r.Close()
//	defer w.Close()
//
//	// Save original stdin and replace with our pipe
//	oldStdin := os.Stdin
//	os.Stdin = r
//	defer func() { os.Stdin = oldStdin }() // Restore stdin when test completes
//	resultCh := make(chan any)
//	errCh := make(chan error)
//
//	go func() {
//		res, err := list.Select()
//		if err != nil {
//			errCh <- err
//		}
//		resultCh <- res
//	}()
//	if err != nil {
//		t.Fatal(err)
//	}
//	time.Sleep(100 * time.Millisecond)
//
//	w.Write([]byte{0x1b, '[', 'B'}) // DOWN arrow
//	time.Sleep(50 * time.Millisecond)
//	w.Write([]byte{0x1b, '[', 'B'}) // DOWN arrow again
//	time.Sleep(50 * time.Millisecond)
//	w.Write([]byte{0x0a}) // ENTER key
//	select {
//	case result := <-resultCh:
//		// Check if we got the correct row (Charlie's data)
//		switch result.(type) {
//		case string:
//			if result.(bool) != list.Dataset.([]any)[2].(bool) {
//				t.Fatalf("Expected '%v', got '%v'", list.Dataset.([]any)[2], result)
//			}
//		default:
//			t.Fatalf("Expected result to be a string, got %T", result)
//		}
//	case e := <-errCh:
//		t.Fatalf("Expected error channel to be nil, got %v", e)
//	case <-time.After(3 * time.Second):
//		t.Fatal("Test timed out waiting for selection")
//	}
//}

func TestSelector_Select_MixedTable(t *testing.T) {

	strPointed := "WORD"
	otherStrPointed := "VERB"
	list := select5.Dataset{
		Header: nil,
		Data: [][]any{
			{"AGI", 3, true, 3.58, nil},
			{"BGM", 2, true, 0.9, nil},
			{"CGC", 1, false, -93.20, &strPointed},
			{"DPO", 1829, true, 3.58, &otherStrPointed},
		},
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
	resultCh := make(chan any)
	errCh := make(chan error)

	go func() {
		res, err := list.Select()
		if err != nil {
			errCh <- err
		}
		resultCh <- res
	}()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	w.Write([]byte{0x1b, '[', 'A'}) // DOWN arrow
	time.Sleep(50 * time.Millisecond)
	w.Write([]byte{0x1b, '[', 'A'}) // DOWN arrow again
	time.Sleep(50 * time.Millisecond)
	w.Write([]byte{0x0a}) // ENTER key
	select {
	case result := <-resultCh:
		// Check if we got the correct row (Charlie's data)
		switch result.(type) {
		case []any:
			switch result.([]any)[0].(type) {
			case string:
				for i, _ := range result.([]any) {
					switch result.([]any)[i].(type) {
					case string:
						if result.([]any)[i].(string) != list.Data.([][]any)[2][i].(string) {
							t.Fatalf("Expected '%v', got '%v'", list.Data.([][]any)[2][i], result)
						}
					case int:
						if result.([]any)[i].(int) != list.Data.([][]any)[2][i].(int) {
							t.Fatalf("Expected '%v', got '%v'", list.Data.([][]any)[2][i], result)
						}
					case float32:
						if result.([]any)[i].(float32) != list.Data.([][]any)[2][i].(float32) {
							t.Fatalf("Expected '%v', got '%v'", list.Data.([][]any)[2][i], result)
						}
					case bool:
						if result.([]any)[i].(bool) != list.Data.([][]any)[2][i].(bool) {
							t.Fatalf("Expected '%v', got '%v'", list.Data.([][]any)[2][i], result)
						}
					case *string:
						if result.([]any)[i].(*string) != list.Data.([][]any)[2][i].(*string) {
							t.Fatalf("Expected '%v', got '%v'", list.Data.([][]any)[2][i], result)
						}
					}
				}
			default:
				t.Fatalf("Expected result to be a string, got %T", result)
			}
		default:
			t.Fatalf("Expected result to be a string, got %T", result)
		}
	case e := <-errCh:
		t.Fatalf("Expected error channel to be nil, got %v", e)
	case <-time.After(3 * time.Second):
		t.Fatal("Test timed out waiting for selection")
	}
}

func TestSelector_Type(t *testing.T) {
	type fields struct {
		Header []string
		Data   any
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		{
			name: "simple string list",
			fields: fields{
				Data: []any{
					"a", "ho", "ko", "ra", "s", "i", "c", "k",
				},
			},
			want: select5.IsList | select5.IsString,
		},
		{
			name: "float64 list",
			fields: fields{
				Data: []any{
					1.414, 3.142, 2.718,
				},
			},
			want: select5.IsList | select5.IsFloat64,
		},
		{
			name: "any number list",
			fields: fields{
				Data: []any{
					1, 5, 1, 8, 1.573, 3.1412, 2.718,
				},
			},
			want: select5.IsList | select5.IsAny,
		},
		{
			name: "mixed list",
			fields: fields{
				Data: []any{
					1, 5, "Acknowledged", 8, 1.573, 3.1412, nil,
				},
			},
			want: select5.IsList | select5.IsAny,
		},
		{
			name: "simple string table",
			fields: fields{
				Data: [][]any{
					{"a", "ho", "ko", "ra"},
					{"s", "i", "c", "k"},
				},
			},
			want: select5.IsTable | select5.IsString,
		},
		{
			name: "int table",
			fields: fields{
				Data: [][]any{
					{123, 345, 6789},
					{101, 1, 110},
				},
			},
			want: select5.IsTable | select5.IsInt,
		},
		{
			name: "int64 table",
			fields: fields{
				Data: [][]any{
					{int64(123), int64(345), int64(678)},
					{int64(101), int64(1), int64(123)},
				},
			},
			want: select5.IsTable | select5.IsInt64,
		},
		{
			name: "float64 table",
			fields: fields{
				Data: [][]any{
					{1.23, 3.45, 6.78},
					{1.01, 0.01, 1.10},
				},
			},
			want: select5.IsTable | select5.IsFloat64,
		},
		{
			name: "mixed number table",
			fields: fields{
				Data: [][]any{
					{1.23, 3.45, -6.7890000190},
					{1, 0.00001, 0.00000},
				},
			},
			want: select5.IsTable | select5.IsAny,
		},
		{
			name: "mixed primitives table",
			fields: fields{
				Data: [][]any{
					{1.23, "直木", 35},
					{true, nil, 0.00000},
				},
			},
			want: select5.IsTable | select5.IsAny,
		},
		{
			name: "struct contained mixed table",
			fields: fields{
				Data: [][]any{
					{struct{}{}, "しゅわ", 3.5},
					{true, nil, 0.00000},
				},
			},
			want: select5.IsTable | select5.IsAny,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &select5.Dataset{
				Header: tt.fields.Header,
				Data:   tt.fields.Data,
			}
			if got := s.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelector_IsList_IsTable(t *testing.T) {
	type fields struct {
		Data any
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "simple string list",
			fields: fields{
				Data: []any{
					"a", "ho", "ko", "ra", "s", "i", "c", "k",
				},
			},
			want: true,
		},
		{
			name: "mixed list",
			fields: fields{
				Data: []any{
					1, 5, "Acknowledged", 8, 1.573, 3.1412, nil,
				},
			},
			want: true,
		},
		{
			name: "float64 table",
			fields: fields{
				Data: [][]any{
					{1.23, 3.45, 6.78},
					{1.01, 0.01, 1.10},
				},
			},
			want: false,
		},
		{
			name: "mixed primitives table",
			fields: fields{
				Data: [][]any{
					{1.23, "直木", 35},
					{true, nil, 0.00000},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &select5.Dataset{
				Data: tt.fields.Data,
			}
			if got := s.IsList(); got != tt.want {
				t.Errorf("IsList() = %v, want %v", got, tt.want)
			}
			if got := !s.IsTable(); got != tt.want {
				t.Errorf("IsTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelector_TypeList(t *testing.T) {
	type fields struct {
		Data any
	}
	tests := []struct {
		name      string
		fields    fields
		want      []byte
		wantError bool
	}{
		{
			name: "simple string list",
			fields: fields{
				Data: []any{
					"a", "ho", "ko", "ra", "s", "i", "c", "k",
				},
			},
			want: []byte{select5.IsString, select5.IsString, select5.IsString, select5.IsString, select5.IsString, select5.IsString, select5.IsString, select5.IsString},
		},
		{
			name: "mixed list",
			fields: fields{
				Data: []any{
					1, 5, "Acknowledged", 8, 1.573, 3.1412, nil,
				},
			},
			want: []byte{select5.IsInt, select5.IsInt, select5.IsString, select5.IsInt, select5.IsFloat64, select5.IsFloat64, select5.IsAny},
		},
		{
			name: "float64 table",
			fields: fields{
				Data: [][]any{
					{1.23, 3.45, 6.78},
					{1.01, 0.01, 1.10},
				},
			},
			wantError: true,
		},
		{
			name: "mixed primitives table",
			fields: fields{
				Data: [][]any{
					{1.23, "直木", 35},
					{true, nil, 0.00000},
				},
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &select5.Dataset{
				Data: tt.fields.Data,
			}
			res, err := s.TypeList()
			if tt.wantError && err == nil {
				t.Errorf("TypeList() error = %v, wantErr %v", err, tt.wantError)
			}
			for i, got := range res {
				if got != tt.want[i] {
					t.Errorf("IsList() = %v, want %v", got, tt.want[i])
				}
			}
		})
	}
}

func TestSelector_TypeTable(t *testing.T) {
	type fields struct {
		Data any
	}
	tests := []struct {
		name      string
		fields    fields
		want      [][]byte
		wantError bool
	}{
		{
			name: "simple string list",
			fields: fields{
				Data: []any{
					"a", "ho", "ko", "ra", "s", "i", "c", "k",
				},
			},
			wantError: true,
		},
		{
			name: "mixed list",
			fields: fields{
				Data: []any{
					1, 5, "Acknowledged", 8, 1.573, 3.1412, nil,
				},
			},
			wantError: true,
		},
		{
			name: "float64 table",
			fields: fields{
				Data: [][]any{
					{1.23, 3.45, 6.78},
					{1.01, 0.01, 1.10},
				},
			},
			want: [][]byte{
				{select5.IsFloat64, select5.IsFloat64, select5.IsFloat64},
				{select5.IsFloat64, select5.IsFloat64, select5.IsFloat64},
			},
		},
		{
			name: "mixed primitives table",
			fields: fields{
				Data: [][]any{
					{1.23, "直木", 35},
					{true, nil, 0.00000},
				},
			},
			want: [][]byte{
				{select5.IsFloat64, select5.IsString, select5.IsInt},
				{select5.IsBool, select5.IsAny, select5.IsFloat64},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &select5.Dataset{
				Data: tt.fields.Data,
			}
			res, err := s.TypeTable()
			if tt.wantError && err == nil {
				t.Errorf("TypeList() error = %v, wantErr %v", err, tt.wantError)
			}
			for i, gotOut := range res {
				for j, got := range gotOut {
					if got != tt.want[i][j] {
						t.Errorf("IsList() = %v, want %v", got, tt.want[i][j])
					}
				}
			}
		})
	}
}
