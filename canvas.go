package drawille

import (
	"fmt"
	"image"
	"strings"
)

// Canvas is a plot of braille characters
type Canvas struct {
	// various settings accessible outside this object
	LineColors []Color
	LabelColor Color
	AxisColor  Color
	ShowAxis   bool

	// a list of labels the canvas will print for the x and y axis
	// horizontal labels must be provided by the caller. too lazy
	// to come up with a good way to print stuff so offloading some
	// of that work to the user. when the horizontal labels arent
	// provided an empty line is printed
	//     HorizontalLabels []string
	verticalLabels []string

	// the bounds of the canvas
	area image.Rectangle

	graphWidth, graphHeight int

	// a map of the entire braille grid
	points map[image.Point]Cell

	horizontalOffset int
	verticalOffset   int
}

// NewCanvas creates a default canvas
func NewCanvas(width, height int) Canvas {
	c := Canvas{
		AxisColor:      Default,
		LabelColor:     Default,
		LineColors:     []Color{},
		ShowAxis:       true,
		area:           image.Rect(0, 0, width, height),
		points:         make(map[image.Point]Cell),
		verticalLabels: []string{},
	}
	return c
}

// Plot sets the Canvas and return the string representation of it
func (c *Canvas) Plot(data [][]float64) string {
	if len(data) == 0 {
		return ""
	}
	c.clear()
	maxDataPoint := getMaxFloat64From2dSlice(data)
	graphHeight := c.area.Dy()

	// setup y-axis labels
	if c.ShowAxis {
		graphHeight--
		verticalScale := maxDataPoint / float64(graphHeight-1)
		lenMaxDataPoint := len(fmt.Sprintf("%.2f", maxDataPoint))
		c.horizontalOffset = lenMaxDataPoint + 2 // y-axis plus space before it
		//         if len(c.HorizontalLabels) != 0 && len(c.HorizontalLabels) <= c.area.Dx()-c.horizontalOffset {
		//             graphHeight--
		//         }
		for i := 0; i < graphHeight; i++ {
			val := fmt.Sprintf("%.2f", float64(i)*verticalScale)
			padStr := ""
			if len(val) < lenMaxDataPoint {
				padStr = padding(lenMaxDataPoint - len(val))
			}
			c.verticalLabels = append(c.verticalLabels, fmt.Sprintf(
				"%s%s %s",
				padStr,
				wrap(val, c.LabelColor),
				wrap("┤", c.AxisColor)),
			)
		}
	}
	c.graphWidth = c.area.Dx() - c.horizontalOffset
	c.graphHeight = graphHeight
	c.verticalOffset = graphHeight

	// plot the data
	for i, line := range data {
		if len(line) == 0 {
			continue
		} else if len(line) > c.graphWidth {
			start := len(line) - c.graphWidth
			line = line[start:]
		}
		previousHeight := int((line[0] / maxDataPoint) * float64(c.graphHeight-1))
		for j, val := range line {
			height := int((val / maxDataPoint) * float64(c.graphHeight-1))
			c.setLine(
				image.Pt(
					(c.horizontalOffset+j)*2,
					(c.graphHeight-previousHeight-1)*4,
				),
				image.Pt(
					(c.horizontalOffset+j+1)*2,
					(c.graphHeight-height-1)*4,
				),
				c.lineColor(i),
			)
			previousHeight = height
		}
	}
	return c.String()
}

// String allows the Canvas to implement the Stringer interface
func (c Canvas) String() string {
	var b strings.Builder
	cells := c.getCells()

	// go through each row of the canvas and print the lines
	for row := 0; row < c.graphHeight; row++ {
		if c.ShowAxis {
			if idx := len(c.verticalLabels) - 1 - row; idx >= 0 {
				b.WriteString(wrap(c.verticalLabels[idx], c.LabelColor))
			}
		}
		for col := c.horizontalOffset; col < c.area.Dx(); col++ {
			b.WriteString(cells[image.Pt(col, row)].String())
		}
		if row < c.graphHeight-1 {
			b.WriteRune('\n')
		}
	}

	if c.ShowAxis {
		b.WriteRune('\n')
		b.WriteString(padding(c.horizontalOffset - 1))
		b.WriteString(wrap(fmt.Sprintf("╰%s", strings.Repeat("─", c.graphWidth)), c.AxisColor))

		// no labels for the x-axis so just draw a line
		// or caller didnt properly update the x-axis labels
		//         if len(c.HorizontalLabels) == 0 || len(c.HorizontalLabels) > c.graphWidth {
		//             b.WriteString(wrap(fmt.Sprintf("╰%s", strings.Repeat("─", c.graphWidth)), c.AxisColor))
		//             return b.String()
		//         }

		//         var axisStr, labelStr strings.Builder
		//         axisStr.WriteString("╰─")
		//         labelStr.WriteString(padding(c.horizontalOffset + 1)) // y-axis line plus the padding
		//         pos := 0
		//         remaining := c.graphWidth
		//         for remaining > 0 {
		//             labelToAdd := c.HorizontalLabels[pos]
		//             if len(labelToAdd)+1 > remaining {
		//                 axisStr.WriteString(strings.Repeat("─", remaining))
		//                 break
		//             }
		//             labelStr.WriteString(wrap("└", c.AxisColor) + wrap(labelToAdd, c.LabelColor))
		//             axisStr.WriteString("┬" + strings.Repeat("─", len(labelToAdd)))
		//             remaining -= len(labelToAdd) + 1
		//             if remaining < 2 {
		//                 axisStr.WriteString(strings.Repeat("─", remaining))
		//                 break
		//             }
		//             labelStr.WriteString("  ")
		//             remaining -= 2
		//             pos += len(labelToAdd) + 3
		//             if pos >= len(c.HorizontalLabels) {
		//                 axisStr.WriteString(strings.Repeat("─", remaining))
		//                 break
		//             }
		//             axisStr.WriteString("──")
		//         }
		//
		//         b.WriteString(wrap(axisStr.String(), c.AxisColor) + "\n")
		//         b.WriteString(labelStr.String())
	}
	return b.String()
}

// SetSize changes the size of the Canvas dimensions
func (c *Canvas) SetSize(width, height int) {
	c.area = image.Rect(0, 0, width, height)
}

// GetSize returns the total canvas width and height, including
// labels, axes, padding, etc.
func (c Canvas) GetSize() (int, int) {
	return c.area.Dx(), c.area.Dy()
}

// GetPlotSize returns the width and height of the area of the
// rectangle that can contain plot data points
func (c Canvas) GetPlotSize() (int, int) {
	width, height := c.GetSize()
	return width - c.horizontalOffset, height - c.verticalOffset
}

func (c *Canvas) clear() {
	c.points = make(map[image.Point]Cell)
	c.verticalLabels = []string{}
}

func (c *Canvas) setPoint(p image.Point, color Color) {
	point := image.Pt(p.X/2, p.Y/4)
	c.points[point] = Cell{
		c.points[point].Rune | BRAILLE[p.Y%4][p.X%2],
		color,
	}
}

func (c *Canvas) setLine(p0, p1 image.Point, color Color) {
	for _, p := range line(p0, p1) {
		c.setPoint(p, color)
	}
}

func (c *Canvas) getCells() map[image.Point]Cell {
	points := make(map[image.Point]Cell)
	for point, cell := range c.points {
		points[point] = Cell{cell.Rune + BRAILLE_OFFSET, cell.color}
	}
	return points
}

func (c Canvas) lineColor(i int) Color {
	if len(c.LineColors) == 0 || i > len(c.LineColors)-1 {
		return Default
	}
	return c.LineColors[i]
}
