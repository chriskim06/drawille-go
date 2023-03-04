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
	HorizontalLabels []string
	verticalLabels   []string

	// this value is used to determine the horizontal scale when
	// graphing points in the plot. with braille characters, each
	// cell represents 2 points along the x axis. so if the total
	// graphable area in the canvas is 50 cells then we can plot
	// 100 points of data. if we only want to plot 50 data points
	// then the horizontal scale needs to be 2 in order to utilize
	// the whole graphable area. since image uses ints for points
	// the graphing gets a little weird. ex:
	// graphable area = 88
	// num of data points = 50
	// horizontal scale = 88/50 = 1.76
	// 0, 1.76, 3.52, 5.28, 7.04, 8.8, 10.56, 12.32, 14.08
	// this would map points at 0,1,3,5,7,8,10,12,14
	NumDataPoints int

	// the bounds of the canvas
	area image.Rectangle

	graphWidth, graphHeight int

	// a map of the entire braille grid
	points map[image.Point]Cell

	horizontalOffset int
	verticalOffset   int
	horizontalScale  float64
	verticalScale    float64
	maxX             int
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
	c.graphHeight = c.area.Dy()
	maxDataPoint := getMaxFloat64From2dSlice(data)
	dataLen := 0
	for _, line := range data {
		if len(line) > dataLen {
			dataLen = len(line)
		}
	}

	// setup y-axis labels
	if c.ShowAxis {
		c.graphHeight--
		lenMaxDataPoint := len(fmt.Sprintf("%.2f", maxDataPoint))
		c.horizontalOffset = lenMaxDataPoint + 2 // y-axis plus space before it
		if len(c.HorizontalLabels) != 0 && len(c.HorizontalLabels) <= c.area.Dx()-c.horizontalOffset {
			c.graphHeight--
		}
		verticalScale := maxDataPoint / float64(c.graphHeight-1)
		for i := c.graphHeight - 1; i >= 0; i-- {
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
	c.graphWidth = (c.area.Dx() - c.horizontalOffset) * 2

	// plot the data
	c.horizontalScale = 1.0
	if c.NumDataPoints > 0 {
		c.horizontalScale = float64(c.graphWidth) / float64(c.NumDataPoints)
	}
	for i, line := range data {
		if len(line) == 0 {
			continue
		} else if c.NumDataPoints > 0 && len(line) > c.NumDataPoints {
			start := len(line) - c.NumDataPoints
			line = line[start:]
		} else if len(line) > c.graphWidth {
			start := len(line) - c.graphWidth
			line = line[start:]
		}
		previousHeight := int((line[0] / maxDataPoint) * float64(c.graphHeight-1))
		length := 0
		for j, val := range line {
			height := int((val / maxDataPoint) * float64(c.graphHeight-1))
			x0 := int(float64(j) * c.horizontalScale)
			x1 := int(float64(j+1) * c.horizontalScale)
			c.setLine(
				image.Pt(
					//                     (c.horizontalOffset+int(float64(j)*c.horizontalScale)),
					c.horizontalOffset+x0,
					(c.graphHeight-previousHeight-1)*4,
				),
				image.Pt(
					//                     (c.horizontalOffset+int(float64(j+1)*c.horizontalScale)),
					c.horizontalOffset+x1,
					(c.graphHeight-height-1)*4,
				),
				c.lineColor(i),
			)
			length += x1 - x0
			previousHeight = height
			//             if int(float64(j+1)*c.horizontalScale) > c.maxX {
			//                 c.maxX = int(float64(j+1)*c.horizontalScale) / 2
			//             }
		}
		if length > c.maxX {
			c.maxX = length
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
			b.WriteString(c.verticalLabels[row])
		}
		for col := c.horizontalOffset; col < c.area.Dx(); col++ {
			b.WriteString(cells[image.Pt(col, row)].String())
		}
		if row < c.graphHeight-1 {
			b.WriteString("\n")
		}
	}

	if c.ShowAxis {
		b.WriteString("\n")
		//         b.WriteString(padding(c.horizontalOffset - 1))
		//         b.WriteString(wrap(fmt.Sprintf("╰%s", strings.Repeat("─", c.graphWidth)), c.AxisColor))

		// no labels for the x-axis so just draw a line
		// or caller didnt properly update the x-axis labels
		if len(c.HorizontalLabels) == 0 || len(c.HorizontalLabels) > c.graphWidth {
			b.WriteString(padding(c.horizontalOffset - 1))
			b.WriteString(wrap(fmt.Sprintf("╰%s", strings.Repeat("─", c.graphWidth)), c.AxisColor))
			return b.String()
		}

		var axisStr, labelStr strings.Builder
		axisStr.WriteString(fmt.Sprintf("%s╰", padding(c.horizontalOffset-1)))
		labelStr.WriteString(padding(c.horizontalOffset)) // y-axis line plus the padding

		// need to put proper labels on the x axis based on how many points have been
		// plotted and the horizontal scale. if we have a graph thats 103 in width
		// then there are 206 plottable points. if the max number we're plotting is 50
		// then the horizontal scale is 206/50 = 4.12. that means we plot points at
		// 0, 4, 8, 12, 16, 20, 24, 28, 32, 37, ...
		// label 0 is '└12:34:56  ', 11 chars wide so next label can be at x coordinate
		// 11. this is x=22 in terms of braille dots or graph width. we have labels
		// for points plotted at 20 and 24 so we should use the label corresponding to 20
		// or the 5th item in the labels array
		//
		// seems ok in topui so far. need to figure out why its not plotting full width
		// also doesnt immediately show data and some pods mem seems behind by 1
		xCoordinate := 0
		pos := 0
		remaining := c.graphWidth / 2
		for remaining > 0 {
			labelToAdd := c.HorizontalLabels[pos]
			xCoordinate += len(labelToAdd) + 3
			if len(labelToAdd)+1 > remaining || xCoordinate > c.maxX/2 {
				axisStr.WriteString(strings.Repeat("─", remaining))
				break
			}
			labelStr.WriteString(wrap("└", c.AxisColor) + wrap(labelToAdd, c.LabelColor))
			axisStr.WriteString("┬" + strings.Repeat("─", len(labelToAdd)))
			remaining -= len(labelToAdd) + 1
			if remaining < 2 {
				axisStr.WriteString(strings.Repeat("─", remaining))
				break
			}
			labelStr.WriteString("  ")
			axisStr.WriteString("──")
			remaining -= 2
			//             f := float64((len(labelToAdd)+3)*2) / c.horizontalScale
			pos += ((len(labelToAdd) + 3) * 2) / int(c.horizontalScale)
			//             if i := int(float64(pos) + f + 0.5); i < len(c.HorizontalLabels) {
			//                 pos = i
			//             } else {
			//                 pos += int(f)
			//             }
			//             pos += int(float64(len(labelToAdd)+3) / c.horizontalScale)
			if pos >= len(c.HorizontalLabels) {
				axisStr.WriteString(strings.Repeat("─", remaining))
				break
			}
		}

		b.WriteString(wrap(axisStr.String(), c.AxisColor) + "\n")
		b.WriteString(labelStr.String())
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
