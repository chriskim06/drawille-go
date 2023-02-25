package drawille

import (
	"math"
	"strings"
)

// Clear all pixels
func (c *Canvas) Clear() {
	c.chars = make(map[int]map[int]Cell)
}

func (c Canvas) MaxY() int {
	max := 0
	for k := range c.chars {
		if k > max {
			max = k
		}
	}
	return max * 4
}

func (c Canvas) MinY() int {
	min := 0
	for k := range c.chars {
		if k < min {
			min = k
		}
	}
	return min * 4
}

func (c Canvas) MaxX() int {
	max := 0
	for _, v := range c.chars {
		for k := range v {
			if k > max {
				max = k
			}
		}
	}
	return max * 2
}

func (c Canvas) MinX() int {
	min := 0
	for _, v := range c.chars {
		for k := range v {
			if k < min {
				min = k
			}
		}
	}
	return min * 2
}

// Set a pixel of c
func (c *Canvas) Set(lineNum, x, y int) {
	px, py := getPos(x, y)
	if m := c.chars[py]; m == nil {
		c.chars[py] = make(map[int]Cell)
	}
	color := Default
	if lineNum >= 0 && lineNum < len(c.LineColors) {
		color = c.LineColors[lineNum]
	}
	cell := c.chars[py][px]
	c.chars[py][px] = Cell{
		val:   cell.val | getPixel(y, x),
		color: color,
	}
}

// Unset a pixel of c
func (c *Canvas) Unset(x, y int) {
	px, py := getPos(x, y)
	x, y = int(math.Abs(float64(x))), int(math.Abs(float64(y)))
	if m := c.chars[py]; m == nil {
		c.chars[py] = make(map[int]Cell)
	}
	dot := c.chars[py][px]
	dot.val = dot.val &^ getPixel(y, x)
	c.chars[py][px] = dot
}

// Toggle a point
func (c *Canvas) Toggle(lineNum, x, y int) {
	px, py := getPos(x, y)
	if (c.chars[py][px].val & getPixel(y, x)) != 0 {
		c.Unset(x, y)
	} else {
		c.Set(lineNum, x, y)
	}
}

// Set text to the given coordinates
func (c *Canvas) SetText(x, y int, text string) {
	px, py := x/2, y/4
	if m := c.chars[py]; m == nil {
		c.chars[y] = make(map[int]Cell)
	}
	for i, char := range text {
		c.chars[py][px+i] = Cell{
			val:   int(char) - braille_char_offset,
			color: Default,
		}
	}
}

// Retrieve the rows from a given view
func (c Canvas) Rows(minX, minY, maxX, maxY int) []string {
	minrow, maxrow := minY/4, (maxY)/4
	mincol, maxcol := minX/2, (maxX)/2

	s := make([]string, 0)
	for rownum := minrow; rownum < (maxrow + 1); rownum = rownum + 1 {
		color := Default
		var b strings.Builder
		for x := mincol; x < (maxcol + 1); x = x + 1 {
			dot := c.chars[rownum][x]
			if dot.color != color {
				color = dot.color
				b.WriteString(dot.color.String())
			}
			b.WriteString(string(rune(dot.val + braille_char_offset)))
			if color != Default {
				b.WriteString(Default.String())
			}
		}
		s = append(s, b.String())
	}
	return s
}

// Retrieve a string representation of the frame at the given parameters
func (c Canvas) Frame(minX, minY, maxX, maxY int) string {
	var b strings.Builder
	for _, row := range c.Rows(minX, minY, maxX, maxY) {
		b.WriteString(row)
		b.WriteString(c.LineEnding)
	}
	return b.String()
}

func (c Canvas) String() string {
	// need to be able to deal with setting a fixed canvas height/width
	return c.Frame(c.minX, c.minY, c.maxX, c.maxY)
}

func (c *Canvas) DrawLine(lineNum int, x1, y1, x2, y2 float64) {
	xdiff := math.Abs(x1 - x2)
	ydiff := math.Abs(y2 - y1)

	var xdir, ydir float64
	if x1 <= x2 {
		xdir = 1
	} else {
		xdir = -1
	}
	if y1 <= y2 {
		ydir = 1
	} else {
		ydir = -1
	}

	r := math.Max(xdiff, ydiff)
	for i := 0; i < round(r)+1; i++ {
		x, y := x1, y1
		if ydiff != 0 {
			y += (float64(i) * ydiff) / (r * ydir)
		}
		if xdiff != 0 {
			x += (float64(i) * xdiff) / (r * xdir)
		}
		c.Toggle(lineNum, round(x), round(y))
	}
}

// Plot takes a 2d array of data points, with each inner
// array being a separate line to graph on the Canvas "┤"
func (c *Canvas) Plot(data [][]float64) string {
	// need to get the largest number in the data set to calculate
	// the offset needed for the y-axis
	// need to correctly print y-axis labels based on height and data
	//   for i := 0; i < c.height
	//     label val should be i * (abs(ymax-ymin) / c.height) + ymin
	//       - so if our data has max/min values of 30/10 and canvas
	//         height is set to 10 then the label values going up
	//         would be 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30
	//       - y-axis labels should also include " ┤ " so that when
	//         they stack on top of each other it creates the axis
	maxDataPoint := GetMaxFloat64From2dSlice(data)
	return c.String()
}

// Convert x, y to cols, rows
func getPos(x, y int) (int, int) {
	return (x / 2), (y / 4)
}

func round(x float64) int {
	return int(x + 0.5)
}
