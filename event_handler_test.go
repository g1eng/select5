package select5_test

import (
	"github.com/g1eng/select5"
	"os"
	"testing"
	"time"
)

func TestCaptureKeyboardEventsA(t *testing.T) {
	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	keyChannel := select5.CaptureKeyboardEvents(r)

	key := 'a'
	_, err = w.Write([]byte{byte(key)})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case k := <-keyChannel:
		if k.Key != key {
			t.Fatalf("invalid key code: %c, expected %c", k.Key, key)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for key event")
	}

}

func TestCaptureKeyboardEventsEnter(t *testing.T) {
	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	keyChannel := select5.CaptureKeyboardEvents(r)

	key := '\n'
	_, err = w.Write([]byte{byte(key)})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case k := <-keyChannel:
		if k.Special != "ENTER" {
			t.Fatalf("invalid key code: %c, expected DEL", k.Key)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for key event")
	}
}

func TestCaptureKeyboardEventsUpArrow(t *testing.T) {
	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	keyChannel := select5.CaptureKeyboardEvents(r)

	// Send the UP arrow key escape sequence: ESC [ A
	_, err = w.Write([]byte{0x1b, '[', 'A'})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case k := <-keyChannel:
		if k.Special != "UP" {
			t.Fatalf("invalid special key: %v, expected UP", k.Special)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for key event")
	}
}

func TestCaptureKeyboardEventsDownArrow(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	keyChannel := select5.CaptureKeyboardEvents(r)

	_, err = w.Write([]byte{0x1b, '[', 'B'})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case k := <-keyChannel:
		if k.Special != "DOWN" {
			t.Fatalf("invalid special key: %v, expected DOWN", k.Special)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for key event")
	}
}

func TestCaptureKeyboardEventsLeftArrow(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	keyChannel := select5.CaptureKeyboardEvents(r)

	_, err = w.Write([]byte{0x1b, '[', 'D'})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case k := <-keyChannel:
		if k.Special != "LEFT" {
			t.Fatalf("invalid special key: %v, expected LEFT", k.Special)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for key event")
	}
}

func TestCaptureKeyboardEventsRightArrow(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	keyChannel := select5.CaptureKeyboardEvents(r)

	_, err = w.Write([]byte{0x1b, '[', 'C'})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case k := <-keyChannel:
		if k.Special != "RIGHT" {
			t.Fatalf("invalid special key: %v, expected RIGHT", k.Special)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for key event")
	}
}

func TestCaptureKeyboardEventsPageUp(t *testing.T) {
	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	keyChannel := select5.CaptureKeyboardEvents(r)

	// Send the PageUp key escape sequence: ESC [ 5 ~
	_, err = w.Write([]byte{0x1b, '[', '5', '~'})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case k := <-keyChannel:
		if k.Special != "PAGEUP" {
			t.Fatalf("invalid special key: %v, expected PAGEUP", k.Special)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for key event")
	}
}

func TestCaptureKeyboardEventsPageDown(t *testing.T) {
	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	keyChannel := select5.CaptureKeyboardEvents(r)

	// Send the PageDown key escape sequence: ESC [ 6 ~
	_, err = w.Write([]byte{0x1b, '[', '6', '~'})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case k := <-keyChannel:
		if k.Special != "PAGEDOWN" {
			t.Fatalf("invalid special key: %v, expected PAGEDOWN", k.Special)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for key event")
	}
}

func TestCaptureKeyboardEventsHome(t *testing.T) {
	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	keyChannel := select5.CaptureKeyboardEvents(r)

	// Send the Home key escape sequence: ESC [ H
	_, err = w.Write([]byte{0x1b, '[', 'H'})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case k := <-keyChannel:
		if k.Special != "HOME" {
			t.Fatalf("invalid special key: %v, expected HOME", k.Special)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for key event")
	}
}

func TestCaptureKeyboardEventsEnd(t *testing.T) {
	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	keyChannel := select5.CaptureKeyboardEvents(r)

	// Send the End key escape sequence: ESC [ F
	_, err = w.Write([]byte{0x1b, '[', 'F'})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case k := <-keyChannel:
		if k.Special != "END" {
			t.Fatalf("invalid special key: %v, expected END", k.Special)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for key event")
	}
}

func TestCaptureKeyboardEventsCtrlD(t *testing.T) {
	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	keyChannel := select5.CaptureKeyboardEvents(r)

	// Send Ctrl+D (ASCII 4)
	_, err = w.Write([]byte{4})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case k := <-keyChannel:
		// Check that it's the correct character code for Ctrl+D
		if k.Key != 4 {
			t.Fatalf("invalid key code: %d, expected 4", k.Key)
		}

		// Verify Ctrl modifier is set
		if !k.Ctrl {
			t.Fatalf("Ctrl modifier not set for Ctrl+D")
		}

		if k.Special != "" {
			t.Fatalf("unexpected special key set: %v", k.Special)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for key event")
	}
}
