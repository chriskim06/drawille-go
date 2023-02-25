package drawille

import "math"

// dots:
//    8x8 cells, 4x2 chars
//    ,___,___,___,___,
//    |1 4|1 4|1 4|1 4|
//    |2 5|2 5|2 5|2 5|
//    |3 6|3 6|3 6|3 6|
//    |7 8|7 8|7 8|7 8|
//    -----------------
//    |1 4|1 4|1 4|1 4|
//    |2 5|2 5|2 5|2 5|
//    |3 6|3 6|3 6|3 6|
//    |7 8|7 8|7 8|7 8|
//    `````````````````
/*
var pixel_map = [4][2]int{
	{0x1, 0x8},
	{0x2, 0x10},
	{0x4, 0x20},
	{0x40, 0x80},
}

// Braille chars start at 0x2800
var braille_char_offset = 0x2800

func getPixel(y, x int) int {
	var cy, cx int
	if y >= 0 {
		cy = y % 4
	} else {
		cy = 3 + ((y + 1) % 4)
	}
	if x >= 0 {
		cx = x % 2
	} else {
		cx = 1 + ((x + 1) % 2)
	}
	return pixel_map[cy][cx]
}
*/

func GetMinMaxFloat64From2dSlice(slices [][]float64) (float64, float64) {
	min, max := math.Inf(1), math.Inf(-1)
	for _, slice := range slices {
		for _, val := range slice {
			if val > max {
				max = val
			}
			if val < min {
				min = val
			}
		}
	}
	return min, max
}

func absInt(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}