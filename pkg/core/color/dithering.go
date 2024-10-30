package color

import (
	"math"

	"github.com/watzon/tide/pkg/core/geometry"
)

// DitherMethod represents different dithering algorithms
type DitherMethod int

const (
	DitherNone DitherMethod = iota
	DitherFloydSteinberg
	DitherOrdered
	DitherBayer
)

// DitherMatrix represents a matrix for ordered dithering
type DitherMatrix [][]float64

var (
	// Bayer4x4 is a 4x4 Bayer dithering matrix
	Bayer4x4 = DitherMatrix{
		{0.0, 8.0, 2.0, 10.0},
		{12.0, 4.0, 14.0, 6.0},
		{3.0, 11.0, 1.0, 9.0},
		{15.0, 7.0, 13.0, 5.0},
	}
)

// ErrorBuffer stores error terms for Floyd-Steinberg dithering
type ErrorBuffer struct {
	errors   map[geometry.Point][3]float64
	minPoint geometry.Point
	maxPoint geometry.Point
}

// NewErrorBuffer creates a new error buffer for the given bounds
func NewErrorBuffer(bounds geometry.Rect) *ErrorBuffer {
	return &ErrorBuffer{
		errors:   make(map[geometry.Point][3]float64),
		minPoint: bounds.Min,
		maxPoint: bounds.Max,
	}
}

// Get retrieves the error terms at a position
func (b *ErrorBuffer) Get(p geometry.Point) [3]float64 {
	if err, ok := b.errors[p]; ok {
		return err
	}
	return [3]float64{0, 0, 0}
}

// Set stores error terms at a position
func (b *ErrorBuffer) Set(p geometry.Point, err [3]float64) {
	if p.X >= b.minPoint.X && p.X < b.maxPoint.X &&
		p.Y >= b.minPoint.Y && p.Y < b.maxPoint.Y {
		b.errors[p] = err
	}
}

// Clear removes all stored error terms
func (b *ErrorBuffer) Clear() {
	b.errors = make(map[geometry.Point][3]float64)
}

// Dither applies the specified dithering method to a color
func (c Color) Dither(method DitherMethod, x, y int, palette []Color, buffer ...*ErrorBuffer) Color {
	// Early return for empty palette
	if len(palette) == 0 {
		return c
	}

	switch method {
	case DitherFloydSteinberg:
		var b *ErrorBuffer
		if len(buffer) > 0 {
			b = buffer[0]
		}
		if b == nil {
			// Create default buffer only if none provided
			b = NewErrorBuffer(geometry.Rect{
				Min: geometry.Point{X: 0, Y: 0},
				Max: geometry.Point{X: x + 2, Y: y + 2}, // +2 to accommodate error diffusion
			})
		}
		return c.floydSteinbergDither(x, y, palette, b)
	case DitherOrdered:
		matrix := Bayer4x4
		if len(matrix) == 0 {
			return c.nearestColor(palette)
		}
		return c.orderedDither(x, y, palette, matrix)
	case DitherBayer:
		return c.bayerDither(x, y, palette)
	default:
		return c.nearestColor(palette)
	}
}

// nearestColor finds the closest color in the palette
func (c Color) nearestColor(palette []Color) Color {
	if len(palette) == 0 {
		return c // Return original color for empty palette
	}

	if len(palette) == 1 {
		return palette[0] // Return the only color for single-color palette
	}

	minDist := float64(math.MaxFloat64)
	var nearest Color

	for _, p := range palette {
		dist := ColorDistance(c, p)
		if dist < minDist {
			minDist = dist
			nearest = p
		}
	}

	return nearest
}

// floydSteinbergDither implements Floyd-Steinberg dithering
func (c Color) floydSteinbergDither(x, y int, palette []Color, buffer *ErrorBuffer) Color {
	if buffer == nil {
		return c.nearestColor(palette)
	}

	p := geometry.Point{X: x, Y: y}
	err := buffer.Get(p)

	// Apply stored error
	adjusted := Color{
		R: uint8(math.Max(0, math.Min(255, float64(c.R)+err[0]))),
		G: uint8(math.Max(0, math.Min(255, float64(c.G)+err[1]))),
		B: uint8(math.Max(0, math.Min(255, float64(c.B)+err[2]))),
		A: c.A,
	}

	// Find nearest color in palette
	nearest := adjusted.nearestColor(palette)

	// Calculate new error
	newErr := [3]float64{
		float64(adjusted.R) - float64(nearest.R),
		float64(adjusted.G) - float64(nearest.G),
		float64(adjusted.B) - float64(nearest.B),
	}

	// Distribute error to neighboring pixels
	// Floyd-Steinberg distribution pattern:
	//     X   7/16
	// 3/16  5/16  1/16
	neighbors := []struct {
		offset geometry.Point
		weight float64
	}{
		{geometry.Point{X: 1, Y: 0}, 7.0 / 16.0},
		{geometry.Point{X: -1, Y: 1}, 3.0 / 16.0},
		{geometry.Point{X: 0, Y: 1}, 5.0 / 16.0},
		{geometry.Point{X: 1, Y: 1}, 1.0 / 16.0},
	}

	for _, neighbor := range neighbors {
		neighborPoint := geometry.Point{X: x + neighbor.offset.X, Y: y + neighbor.offset.Y}
		buffer.Set(neighborPoint, [3]float64{
			buffer.Get(neighborPoint)[0] + newErr[0]*neighbor.weight,
			buffer.Get(neighborPoint)[1] + newErr[1]*neighbor.weight,
			buffer.Get(neighborPoint)[2] + newErr[2]*neighbor.weight,
		})
	}

	return nearest
}

// orderedDither implements ordered dithering using a given matrix
func (c Color) orderedDither(x, y int, palette []Color, matrix DitherMatrix) Color {
	if len(matrix) == 0 {
		return c.nearestColor(palette)
	}

	// Get threshold from matrix, wrapping around matrix dimensions
	mx := x % len(matrix)
	my := y % len(matrix[0])

	// Normalize threshold to -0.5 to 0.5 range for proper distribution
	threshold := (matrix[my][mx] / 16.0) - 0.5

	// For 50% gray input with black/white palette, this should produce
	// a checkerboard pattern. Scale factor of 32 helps achieve this.
	adjusted := Color{
		R: uint8(math.Max(0, math.Min(255, float64(c.R)+threshold*32))),
		G: uint8(math.Max(0, math.Min(255, float64(c.G)+threshold*32))),
		B: uint8(math.Max(0, math.Min(255, float64(c.B)+threshold*32))),
		A: c.A,
	}

	return adjusted.nearestColor(palette)
}

// bayerDither implements Bayer matrix dithering
func (c Color) bayerDither(x, y int, palette []Color) Color {
	return c.orderedDither(x, y, palette, Bayer4x4)
}
