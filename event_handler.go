package select5

import (
	"golang.org/x/term"
	"os"
)

// CaptureKeyboardEvents starts capturing keyboard events in a background goroutine.
// Returns a channel that delivers KeyEvent structs.
// The channel should be properly handled and the goroutine will exit when the channel is closed.
func CaptureKeyboardEvents(stdin *os.File) chan KeyEvent {
	keyChannel := make(chan KeyEvent, 10)
	go func() {
		// Check if stdin is a terminal and set up raw mode if it is
		var oldState *term.State
		isTerm := term.IsTerminal(int(stdin.Fd()))
		if isTerm {
			var err error
			oldState, err = term.MakeRaw(int(stdin.Fd()))
			if err != nil {
				close(keyChannel)
				return
			}
			defer term.Restore(int(stdin.Fd()), oldState)
		}

		buffer := make([]byte, 6)
		for {
			n, err := stdin.Read(buffer)
			if err != nil {
				close(keyChannel)
				return
			}

			if n > 0 {
				key := KeyEvent{
					Key:  rune(buffer[0]),
					Code: int(buffer[0]),
				}

				// Check for special keys or modifiers
				if buffer[0] == 27 { // Escape sequence
					// Check for arrow keys and other special keys
					if n >= 3 && buffer[1] == '[' {
						switch buffer[2] {
						case 'A':
							key.Special = "UP"
						case 'B':
							key.Special = "DOWN"
						case 'C':
							key.Special = "RIGHT"
						case 'D':
							key.Special = "LEFT"
						case '5':
							key.Special = "PAGEUP"
						case '6':
							key.Special = "PAGEDOWN"
						case 'F':
							key.Special = "END"
						case 'H':
							key.Special = "HOME"
						}
					}
					key.Alt = true
				} else if buffer[0] == 13 || buffer[0] == 10 {
					// Handle Enter key
					key.Special = "ENTER"
				} else if buffer[0] < 32 {
					// Control characters
					key.Ctrl = true
				}

				keyChannel <- key
			}

			if buffer[0] == 3 { // Ctrl+C
				close(keyChannel)
				return
			}
		}
	}()
	return keyChannel
}
