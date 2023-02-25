package main

import (
	"fmt"
	"math"

	"github.com/chriskim06/drawille-go"
)

func main() {
	s := drawille.NewCanvas()
	s.LineColors = []drawille.AnsiColor{
		drawille.Red,
		drawille.RoyalBlue,
	}
	// for x := 0; x < (1800); x = x + 1 {
	// 	y := int(math.Sin((math.Pi/180)*float64(x))*10 + 0.5)
	// 	s.Set(x/10, y, 0)
	// }
	// fmt.Print(s)

	// s.Clear()

	for x := 0; x < 1800; x = x + 10 {
		s.Set(0, x/10, int(10+math.Sin((math.Pi/180)*float64(x))*10+0.5))
		s.Set(1, x/10, int(10+math.Cos((math.Pi/180)*float64(x))*10+0.5))
	}
	fmt.Print(s)

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
