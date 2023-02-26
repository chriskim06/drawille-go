package drawille

import (
	"fmt"
	"image"
	"strings"
)

func getMaxFloat64From2dSlice(slices [][]float64) float64 {
	var max float64
	for _, slice := range slices {
		for _, val := range slice {
			if val > max {
				max = val
			}
		}
	}
	return max
}

func padding(i int) string {
	return strings.Repeat(" ", i)
}

func wrap(s string, c Color) string {
	if c == Default {
		return s
	}
	return fmt.Sprintf("%s%s%s", c, s, reset)
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

func absInt(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}
