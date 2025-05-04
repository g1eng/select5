package select5

// Terminal control escape sequences
const (
	ClearScreen           = "\x1b[2J"     // Clear entire screen use with print functions
	ClearScreenFromCursor = "\x1b[J"      // Clear all character after the cursor in the current line
	ClearLine             = "\x1b[2K"     // Clear current line with print functions
	ClearLineFromCursor   = "\x1b[K"      // Clear all character after the cursor in the current line
	ResetCursor           = "\x1b[H"      // Move cursor to top-left cursor position
	CursorUp              = "\x1b[1A"     // Move cursor up one line
	CursorDown            = "\x1b[1B"     // Move cursor down one line
	CursorRight           = "\x1b[1C"     // Move cursor right one character
	CursorLeft            = "\x1b[1D"     // Move cursor left one character
	HideCursor            = "\x1b[?25l"   // Hide cursor with print functions
	ShowCursor            = "\x1b[?25h"   // Show cursor with print functions
	MoveTo                = "\x1b[%d;%dH" // Move cursor to position with fmt.Printf

	BS       = 0x08
	ENTER    = 0x0a
	ESC      = 0x1b
	DEL      = 0x7f
	UP       = 0x1b5b41
	DOWN     = 0x1b5b42
	RIGHT    = 0x1b5b43
	LEFT     = 0x1b5b44
	END      = 0x1b5b46
	HOME     = 0x1b5b48
	PAGEUP   = 0x1b357e
	PAGEDOWN = 0x1b367e

	CtrlA = 0x01
	CtrlB = 0x01
	CtrlC = 0x03
	CtrlD = 0x04
	CtrlE = 0x05
	CtrlN = 0x0e
	CtrlP = 0x10
	CtrlX = 0x18
	CtrlY = 0x19
	CtrlZ = 0x1a
)
