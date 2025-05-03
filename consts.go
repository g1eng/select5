package select5

// Terminal control escape sequences
const (
	ClearScreen = "\033[2J"     // Clear entire screen
	ClearLine   = "\033[2K"     // Clear current line
	ResetCursor = "\033[H"      // Move cursor to top-left corner
	CursorUp    = "\033[1A"     // Move cursor up one line
	CursorDown  = "\033[1B"     // Move cursor down one line
	HideCursor  = "\033[?25l"   // Hide cursor
	ShowCursor  = "\033[?25h"   // Show cursor
	MoveTo      = "\033[%d;%dH" // Move cursor to position
)
