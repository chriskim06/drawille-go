package main

import (
	"fmt"
	"time"

	"github.com/chriskim06/drawille-go"
)

var t = time.Now()

func main() {
	s := drawille.NewCanvas(100, 25)
	s.LineColors = []drawille.Color{
		drawille.Red,
		drawille.RoyalBlue,
	}
	s.LabelColor = drawille.Purple
	s.AxisColor = drawille.SeaGreen

	i := 0
	labels := []string{}
	data := [][]float64{{}, {}}
	for x := 0; x < 16; x++ {
		data[0] = append(data[0], 150)
		data[1] = append(data[1], 25)
		update(i, &labels)
		i++
	}
	for x := 0; x < 20; x++ {
		data[0] = append(data[0], 256)
		data[1] = append(data[1], 40)
		update(i, &labels)
		i++
	}
	for x := 0; x < 64; x++ {
		data[0] = append(data[0], 140)
		data[1] = append(data[1], 17)
		update(i, &labels)
		i++
	}
	//     for x := 0; x < 22; x++ {
	//         data[0] = append(data[0], 256)
	//         data[1] = append(data[1], 17)
	//         update(i, &labels)
	//         i++
	//     }
	//     for x := 0; x < 18; x++ {
	//         data[0] = append(data[0], 355)
	//         data[1] = append(data[1], 17)
	//         update(i, &labels)
	//         i++
	//     }
	s.NumDataPoints = 100
	s.HorizontalLabels = labels
	fmt.Print(s.Plot(data))
}

func update(i int, labels *[]string) {
	ti := t.Add(time.Duration(i) * time.Second)
	t = ti
	*labels = append(*labels, fmt.Sprintf("%02d:%02d:%02d", ti.Hour(), ti.Minute(), ti.Second()))
}
