package drawille

type Cell struct {
	val   int
	color AnsiColor
}

type Canvas struct {
	LineEnding string
	LineColors []AnsiColor

	width, height int
	minX, minY    int
	maxX, maxY    int

	// a map of the entire braille grid
	chars map[int]map[int]Cell
}

// Make a new canvas
func NewCanvas(width, height int) Canvas {
	c := Canvas{
		LineEnding: "\n",
		chars:      make(map[int]map[int]Cell),
		width:      width,
		height:     height,
		minX:       0,
		minY:       0,
	}
	return c
}
