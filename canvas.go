package drawille

// Cell represents the braille character at some coordinate in the canvas
type Cell struct {
	val   int
	color AnsiColor
}

// Canvas is a plot of braille characters
type Canvas struct {
	LineEnding string
	LineColors []AnsiColor

	width, height int
	minX, minY    int
	maxX, maxY    int

	// a map of the entire braille grid
	// the map is map[row][col] = cell
	chars map[int]map[int]Cell
}

// Make a new canvas
func NewCanvas(width, height int) Canvas {
	c := Canvas{
		LineEnding: "\n",
		chars:      make(map[int]map[int]Cell),
		width:      width,
		height:     height,
	}
	return c
}
