package drawille

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

func GetMaxFloat64From2dSlice(slices [][]float64) float64 {
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