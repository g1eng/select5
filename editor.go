package select5

import (
	"fmt"
	"golang.org/x/term"
	"io"
	"log"
	"os"
	"strings"
	"syscall"
	"unicode/utf8"
)

// CursorPosition represents the X,Y coordinates of the cursor in the editor
type CursorPosition struct {
	X, Y int
}

// Editor provides a simple terminal-based text editor
type Editor struct {
	Cursor CursorPosition
	In     io.Reader
	Out    io.Writer
	Line   []string
}

// NewEditor creates a new Editor instance with default settings
func NewEditor() *Editor {
	return &Editor{
		CursorPosition{0, 0},
		os.Stdin,
		os.Stdout,
		[]string{""},
	}
}

// Edit starts the editing session and returns the edited text when complete (with Ctrl-D)
func (e *Editor) Edit() string {
	fmt.Fprint(e.Out, HideCursor)
	fmt.Fprint(e.Out, ClearScreen)
	fmt.Fprint(e.Out, ResetCursor)
	fmt.Fprint(e.Out, ShowCursor)

	keyCh, sigCh := CaptureKeyboardEvents()
	var oldState *term.State
	isTerm := term.IsTerminal(int(os.Stdin.Fd()))
	if isTerm {
		var err error
		oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			close(keyCh)
			return ""
		}

		defer term.Restore(int(os.Stdin.Fd()), oldState)
		defer term.Restore(int(os.Stdout.Fd()), oldState)
	}
	for {
		select {
		case sig := <-sigCh:
			switch sig {
			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
				term.Restore(int(os.Stdin.Fd()), oldState)
				term.Restore(int(os.Stdout.Fd()), oldState)
				os.Exit(0)
			}
		case key := <-keyCh:

			switch key.Special {
			case 0:
				if key.Ctrl == false {
					if key.IsRuneStart {
						runes, err := key.Utf8Char()
						if err != nil {
							log.Fatal(err)
						}
						e.PutS(runes)
					} else {
						e.PutS([]byte{byte(key.Key)})
					}
				} else if key.Ctrl == true {
					switch key.Key {
					case CtrlA:
						e.Cursor.X = 0
					case CtrlD:
						//Ctrl-D -> end without clear screen
						fmt.Fprint(e.Out, ResetCursor)
						return strings.Join(e.Line, "\n")
					case CtrlE:
						e.Cursor.X = e.GetLineMaxX()
					case CtrlP:
						e.Up()
					case CtrlJ:
						e.PutEnter()
					case CtrlN:
						e.Down()
					}
					e.Reposition()
				}
			case BS:
				fallthrough
			case DEL:
				e.PutBackspace()
			case ENTER:
				e.PutEnter()
			case UP:
				e.Up()
			case DOWN:
				e.Down()
			case RIGHT:
				e.Right()
			case LEFT:
				e.Left()
			}
		}
	}
}

// IsOnBlankLine returns true if the current line is empty
func (e *Editor) IsOnBlankLine() bool {
	return len(e.Line[e.Cursor.Y]) == 0
}

// IsNotOnLastLine returns true if the cursor is not on the last line
func (e *Editor) IsNotOnLastLine() bool {
	return e.Cursor.Y < len(e.Line)-1
}

// IsOnLineHead returns true if the cursor is at the beginning of a line
func (e *Editor) IsOnLineHead() bool {
	return e.Cursor.X == 0
}

// IsOnUtf8Boundary returns true if the cursor is on a valid UTF-8 character boundary
func (e *Editor) IsOnUtf8Boundary() bool {
	if e.IsOnLineEnd() {
		return true
	}
	return utf8.RuneStart(e.Line[e.Cursor.Y][e.Cursor.X])
}

// IsOnLineEnd returns true if the cursor is at the end of a line
func (e *Editor) IsOnLineEnd() bool {
	return e.Cursor.X >= e.GetLineMaxX()
}

// GetCurrentLine returns the text of the current line
func (e *Editor) GetCurrentLine() string {
	return e.Line[e.Cursor.Y]
}

// GetCurrentLineLength returns the length of the current line in bytes
func (e *Editor) GetCurrentLineLength() int {
	return len(e.Line[e.Cursor.Y])
}

// GetLineMaxX returns the maximum valid X position for the current line
func (e *Editor) GetLineMaxX() int {
	return len(e.Line[e.Cursor.Y])
}

// GetTextMaxY returns the maximum valid Y position (last line index)
func (e *Editor) GetTextMaxY() int {
	return len(e.Line) - 1
}

// Reposition moves the terminal cursor to match the editor's cursor position
func (e *Editor) Reposition() {
	fmt.Fprintf(e.Out, MoveTo, e.Cursor.Y+1, e.GetLineVisibleXPosition())
}

// GoToLineHead moves the cursor to the beginning of the current line
func (e *Editor) GoToLineHead() {
	fmt.Fprintf(e.Out, MoveTo, e.Cursor.Y+1, 1)
}

// PutS inserts a string at the current cursor position
func (e *Editor) PutS(key []byte) {
	//log.Fprintf(e.Out,"%x\n", key)
	var str string
	if len(key) == 0 {
		panic("key is empty")
	} else {
		str = string(key)
	}

	if e.IsOnBlankLine() || e.IsOnLineEnd() {
		e.Line[e.Cursor.Y] += str
		fmt.Fprint(e.Out, str)
	} else if e.IsOnLineHead() {
		fmt.Fprint(e.Out, str, ClearLineFromCursor, e.GetCurrentLine())
		e.Line[e.Cursor.Y] = str + e.GetCurrentLine()
	} else {
		var pre, post string
		pre = e.GetCurrentLine()[:e.Cursor.X]
		if e.IsOnLineEnd() {
			post = ""
		} else {
			post = e.GetCurrentLine()[e.Cursor.X:]
		}
		e.Line[e.Cursor.Y] = pre + string(key) + post
		fmt.Fprint(e.Out, ClearLineFromCursor, string(key), post)
	}
	e.Cursor.X += len(key)
	e.Reposition()
}

// PutEnter inserts a new line at the current cursor position
func (e *Editor) PutEnter() {
	// Determine the content for the current line and the new line
	var currentLineContent, newLineContent string
	//
	//if !e.IsOnUtf8Boundary() {
	//	for !e.IsOnUtf8Boundary() {
	//		e.Cursor.X--
	//	}
	//}

	if e.IsOnLineHead() {
		currentLineContent = ""
		newLineContent = e.GetCurrentLine()
	} else if e.IsOnLineEnd() {
		currentLineContent = e.GetCurrentLine()
		newLineContent = ""
	} else {
		originalLine := e.GetCurrentLine()
		for i, c := range originalLine {
			if i < e.Cursor.X {
				currentLineContent += string(c)
			} else {
				newLineContent += string(c)
			}
		}
	}

	if e.IsNotOnLastLine() {
		newLines := make([]string, len(e.Line)+1)
		for i := range newLines {
			if i < e.Cursor.Y+1 {
				newLines[i] = e.Line[i]
			} else if i == e.Cursor.Y+1 {
				newLines[i] = newLineContent
			} else {
				newLines[i] = e.Line[i-1]
			}
		}
		e.Line = newLines
	} else {
		// Insert at the end of the document
		e.Line = append(e.Line, newLineContent)
	}
	// Update the current line
	e.Line[e.Cursor.Y] = currentLineContent

	fmt.Fprintf(e.Out, ClearScreenFromCursor)
	// Update cursor position
	e.Cursor.Y++
	e.Cursor.X = 0
	e.Reposition()

	fmt.Fprintf(e.Out, MoveTo, e.Cursor.Y, 1)
	fmt.Fprint(e.Out, e.Line[e.Cursor.Y-1])
	for i, line := range e.Line[e.Cursor.Y:] {
		fmt.Fprintf(e.Out, MoveTo, e.Cursor.Y+i+1, 1)
		fmt.Fprint(e.Out, line)
	}

	fmt.Fprintf(e.Out, MoveTo, e.Cursor.Y+1, 1)
}

// Up moves the cursor up one line, adjusting X position if needed
func (e *Editor) Up() {
	if e.Cursor.Y > 0 {
		e.Cursor.Y--
		if e.IsOnBlankLine() {
			e.Cursor.X = 0
		} else if e.IsOnLineEnd() {
			e.Cursor.X = e.GetLineMaxX()
		} else if !e.IsOnUtf8Boundary() {
			for !e.IsOnUtf8Boundary() {
				e.Cursor.X--
			}
		}
	}
	e.Reposition()
}

// Down moves the cursor down one line, adjusting X position if needed
func (e *Editor) Down() {
	if e.IsNotOnLastLine() {
		e.Cursor.Y++
		if e.IsOnBlankLine() {
			e.Cursor.X = 0
		} else if e.Cursor.X > e.GetLineMaxX() {
			e.Cursor.X = e.GetLineMaxX()
		} else if !e.IsOnUtf8Boundary() {
			for !e.IsOnUtf8Boundary() {
				e.Cursor.X--
			}
		}
	}
	e.Reposition()
}

// Right moves the cursor one position to the right, respecting UTF-8 boundaries
func (e *Editor) Right() {
	if !e.IsOnLineEnd() {
		e.Cursor.X++
		if !e.IsOnUtf8Boundary() {
			for !e.IsOnUtf8Boundary() {
				e.Cursor.X++
			}
		}
	}
	e.Reposition()
}

// Left moves the cursor one position to the right, respecting UTF-8 boundaries
func (e *Editor) Left() {
	if e.IsOnBlankLine() {
		return
	} else if e.IsOnLineHead() {
		e.Cursor.Y--
		e.Cursor.X = e.GetLineMaxX()
	}
	if e.Cursor.X > 0 {
		e.Cursor.X--
		if !e.IsOnUtf8Boundary() {
			for !e.IsOnUtf8Boundary() {
				e.Cursor.X--
			}
		}
	}
	e.Reposition()
}

// IsDocumentHead returns true if the cursor at the beginning of the document
func (e *Editor) IsDocumentHead() bool {
	return e.IsOnLineHead() && e.Cursor.Y == 0
}

// PutDelete removes the character at the current cursor position
func (e *Editor) PutDelete() {
	if !e.IsOnLineEnd() && !e.IsOnBlankLine() {
		if e.IsOnLineHead() {
			del := func() { e.Line[e.Cursor.Y] = e.Line[e.Cursor.Y][1:] }
			del()
			if len(e.Line[e.Cursor.Y]) != 0 && !e.IsOnUtf8Boundary() {
				for !e.IsOnUtf8Boundary() {
					del()
				}
			}
		} else {
			del := func() {
				post := e.Line[e.Cursor.Y][e.Cursor.X+1:]
				//log.Println("post:", post)
				e.Line[e.Cursor.Y] = e.Line[e.Cursor.Y][:e.Cursor.X] + post
			}
			del()
			for !e.IsOnUtf8Boundary() {
				del()
			}
		}
	}
}

// PutBackspace removes the previous character at the current cursor position.
// If it is performed at the beginning of a line, two lines are concatenated to one.
func (e *Editor) PutBackspace() {
	if !e.IsDocumentHead() {
		if e.IsOnLineHead() {
			pre, postLines := e.Line[:e.Cursor.Y], e.Line[e.Cursor.Y:]
			nextX := len(pre[len(pre)-1])
			pre[len(pre)-1] = pre[len(pre)-1] + postLines[0]
			e.Line = append(pre, postLines[1:]...)
			e.Cursor.Y--
			e.Cursor.X = nextX
		} else {
			tmpX := e.Cursor.X - 1
			for !utf8.RuneStart(e.Line[e.Cursor.Y][tmpX]) {
				tmpX--
			}
			e.Line[e.Cursor.Y] = e.Line[e.Cursor.Y][:tmpX] + e.Line[e.Cursor.Y][e.Cursor.X:]
			e.Cursor.X = tmpX
		}
		e.Reposition()
		fmt.Fprintf(e.Out, ClearLineFromCursor)
		fmt.Fprintf(e.Out, "%s", e.Line[e.Cursor.Y][e.Cursor.X:])
		e.Reposition()
	}
}

// GetLineVisibleLength returns the visible line length on the console.
// It can be used to calculate the real width for the line, which
// consists of multibyte (and double-width) characters.
func (e *Editor) GetLineVisibleLength() int {
	count := 0
	line := e.Line[e.Cursor.Y]
	continuousBit := false
	for _, r := range []byte(line) {
		if utf8.RuneStart(r) {
			continuousBit = false
			count++
		} else {
			if !continuousBit {
				continuousBit = true
				count++
			}
		}
	}
	return count
}

// GetLineVisibleXPosition returns the visible cursor position on the console.
// It can be used to calculate the real cursor position for the line, which
// consists of multibyte (and double-width) characters.
func (e *Editor) GetLineVisibleXPosition() int {
	count := 0
	line := e.Line[e.Cursor.Y]
	continuousBit := false
	if e.IsOnBlankLine() {
		return 0
	}
	for i, r := range []byte(line) {
		if i == e.Cursor.X {
			break
		}
		if utf8.RuneStart(r) {
			continuousBit = false
			count++
		} else {
			//add just 1 width for multibyte characters
			if !continuousBit {
				count++
			}
			continuousBit = true
		}
	}
	return count + 1
}
