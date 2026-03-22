package ansi

const (
	ArrowKeyUp     = "\033[A"
	ArrowKeyDown   = "\033[B"
	PageUpKey      = "\033[5~"
	PageDownKey    = "\033[6~"
	ClearScreen    = "\033[2J"
	CursorHome     = "\033[H"
	CursorToPos    = "\033[%d;%dH" // use with fmt.Sprintf, the two ds are for the row and column coordinates
	EnterAltScreen = "\033[?1049h"
	ExitAltScreen  = "\033[?1049l"
	Red            = 31
	Green          = 32
	BrightGreen    = 92
	Yellow         = 33
	Blue           = 34
	BrightBlue     = 94
	Magenta        = 35
	Cyan           = 36
	Bold           = 1
	DefaultColor   = 39
	Reset          = 0
)
