package select5_test

import (
	"github.com/g1eng/select5"
	"log"
	"os"
	"testing"
	"time"
)

func TestCaptureKeyboardEvents(t *testing.T) {

	targetChar := []rune{'a', 'b', 'C', '1', '#', '$', ' '}

	for _, key := range targetChar {
		// Create a pipe
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()
		defer w.Close()

		oldStdin := os.Stdin
		os.Stdin = r
		defer func() { os.Stdin = oldStdin }() // Restore stdin when test completes

		keyChannel := make(chan select5.KeyEvent)
		sigChan := make(chan os.Signal)

		go func() {
			keyChannel, sigChan = select5.CaptureKeyboardEvents()
		}()

		time.Sleep(100 * time.Millisecond)

		_, err = w.Write([]byte{byte(key)})
		if err != nil {
			t.Fatal(err)
		}

		select {
		case k := <-keyChannel:
			if k.Key != key {
				t.Fatalf("invalid key code: %c, expected %c", k.Key, key)
			}
		case sig := <-sigChan:
			t.Fatalf("os.Signal returned: %s", sig.String())
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for key event")
		}

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

	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	keyChannel := make(chan select5.KeyEvent)
	sigChan := make(chan os.Signal)
	go func() {
		keyChannel, sigChan = select5.CaptureKeyboardEvents()
	}()
	time.Sleep(100 * time.Millisecond)

	key := '\n'
	_, err = w.Write([]byte{byte(key)})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case k := <-keyChannel:
		if k.Special != select5.ENTER {
			t.Fatalf("invalid key code: %c, expected DEL", k.Key)
		}
	case sig := <-sigChan:
		t.Fatalf("os.Signal returned: %s", sig.String())
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for key event")
	}
}

func TestCaptureKeyboardEventsSpecialChar(t *testing.T) {

	testTarget := map[int][]byte{
		select5.ENTER:    {'\n'},
		select5.BS:       {'\b'},
		select5.DEL:      {0x7f},
		select5.UP:       {0x1b, '[', 'A'},
		select5.DOWN:     {0x1b, '[', 'B'},
		select5.RIGHT:    {0x1b, '[', 'C'},
		select5.LEFT:     {0x1b, '[', 'D'},
		select5.HOME:     {0x1b, '[', 'H'},
		select5.END:      {0x1b, '[', 'F'},
		select5.PAGEUP:   {0x1b, '[', '5', '~'},
		select5.PAGEDOWN: {0x1b, '[', '6', '~'},
	}
	for keyCode, v := range testTarget {
		// Create a pipe
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()
		defer w.Close()

		oldStdin := os.Stdin
		os.Stdin = r
		defer func() { os.Stdin = oldStdin }()

		keyChannel := make(chan select5.KeyEvent)
		sigChan := make(chan os.Signal)
		go func() {
			keyChannel, sigChan = select5.CaptureKeyboardEvents()
		}()

		time.Sleep(100 * time.Millisecond)

		// Send the UP arrow key escape sequence: ESC [ A
		_, err = w.Write(v)
		if err != nil {
			t.Fatal(err)
		}

		select {
		case k := <-keyChannel:
			if k.Special != keyCode {
				t.Fatalf("invalid special key: %x, expected %x", k.Special, keyCode)
			}
		case sig := <-sigChan:
			t.Fatalf("os.Signal returned: %s", sig.String())
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for key event")
		}
	}
}

func TestCaptureKeyboardEventsNonInterruptCtrl(t *testing.T) {
	testTarget := []byte{select5.CtrlA, select5.CtrlB, select5.CtrlD, select5.CtrlE, select5.CtrlN, select5.CtrlP, select5.CtrlX, select5.CtrlY}
	for i, key := range testTarget {
		// Create a pipe
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()
		defer w.Close()

		oldStdin := os.Stdin
		os.Stdin = r
		defer func() { os.Stdin = oldStdin }()
		keyChannel := make(chan select5.KeyEvent)
		sigChan := make(chan os.Signal)

		go func() {
			keyChannel, sigChan = select5.CaptureKeyboardEvents()
		}()

		time.Sleep(100 * time.Millisecond)

		// Send Ctrl+?
		_, err = w.Write([]byte{key})
		if err != nil {
			t.Fatal(err)
		}

		log.Println(i)
		select {
		case k := <-keyChannel:
			// Check that its the correct character code for Ctrl+?
			if k.Key != rune(key) {
				t.Fatalf("invalid key code: %d, expected %d", k.Key, rune(key))
			}
			// Verify Ctrl modifier is set
			if k.Ctrl != true {
				t.Fatalf("Ctrl modifier not set for Ctrl+%c", rune(key+0x60))
			}
			if k.Special != 0 {
				t.Fatalf("unexpected special key set: %v", k.Special)
			}
		case sig := <-sigChan:
			t.Fatalf("os.Signal returned: %s", sig.String())
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for key event")
		}
	}
}

func TestCaptureKeyboardEventsUtf8Char(t *testing.T) {

	testTarget := map[string][]byte{
		"あ": {0xe3, 0x81, 0x82, 0, 0, 0},
		"猫": {0xe7, 0x8c, 0xab, 0, 0, 0},
	}
	for target, v := range testTarget {
		// Create a pipe
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()
		defer w.Close()

		oldStdin := os.Stdin
		os.Stdin = r
		defer func() { os.Stdin = oldStdin }()

		keyChannel := make(chan select5.KeyEvent)
		sigChan := make(chan os.Signal)
		go func() {
			keyChannel, sigChan = select5.CaptureKeyboardEvents()
		}()

		time.Sleep(100 * time.Millisecond)

		// Send the UP arrow key escape sequence: ESC [ A
		_, err = w.Write(v)
		if err != nil {
			t.Fatal(err)
		}

		select {
		case k := <-keyChannel:
			if s, err := k.Utf8Char(); err != nil {
				t.Fatal(err)
			} else if string(s) != target {
				t.Fatalf("invalid utf-8 character: %s, expected %s", k.Runes, target)
			} else if !k.IsRuneStart {
				t.Fatalf("key should be an UTF-8 character")
			}
		case sig := <-sigChan:
			t.Fatalf("os.Signal returned: %s", sig.String())
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for key event")
		}
	}
}

func TestKeyEvent_Size(t *testing.T) {
	tt := []struct {
		name  string
		event select5.KeyEvent
		want  int
	}{
		{
			name: "alphabet key event",
			event: select5.KeyEvent{
				Key:         'a',
				Code:        int('a'),
				Ctrl:        false,
				Alt:         false,
				Shift:       false,
				Special:     0,
				IsRuneStart: false,
				Runes:       []byte{'a', 0, 0, 0, 0, 0},
			},
			want: 1,
		},
		{
			name: "control-a event",
			event: select5.KeyEvent{
				Key:         0x01,
				Code:        0x01,
				Ctrl:        true,
				Alt:         false,
				Shift:       false,
				Special:     0,
				IsRuneStart: false,
				Runes:       []byte{0x01, 0, 0, 0, 0, 0},
			},
			want: 1,
		},
		{
			name: "arrow up event",
			event: select5.KeyEvent{
				Key:         select5.ESC,
				Code:        select5.ESC,
				Ctrl:        true,
				Alt:         false,
				Shift:       false,
				Special:     select5.UP,
				IsRuneStart: false,
				Runes:       []byte{select5.ESC, '[', 'A', 0, 0, 0},
			},
			want: 3,
		},
		{
			name: "UTF-8 event `あ`",
			event: select5.KeyEvent{
				Key:         select5.ESC,
				Code:        select5.ESC,
				Ctrl:        true,
				Alt:         false,
				Shift:       false,
				Special:     select5.UP,
				IsRuneStart: false,
				Runes:       []byte{0xe3, 0x81, 0x82, 0, 0, 0},
			},
			want: 3,
		},
		{
			name: "UTF-8 event `猫`",
			event: select5.KeyEvent{
				Key:         select5.ESC,
				Code:        select5.ESC,
				Ctrl:        true,
				Alt:         false,
				Shift:       false,
				Special:     select5.UP,
				IsRuneStart: false,
				Runes:       []byte{0xe7, 0x8c, 0xab, 0, 0, 0},
			},
			want: 3,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.event.Size(); got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestKeyEvent_String(t *testing.T) {
	tt := []struct {
		name      string
		event     select5.KeyEvent
		want      string
		wantError bool
	}{
		{
			name: "alphabet key event",
			event: select5.KeyEvent{
				Key:         'a',
				Code:        int('a'),
				Ctrl:        false,
				Alt:         false,
				Shift:       false,
				Special:     0,
				IsRuneStart: false,
				Runes:       []byte{'a', 0, 0, 0, 0, 0},
			},
			want:      "a",
			wantError: false,
		},
		{
			name: "control-a event",
			event: select5.KeyEvent{
				Key:         0x01,
				Code:        0x01,
				Ctrl:        true,
				Alt:         false,
				Shift:       false,
				Special:     0,
				IsRuneStart: false,
				Runes:       []byte{0x01, 0, 0, 0, 0, 0},
			},
			wantError: true,
		},
		{
			name: "arrow up event",
			event: select5.KeyEvent{
				Key:         select5.ESC,
				Code:        select5.ESC,
				Ctrl:        true,
				Alt:         false,
				Shift:       false,
				Special:     select5.UP,
				IsRuneStart: false,
				Runes:       []byte{select5.ESC, '[', 'A', 0, 0, 0},
			},
			wantError: true,
		},
		{
			name: "UTF-8 event `あ`",
			event: select5.KeyEvent{
				Key:         0xe3,
				Code:        0xe3,
				Ctrl:        false,
				Alt:         false,
				Shift:       false,
				Special:     0,
				IsRuneStart: true,
				Runes:       []byte{0xe3, 0x81, 0x82, 0, 0, 0},
			},
			want:      "あ",
			wantError: false,
		},
		{
			name: "UTF-8 event `猫`",
			event: select5.KeyEvent{
				Key:         0xe7,
				Code:        0xe7,
				Ctrl:        false,
				Alt:         false,
				Shift:       false,
				Special:     0,
				IsRuneStart: true,
				Runes:       []byte{0xe7, 0x8c, 0xab, 0, 0, 0},
			},
			want:      "猫",
			wantError: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantError {
				if _, err := tc.event.Utf8Char(); err == nil {
					t.Errorf("want error, got nil")
				}
			} else {
				if got, err := tc.event.Utf8Char(); err != nil {
					t.Fatal(err)
				} else if string(got) != tc.want {
					log.Println("got size: ", len(got))
					t.Errorf("got %s (%x), want %s (%x%x%x)", string(got), got[0], tc.want, []byte(tc.want)[0], []byte(tc.want)[1], []byte(tc.want)[2])
				}
			}
		})
	}
}
