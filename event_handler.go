package select5

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"os/signal"
	"syscall"
	"unicode/utf8"
)

// KeyEvent represents a keyboard input event with information about special keys and modifiers.
type KeyEvent struct {
	Key         rune   // The character pressed
	Code        int    // Numeric key code
	Ctrl        bool   // Whether Ctrl was pressed
	Alt         bool   // Whether Alt was pressed
	Shift       bool   // Whether Shift was pressed
	Special     int    // Special key name (UP, DOWN, ENTER, etc.)
	IsRuneStart bool   // Whether the character is UTF-8 multibyte character or not
	Runes       []byte // Raw key bytes
}

// Utf8Char returns byte representation for the UTF-8 character.
// If it is an ascii character, simply return the first byte.
func (e KeyEvent) Utf8Char() ([]byte, error) {
	if e.Ctrl {
		return nil, fmt.Errorf("control characters cannot be formated")
	}
	if e.IsRuneStart {
		for i := range e.Runes {
			if e.Runes[i] == 0 {
				return e.Runes[:i], nil
			}
		}
		return e.Runes, nil
	} else {
		return []byte{byte(e.Key)}, nil
	}
}

// Size presents the size in octet order of the UTF-8 character.
func (e KeyEvent) Size() (size int) {
	for i, b := range e.Runes {
		if b == 0 {
			break
		}
		size = i
	}
	return size + 1
}

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

		buffer := make([]byte, 0, 8)
		oneByte := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(oneByte)
			if err != nil {
				close(keyChannel)
				return
			}

			if n > 0 {
				buffer = append(buffer, oneByte[0])
				//key := KeyEvent{
				//	Key:   rune(oneByte[0]),
				//	Code:  int(oneByte[0]),
				//	Runes: oneByte,
				//}
				if buffer[0] == ESC {
					var escapedKey int
					// Check for special keys or modifiers
					// Check for arrow keys and other special keys
					if len(buffer) >= 3 && buffer[1] == '[' {
						complete := false
						if len(buffer) == 3 {
							switch buffer[2] {
							case 'A', 'B', 'C', 'D', 'F', 'H':
								complete = true
							}
						}
						if len(buffer) == 4 && buffer[3] == '~' {
							complete = true
						}
						if complete {
							switch buffer[2] {
							case 'A':
								escapedKey = UP
							case 'B':
								escapedKey = DOWN
							case 'C':
								escapedKey = RIGHT
							case 'D':
								escapedKey = LEFT
							case '5':
								escapedKey = PAGEUP
							case '6':
								escapedKey = PAGEDOWN
							case 'F':
								escapedKey = END
							case 'H':
								escapedKey = HOME
							}
							keyCode := int(buffer[0]<<4) | int(buffer[1])<<2 | int(buffer[2])
							keyChannel <- KeyEvent{
								Key:     rune(buffer[1]),
								Code:    keyCode,
								Special: escapedKey,
							}
							//clear
							buffer = buffer[:0]
							continue
						}

					}
					if len(buffer) < 4 {
						continue
					}
					keyCode := int(buffer[0]<<4) | int(buffer[1])<<2 | int(buffer[2])
					keyChannel <- KeyEvent{
						Key:     rune(buffer[1]),
						Code:    keyCode,
						Special: escapedKey,
					}
					//clear
					buffer = buffer[:0]
					continue
				}
				if buffer[0] == ENTER || buffer[0] == 0x0d {
					keyChannel <- KeyEvent{
						Key:     rune(ENTER),
						Code:    ENTER,
						Special: ENTER,
						Runes:   buffer[:1],
					}
					buffer = buffer[:0] //clear
					continue
				} else if buffer[0]&0x80 != 0x80 {
					//ascii character handling
					receiveASCII(buffer[0], keyChannel, sigChan)
					buffer = buffer[:0] //clear
					continue
				}
				if utf8.FullRune(buffer) {
					r, size := utf8.DecodeRune(buffer)
					if r != utf8.RuneError {
						key := KeyEvent{
							Key:         r,
							Code:        int(buffer[0]),
							IsRuneStart: true,
							Runes:       make([]byte, 6), // TODO: check maxSize = 6?
						}
						copy(key.Runes, buffer[:size])
						keyChannel <- key
						buffer = buffer[size:] // remove processed bytes

						// if we had one or more character,
						if len(buffer) > 0 {
							// remaining bytes back for processing
							remainder := make([]byte, len(buffer))
							copy(remainder, buffer)
							buffer = buffer[:0] //clear

							// process each remaining byte
							for _, b := range remainder {
								buffer = append(buffer, b)
								// Try to process if it's a complete sequence
								if utf8.FullRune(buffer) || buffer[0] < 128 {
									break //reset
								}
							}
							continue
						}

					} else if size > 0 {
						// invalid UTF-8 sequence, but we consumed some bytes
						buffer = buffer[size:] // skip
						continue
					}
				}
				if len(buffer) < 4 {
					continue
				}

				// If buffer is getting too large but no valid character,
				// emit what we have as individual bytes (likely garbage)
				if len(buffer) > 6 {
					for _, b := range buffer {
						receiveASCII(b, keyChannel, sigChan)
					}
					buffer = buffer[:0] // Clear buffer
				}
			}

		}
	}()
	return keyChannel, sigChan
}

// receive ASCII and control characters
func receiveASCII(b byte, keyChannel chan KeyEvent, sigChan chan os.Signal) {
	key := KeyEvent{
		Key:   rune(b),
		Code:  int(b),
		Runes: make([]byte, 6),
	}
	key.Runes[0] = b

	switch b {
	case BS:
		key.Special = BS
	case DEL:
		key.Special = DEL
	default:
		if b < 0x20 {
			// Control characters
			key.Ctrl = true
			switch b {
			case CtrlC:
				sigChan <- syscall.SIGINT
			case CtrlZ:
				sigChan <- syscall.SIGSTOP
			}
		}
	}

	keyChannel <- key

	// Handle Ctrl+C separately to ensure we exit
	if b == 0x03 { // Ctrl+C
		close(keyChannel)
	}
}
