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

func (c Cell) String() string {
	return color(string(c.Rune), c.color)
}

// Canvas is a plot of braille characters
type Canvas struct {
	LineColors []AnsiColor
	LabelColor AnsiColor
	AxisColor  AnsiColor

	// a map of the entire braille grid
	area   image.Rectangle
	points map[image.Point]Cell
	labels []string
	offset int
}

// Make a new canvas
func NewCanvas(width, height int) Canvas {
	c := Canvas{
		AxisColor: Default,
		LabelColor: Default,
		LineColors: []AnsiColor{},
		area:   image.Rect(0, 0, width-1, height-1),
		points: make(map[image.Point]Cell),
		labels: []string{},
	}
	return c
}

func (c *Canvas) Plot(data [][]float64) string {
	_, maxDataPoint := GetMinMaxFloat64From2dSlice(data)
	lenMaxDataPoint := len(fmt.Sprintf("%.2f", maxDataPoint))
	c.offset = lenMaxDataPoint + 3
	verticalScale := maxDataPoint / float64(c.area.Dy())
	for i := 0; i < c.area.Dy(); i++ {
		val := fmt.Sprintf("%.2f", float64(i)*verticalScale)
		padding := ""
		if len(val) < lenMaxDataPoint {
			padding = strings.Repeat(" ", lenMaxDataPoint-len(val))
		}
		label := fmt.Sprintf("%s%s %s ", padding, color(val, c.LabelColor), color("┤", c.AxisColor))
		c.labels = append(c.labels, label)
	}
	for i, line := range data {
		previousHeight := int((line[1] / maxDataPoint) * float64(c.area.Dy()-1))
		for j, val := range line[1:] {
			height := int((val / maxDataPoint) * float64(c.area.Dy()-1))
			startX := c.area.Min.X + c.offset
			c.SetLine(
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
	return c.String()
}

func (c Canvas) String() string {
	var b strings.Builder
	// max := 0
	cells := c.GetCells()
	for row := c.area.Dy() - 1; row >= 0; row-- {
		label := c.labels[row]
		b.WriteString(label)
		for col := c.offset; col < c.area.Dx(); col++ {
			r := " "
			if cell, ok := cells[image.Pt(col, row)]; ok {
				r = cell.String()
			}
			b.WriteString(r)
		}
		b.WriteRune('\n')
	}

	// this is so that the x-axis can start at the y-axis line
	offset := c.offset - 2
	b.WriteString(strings.Repeat(" ", offset))
	b.WriteString(color(string('╰'), c.AxisColor))
	b.WriteString(color(strings.Repeat("─", c.area.Dx()-offset), c.AxisColor))
	return b.String()
}

func (c *Canvas) SetPoint(p image.Point, color AnsiColor) {
	point := image.Pt(p.X/2, p.Y/4)
	c.points[point] = Cell{
		c.points[point].Rune | BRAILLE[p.Y%4][p.X%2],
		color,
	}
}

func (c *Canvas) SetLine(p0, p1 image.Point, color AnsiColor) {
	for _, p := range line(p0, p1) {
		c.SetPoint(p, color)
	}
}

func (c *Canvas) GetCells() map[image.Point]Cell {
	points := make(map[image.Point]Cell)
	for point, cell := range c.points {
		points[point] = Cell{cell.Rune + BRAILLE_OFFSET, cell.color}
	}
	return points
}

func line(p0, p1 image.Point) []image.Point {
	points := []image.Point{}

	leftPoint, rightPoint := p0, p1
	if leftPoint.X > rightPoint.X {
		leftPoint, rightPoint = rightPoint, leftPoint
	}

	xDistance := absInt(leftPoint.X - rightPoint.X)
	yDistance := absInt(leftPoint.Y - rightPoint.Y)
	slope := float64(yDistance) / float64(xDistance)
	slopeSign := 1
	if rightPoint.Y < leftPoint.Y {
		slopeSign = -1
	}

	targetYCoordinate := float64(leftPoint.Y)
	currentYCoordinate := leftPoint.Y
	for i := leftPoint.X; i < rightPoint.X; i++ {
		points = append(points, image.Pt(i, currentYCoordinate))
		targetYCoordinate += (slope * float64(slopeSign))
		for currentYCoordinate != int(targetYCoordinate) {
			points = append(points, image.Pt(i, currentYCoordinate))
			currentYCoordinate += slopeSign
		}
	}

	return points
}
