package select5

import (
	"golang.org/x/term"
	"os"
	"os/signal"
	"syscall"
)

// CaptureKeyboardEvents starts capturing keyboard events in a background goroutine.
// Returns a channel that delivers KeyEvent structs.
// The channel should be properly handled and the goroutine will exit when the channel is closed.
func CaptureKeyboardEvents() (chan KeyEvent, chan os.Signal) {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGCONT, syscall.SIGQUIT)
	keyChannel := make(chan KeyEvent, 10)

	go func() {
		// Check if os.Stdin is a terminal and set up raw mode if it is
		var oldState *term.State
		isTerm := term.IsTerminal(int(os.Stdin.Fd()))
		if isTerm {
			var err error
			oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
			if err != nil {
				close(keyChannel)
				return
			}
			defer term.Restore(int(os.Stdin.Fd()), oldState)
			defer term.Restore(int(os.Stdout.Fd()), oldState)
		}

		buffer := make([]byte, 6)
		for {
			n, err := os.Stdin.Read(buffer)
			if err != nil {
				close(keyChannel)
				return
			}

			if n > 0 {
				key := KeyEvent{
					Key:  rune(buffer[0]),
					Code: int(buffer[0]),
				}
				switch key.Key {
				case ESC: // Escape sequence
					// Check for special keys or modifiers
					// Check for arrow keys and other special keys
					if n >= 3 && buffer[1] == '[' {
						keyCode := int(buffer[0]<<4) | int(buffer[1])<<2 | int(buffer[2])
						if keyCode&0xffff != keyCode {
							key.Code = keyCode
						}
						switch buffer[2] {
						case 'A':
							key.Special = UP
						case 'B':
							key.Special = DOWN
						case 'C':
							key.Special = RIGHT
						case 'D':
							key.Special = LEFT
						case '5':
							key.Special = PAGEUP
						case '6':
							key.Special = PAGEDOWN
						case 'F':
							key.Special = END
						case 'H':
							key.Special = HOME
						}
					}
					key.Alt = true
				case 0x0d:
					fallthrough
				case ENTER:
					// handle enter key
					key.Special = ENTER
				case BS:
					key.Special = BS
				case DEL:
					key.Special = DEL
				default:
					if key.Key < 0x20 {
						// control characters
						key.Ctrl = true
						switch key.Key {
						case CtrlC:
							sigChan <- syscall.SIGINT
						case CtrlZ:
							sigChan <- syscall.SIGSTOP
						}
					}
				}
				keyChannel <- key
			}

			if buffer[0] == 0x03 { // Ctrl+C
				close(keyChannel)
				return
			}
		}
	}()
	return keyChannel, sigChan
}
