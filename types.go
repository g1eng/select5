package select5

// KeyEvent represents a keyboard input event with information about special keys and modifiers.
type KeyEvent struct {
	Key     rune   // The character pressed
	Code    int    // Numeric key code
	Ctrl    bool   // Whether Ctrl was pressed
	Alt     bool   // Whether Alt was pressed
	Shift   bool   // Whether Shift was pressed
	Special string // Special key name (UP, DOWN, ENTER, etc.)
}
