package main

import (
	"fmt"

	"github.com/chriskim06/drawille-go"
)

func main() {
	s := drawille.NewCanvas(100, 25)
	s.LineColors = []drawille.Color{
		drawille.Red,
		drawille.RoyalBlue,
	}
	s.AxisColor = drawille.SeaGreen

	data := [][]float64{{}, {}}
	for x := 0; x < 16; x++ {
		data[0] = append(data[0], 150)
		data[1] = append(data[1], 25)
	}
	for x := 0; x < 20; x++ {
		data[0] = append(data[0], 256)
		data[1] = append(data[1], 40)
	}
	for x := 0; x < 16; x++ {
		data[0] = append(data[0], 140)
		data[1] = append(data[1], 17)
	}
	for x := 0; x < 22; x++ {
		data[0] = append(data[0], 256)
		data[1] = append(data[1], 17)
	}
	for x := 0; x < 22; x++ {
		data[0] = append(data[0], 355)
		data[1] = append(data[1], 17)
	}
	fmt.Println(s.Plot(data, 0, 90))
	// fmt.Println(s.Plot(data))
	// for x := 0; x < 26; x++ {
	// 	data[0] = append(data[0], 355)
	// 	data[1] = append(data[1], 17)
	// }
	// fmt.Println(s.Plot(data, 0, 91))
}
