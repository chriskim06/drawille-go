package drawille

import (
	"fmt"
	"image"
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
	color AnsiColor
}

// String returns the cell's rune wrapped in the color escape strings
func (c Cell) String() string {
	if c.Rune == 0 {
		return color(" ", c.color)
	}
	return color(string(c.Rune), c.color)
}

// Canvas is a plot of braille characters
type Canvas struct {
	LineColors []AnsiColor
	LabelColor AnsiColor
	AxisColor  AnsiColor

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
		LineColors: []AnsiColor{},
		area:       image.Rect(0, 0, width-1, height-1),
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
func (c *Canvas) Plot(data [][]float64) string {
	c.clear()
	maxDataPoint := getMaxFloat64From2dSlice(data)
	lenMaxDataPoint := len(fmt.Sprintf("%.2f", maxDataPoint))
	c.offset = lenMaxDataPoint + 3
	verticalScale := maxDataPoint / float64(c.area.Dy())

	// setup y-axis labels
	for i := 0; i < c.area.Dy(); i++ {
		val := fmt.Sprintf("%.2f", float64(i)*verticalScale)
		padding := ""
		if len(val) < lenMaxDataPoint {
			padding = strings.Repeat(" ", lenMaxDataPoint-len(val))
		}
		label := fmt.Sprintf("%s%s %s ", padding, color(val, c.LabelColor), color("┤", c.AxisColor))
		c.labels = append(c.labels, label)
	}

	// plot the data
	for i, line := range data {
		previousHeight := int((line[1] / maxDataPoint) * float64(c.area.Dy()-1))
		for j, val := range line[1:] {
			height := int((val / maxDataPoint) * float64(c.area.Dy()-1))
			startX := c.area.Min.X + c.offset
			c.setLine(
				image.Pt(
					(startX+j)*2,
					(c.area.Max.Y-previousHeight-1)*4,
				),
				image.Pt(
					(startX+j+1)*2,
					(c.area.Max.Y-height-1)*4,
				),
				c.LineColors[i],
			)
			previousHeight = height
		}
	}
	return c.string()
}

func (c Canvas) string() string {
	var b strings.Builder
	cells := c.getCells()
	for row := c.area.Dy() - 1; row >= 0; row-- {
		label := c.labels[row]
		b.WriteString(label)
		for col := c.offset; col < c.area.Dx(); col++ {
			b.WriteString(cells[image.Pt(col, row)].String())
		}
		b.WriteRune('\n')
	}

	// this is so that the x-axis can start at the y-axis line
	offset := c.offset - 2
	xaxis := fmt.Sprintf(
		"%s%s%s",
		strings.Repeat(" ", offset),
		color(string('╰'), c.AxisColor),
		color(strings.Repeat("─", c.area.Dx()-offset), c.AxisColor),
	)
	b.WriteString(xaxis)
	return b.String()
}

func (c *Canvas) setPoint(p image.Point, color AnsiColor) {
	point := image.Pt(p.X/2, p.Y/4)
	c.points[point] = Cell{
		c.points[point].Rune | BRAILLE[p.Y%4][p.X%2],
		color,
	}
}

func (c *Canvas) setLine(p0, p1 image.Point, color AnsiColor) {
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
