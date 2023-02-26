package drawille

const BRAILLE_OFFSET = '\u2800'

var BRAILLE = [4][2]rune{
	{'\u0001', '\u0008'},
	{'\u0002', '\u0010'},
	{'\u0004', '\u0020'},
	{'\u0040', '\u0080'},
}

// Cell represents the braille character at some coordinate in the canvas
type Cell struct {
	Rune  rune
	color Color
}

// String returns the cell's rune wrapped in the color escape strings
func (c Cell) String() string {
	if c.Rune == 0 {
		return wrap(" ", c.color)
	}
	return wrap(string(c.Rune), c.color)
}
