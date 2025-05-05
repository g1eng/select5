package main

import (
	"fmt"
	"github.com/g1eng/select5"
	"io"
	"log"
	"os"
	"strings"
	"syscall"
	"unicode/utf8"
)

type CursorPosition struct {
	X, Y int
}

type Editor struct {
	Cursor      CursorPosition
	In          io.Reader
	Out         io.Writer
	Line        []string
	PreviousKey select5.KeyEvent
}

func NewEditor() *Editor {
	return &Editor{
		CursorPosition{0, 0},
		os.Stdin,
		os.Stdout,
		[]string{""},
		select5.KeyEvent{
			IsRuneStart: false,
		},
	}
}

func (e *Editor) Edit() string {
	fmt.Fprint(e.Out, select5.HideCursor)
	fmt.Fprint(e.Out, select5.ClearScreen)
	fmt.Fprint(e.Out, select5.ResetCursor)
	fmt.Fprint(e.Out, select5.ShowCursor)

	keyCh, sigCh := select5.CaptureKeyboardEvents()
	for {
		select {
		case sig := <-sigCh:
			switch sig {
			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
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
						e.PreviousKey.IsRuneStart = true
					} else {
						e.PutS([]byte{byte(key.Key)})
						e.PreviousKey.IsRuneStart = false
					}
				} else if key.Ctrl == true {
					switch key.Key {
					case 0x01:
						e.Cursor.X = 0
					case 0x03:
						//Ctrl-C -> end
						fmt.Fprint(e.Out, select5.ClearScreen, select5.ResetCursor)
						os.Exit(0)
						return ""
					case 0x04:
						//Ctrl-D -> end without clear screen
						fmt.Fprint(e.Out, select5.ResetCursor)
						return strings.Join(e.Line, "\n")
					case 0x05:
						e.Cursor.X = e.GetLineMaxX()
					}
					e.Reposition()
				}
			case select5.BS:
				e.PutBackspace()
			case select5.DEL:
				e.PutDelete()
			case select5.ENTER:
				e.PutEnter()
			case select5.UP:
				e.Up()
			case select5.DOWN:
				e.Down()
			case select5.RIGHT:
				e.Right()
			case select5.LEFT:
				e.Left()
			}
		}
	}
}

func (e *Editor) IsOnBlankLine() bool {
	return len(e.Line[e.Cursor.Y]) == 0
}

func (e *Editor) IsNotOnLastLine() bool {
	return e.Cursor.Y < len(e.Line)-1
}

func (e *Editor) IsOnLineHead() bool {
	return e.Cursor.X == 0
}

func (e *Editor) IsOnUtf8Boundary() bool {
	if e.IsOnLineEnd() {
		return true
	}
	return utf8.RuneStart(e.Line[e.Cursor.Y][e.Cursor.X])
}

func (e *Editor) IsOnLineEnd() bool {
	return e.Cursor.X >= e.GetLineMaxX()
}

func (e *Editor) GetCurrentLine() string {
	return e.Line[e.Cursor.Y]
}

func (e *Editor) GetCurrentLineLength() int {
	return len(e.Line[e.Cursor.Y])
}
func (e *Editor) GetLineMaxX() int {
	return len(e.Line[e.Cursor.Y])
}
func (e *Editor) GetTextMaxY() int {
	return len(e.Line) - 1
}

func (e *Editor) Reposition() {
	fmt.Fprintf(e.Out, select5.MoveTo, e.Cursor.Y+1, e.GetLineVisibleXPosition())
}
func (e *Editor) GoToLineHead() {
	fmt.Fprintf(e.Out, select5.MoveTo, e.Cursor.Y+1, 1)
}

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
		fmt.Fprint(e.Out, str, select5.ClearLineFromCursor, e.GetCurrentLine())
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
		fmt.Fprint(e.Out, select5.ClearLineFromCursor, string(key), post)
	}
	e.Cursor.X += len(key)
	e.Reposition()
}

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

	// Update cursor position
	e.Cursor.Y++
	e.Cursor.X = 0

	// Update the display (this part is not needed for the test but kept for actual use)
	fmt.Fprintf(e.Out, select5.MoveTo, e.Cursor.Y+1, 1)
}

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
	fmt.Fprintf(e.Out, select5.MoveTo, e.Cursor.Y+1, e.Cursor.X+1)
}

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

func (e *Editor) Left() {
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

func (e *Editor) IsDocumentHead() bool {
	return e.IsOnLineHead() && e.Cursor.Y == 0
}

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
			bs := func() {
				e.Line[e.Cursor.Y] = e.Line[e.Cursor.Y][:e.Cursor.X-1] + e.Line[e.Cursor.Y][e.Cursor.X:]
				e.Cursor.X--
			}
			e.Cursor.X--
			for !e.IsOnUtf8Boundary() {
				e.Cursor.X++
				bs()
				e.Cursor.X--
			}
			e.Cursor.X++
			bs()
		}
	}
}

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

func (e *Editor) GetLineVisibleXPosition() int {
	count := 0
	line := e.Line[e.Cursor.Y]
	continuousBit := false
	for i, r := range []byte(line) {
		if i == e.Cursor.X {
			break
		}
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
