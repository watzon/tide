package color

import "math"

// Lerp performs linear interpolation between two colors
func Lerp(c1, c2 Color, t float64) Color {
	t = math.Max(0, math.Min(1, t)) // Clamp t between 0 and 1
	return Color{
		R: uint8(float64(c1.R) + t*float64(c2.R-c1.R)),
		G: uint8(float64(c1.G) + t*float64(c2.G-c1.G)),
		B: uint8(float64(c1.B) + t*float64(c2.B-c1.B)),
		A: uint8(float64(c1.A) + t*float64(c2.A-c1.A)),
	}
}

// Gradient generates a slice of colors interpolating between start and end
func Gradient(start, end Color, steps int) []Color {
	if steps < 2 {
		return []Color{start}
	}

	result := make([]Color, steps)
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		result[i] = Lerp(start, end, t)
	}
	return result
}

// Mix blends two colors with the given weight (0.0 to 1.0)
func Mix(c1, c2 Color, weight float64) Color {
	return Lerp(c1, c2, weight)
}
