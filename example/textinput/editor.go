package main

import (
	"fmt"
	"github.com/g1eng/select5"
	"os"
	"strings"
	"syscall"
)

type CursorPosition struct {
	X, Y int
}

type Editor struct {
	Cursor CursorPosition
	Line   []string
}

func NewEditor() *Editor {
	return &Editor{
		CursorPosition{0, 0},
		[]string{""},
	}
}

func (e *Editor) Edit() string {
	fmt.Print(select5.HideCursor)
	fmt.Print(select5.ClearScreen)
	fmt.Print(select5.ResetCursor)
	fmt.Print(select5.ShowCursor)

	keyCh, sigCh := select5.CaptureKeyboardEvents()
	for {
		select {
		case sig := <-sigCh:
			switch sig {
			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
				os.Exit(0)
				//case syscall.SIGSTOP:
				//	fmt.Print(select5.ClearScreen)
				//	fmt.Print(select5.ShowCursor)
				//	proc, _ := os.FindProcess(os.Getpid())
				//	err := proc.Signal(syscall.SIGSTOP)
				//	if err != nil {
				//		panic(err)
				//	}
				//case syscall.SIGCONT:
				//	fmt.Print(select5.ClearScreen)
				//	fmt.Print(select5.ResetCursor)
				//	for i, s := range e.Line {
				//		fmt.Print(s)
				//		fmt.Printf(select5.MoveTo, i+1, 1)
				//	}
				//	e.Reposition()
			}
		case key := <-keyCh:

			switch key.Special {
			case 0:
				if key.Ctrl == false {
					if e.GetCurrentLineLength() == 0 || e.Cursor.X == e.GetLineMaxX() {
						e.Line[e.Cursor.Y] += string(key.Key)
						fmt.Print(string(key.Key))
					} else if e.Cursor.X == 0 {
						e.Line[e.Cursor.Y] = string(key.Key) + e.GetCurrentLine()
						fmt.Print(select5.ClearLineFromCursor, e.GetCurrentLine())
					} else {
						pre, post := e.GetCurrentLine()[:e.Cursor.X], e.GetCurrentLine()[e.Cursor.X:]
						e.Line[e.Cursor.Y] = pre + string(key.Key) + post
						fmt.Print(select5.ClearLineFromCursor, string(key.Key), post)
					}
					e.Cursor.X++
					e.Reposition()
				} else if key.Ctrl == true {
					switch key.Key {
					case 0x01:
						e.Cursor.X = 0
					case 0x03:
						//Ctrl-C -> end
						fmt.Print(select5.ClearScreen, select5.ResetCursor)
						os.Exit(0)
						return ""
					case 0x04:
						//Ctrl-D -> end without clear screen
						fmt.Print(select5.ResetCursor)
						return strings.Join(e.Line, "\n")
					case 0x05:
						e.Cursor.X = e.GetLineMaxX()
					}
					e.Reposition()
				}
			case select5.ENTER:
				var insert string
				if e.GetCurrentLineLength() == 0 || e.Cursor.X == e.GetLineMaxX() {
					insert = ""
				} else {
					pre, post := e.GetCurrentLine()[:e.Cursor.X], e.GetCurrentLine()[e.Cursor.X:]
					insert = post
					e.Line[e.Cursor.Y] = pre
				}
				if e.Cursor.Y == e.GetTextMaxY() {
					fmt.Print(insert)
					e.Line = append(e.Line, insert)
				} else {
					postLines := e.Line[e.Cursor.Y+1:]
					fmt.Printf(select5.MoveTo, e.Cursor.Y+1, 1)
					fmt.Print(select5.ClearScreenFromCursor)
					fmt.Print(e.GetCurrentLine())
					e.Line = append(e.Line[:e.Cursor.Y], e.GetCurrentLine(), insert)
					e.Line = append(e.Line, postLines...)
					for i, s := range e.Line[e.Cursor.Y:] {
						fmt.Printf(select5.MoveTo, e.Cursor.Y+1+i, 1)
						fmt.Print(s)
					}
				}
				e.Cursor.X = 0
				e.Cursor.Y++
				fmt.Printf(select5.MoveTo, e.Cursor.Y+1, 1)
			case select5.UP:
				if e.Cursor.Y > 0 {
					e.Cursor.Y--
					if e.IsOnBlankLine() {
						e.Cursor.X = 0
					} else if e.Cursor.X > e.GetLineMaxX() {
						e.Cursor.X = e.GetLineMaxX()
					}
				}
				fmt.Printf(select5.MoveTo, e.Cursor.Y+1, e.Cursor.X+1)
			case select5.DOWN:
				if e.IsNotOnLastLine() {
					e.Cursor.Y++
					if e.IsOnBlankLine() {
						e.Cursor.X = 0
					} else if e.Cursor.X > e.GetLineMaxX() {
						e.Cursor.X = e.GetLineMaxX()
					}
				}
				e.Reposition()
			case select5.RIGHT:
				if e.IsOnLastChar() {
					e.Cursor.X++
				}
				e.Reposition()
			case select5.LEFT:
				if e.Cursor.X > 0 {
					e.Cursor.X--
				}
				e.Reposition()
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

func (e *Editor) IsOnLastChar() bool {
	return e.Cursor.X < e.GetLineMaxX()
}

func (e *Editor) GetCurrentLine() string {
	return e.Line[e.Cursor.Y]
}

func (e *Editor) GetCurrentLineLength() int {
	return len(e.Line[e.Cursor.Y])
}
func (e *Editor) GetLineMaxX() int {
	return len(e.Line[e.Cursor.Y]) - 1
}
func (e *Editor) GetTextMaxY() int {
	return len(e.Line) - 1
}

func (e *Editor) Reposition() {
	fmt.Printf(select5.MoveTo, e.Cursor.Y+1, e.Cursor.X+1)
}
func (e *Editor) GoToLineHead() {
	fmt.Printf(select5.MoveTo, e.Cursor.Y+1, 1)
}
