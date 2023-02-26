package drawille

import (
	"fmt"
	"image"
	"math"
	"strings"
)

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

// Canvas is a plot of braille characters
type Canvas struct {
	LineColors []Color
	LabelColor Color
	AxisColor  Color
	ShowAxis   bool

	// the bounds of the canvas
	area image.Rectangle

	// a map of the entire braille grid
	points map[image.Point]Cell

	labels []string
	offset int
}

// Make a new canvas
func NewCanvas(width, height int) Canvas {
	c := Canvas{
		AxisColor:  Default,
		LabelColor: Default,
		LineColors: []Color{},
		ShowAxis:   true,
		area:       image.Rect(0, 0, width, height),
		points:     make(map[image.Point]Cell),
		labels:     []string{},
	}
	return c
}

func (c *Canvas) clear() {
	c.points = make(map[image.Point]Cell)
	c.labels = []string{}
}

// Plot takes a list of data points to graph
// each inner slice represents a different line
func (c *Canvas) Plot(data [][]float64, bounds ...float64) string {
	if len(data) == 0 {
		return ""
	}
	c.clear()
	maxDataPoint := getMaxFloat64From2dSlice(data)
	graphHeight := c.area.Dy()
	xlabels := bounds

	// setup y-axis labels
	if c.ShowAxis {
		verticalScale := maxDataPoint / float64(c.area.Dy())
		lenMaxDataPoint := len(fmt.Sprintf("%.2f", maxDataPoint))
		for i := 0; i < c.area.Dy(); i++ {
			val := fmt.Sprintf("%.2f", float64(i)*verticalScale)
			padStr := ""
			if len(val) < lenMaxDataPoint {
				padStr = padding(lenMaxDataPoint-len(val))
			}
			c.labels = append(c.labels, fmt.Sprintf(
				"%s%s %s ",
				padStr,
				wrap(val, c.LabelColor),
				wrap("┤", c.AxisColor)),
			)
		}
		c.offset = lenMaxDataPoint + 3 // y-axis plus spaces around it
		graphHeight--
		if len(xlabels) == 2 {
			graphHeight--
		}
	}

	// plot the data
	graphWidth := c.area.Dx() - c.offset
	for i, line := range data {
		if len(line) == 0 {
			continue
		} else if len(line) > graphWidth {
			start := len(line)-graphWidth
			line = line[start:]
			if len(xlabels) == 2 && xlabels[0] < float64(start) {
				xlabels[0] += float64(start)
			}
		}
		previousHeight := int((line[0] / maxDataPoint) * float64(graphHeight-1))
		for j, val := range line {
			height := int((val / maxDataPoint) * float64(graphHeight-1))
			c.setLine(
				image.Pt(
					(c.offset+j)*2,
					(graphHeight-previousHeight-1)*4,
				),
				image.Pt(
					(c.offset+j+1)*2,
					(graphHeight-height-1)*4,
				),
				c.lineColor(i),
			)
			previousHeight = height
		}
	}
	return c.string(xlabels)
}

func (c Canvas) string(bounds []float64) string {
	var b strings.Builder
	cells := c.getCells()

	// go through each row of the canvas and print the lines
	for row := 0; row < c.area.Dy(); row++ {
		if c.ShowAxis {
			b.WriteString(c.labels[c.area.Dy()-1-row])
		}
		for col := c.offset; col < c.area.Dx(); col++ {
			b.WriteString(cells[image.Pt(col, row)].String())
		}
		if row < c.area.Dy()-1 {
			b.WriteRune('\n')
		}
	}

	if c.ShowAxis {
		b.WriteRune('\n')

		// start at the y-axis line
		xOffset := c.offset - 2
		b.WriteString(padding(xOffset))
		if len(bounds) != 2 {
			b.WriteString(wrap(fmt.Sprintf("╰%s", strings.Repeat("─", c.area.Dx()-xOffset-1)), c.AxisColor))
			return b.String()
		}

		boundsDistance := math.Abs(bounds[1] - bounds[0])
		labelWidth := len(fmt.Sprintf("%.2f", boundsDistance)) + 2
		graphWidth := c.area.Dx() - c.offset
		numLabels := graphWidth / labelWidth
		remaining := graphWidth % labelWidth
		scale := (boundsDistance - float64(remaining)) / float64(numLabels)
		xaxis := fmt.Sprintf(
			"╰%s%s─",
			strings.Repeat(fmt.Sprintf("┬%s", strings.Repeat("─", labelWidth-1)), numLabels),
			strings.Repeat("─", remaining),
		)
		b.WriteString(wrap(xaxis, c.AxisColor))
		var labelStr strings.Builder
		labelStr.WriteString(padding(c.offset-1))
		for i := 0; i < numLabels; i++ {
			l := fmt.Sprintf("%.2f", scale*float64(i) + bounds[0])
			if len(l) < labelWidth {
				l += padding(labelWidth-len(l))
			}
			labelStr.WriteString(l)
		}
		b.WriteRune('\n')
		b.WriteString(labelStr.String())
	}
	return b.String()
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
