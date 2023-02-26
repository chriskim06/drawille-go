package main

import (
	"fmt"

	"github.com/chriskim06/drawille-go"
)

func main() {
	s := drawille.NewCanvas(50, 25)
	s.LineColors = []drawille.Color{
		drawille.Red,
		drawille.RoyalBlue,
	}
	s.AxisColor = drawille.SeaGreen

	// graph width here is 50-9=41
	data := [][]float64{{}, {}}
	fmt.Println(s.Plot(data))
	for x := 0; x < 8; x++ {
		data[0] = append(data[0], 150)
		data[1] = append(data[1], 25)
	}
	for x := 0; x < 10; x++ {
		data[0] = append(data[0], 256)
		data[1] = append(data[1], 40)
	}
	for x := 0; x < 8; x++ {
		data[0] = append(data[0], 140)
		data[1] = append(data[1], 17)
	}
	for x := 0; x < 11; x++ {
		data[0] = append(data[0], 140)
		data[1] = append(data[1], 17)
	}
	fmt.Println(s.Plot(data))
	for x := 0; x < 4; x++ {
		data[0] = append(data[0], 140)
		data[1] = append(data[1], 17)
	}
	fmt.Println(s.Plot(data))
	// for x := 0; x < 3; x++ {
	// 	data[0] = append(data[0], 140)
	// 	data[1] = append(data[1], 17)
	// }
	// fmt.Println(s.Plot(data))
	// for x := 0; x < 25; x++ {
	// 	data[0] = append(data[0], 256)
	// 	data[1] = append(data[1], 30)
	// }
	// fmt.Println(s.Plot(data))
	// s.ShowAxis = false
	// data = [][]float64{{}, {}}
	// for x := 0; x < 2800; x = x + 10 {
	// 	data[0] = append(data[0], (10+math.Sin((math.Pi/180)*float64(x))*10+0.5))
	// 	data[1] = append(data[1], (10+math.Cos((math.Pi/180)*float64(x))*10+0.5))
	// }
	// fmt.Println(s.Plot(data))

	// for x := 0; x < 400; x = x + 10 {
	// 	data[0] = append(data[0], (17))
	// 	data[1] = append(data[1], (15))
	// }
	// fmt.Println(s.Plot(data))

	// s.Clear()

	// for x := 0; x < 3600; x = x + 20 {
	// 	s.Set(1, x/20, int(4+math.Sin((math.Pi/180)*float64(x))*4))
	// }
	// fmt.Print(s)

	// s.Clear()

	// for x := 0; x < 360; x = x + 1 {
	// 	s.Set(0, x/4, int(30+math.Sin((math.Pi/180)*float64(x))*30))
	// }

	// for x := 0; x < 30; x = x + 1 {
	// 	for y := 0; y < 30; y = y + 1 {
	// 		s.Set(1, x, y)
	// 		s.Toggle(1, x+30, y+30)
	// 		s.Toggle(1, x+60, y)
	// 	}
	// }
	// fmt.Print(s)
}
