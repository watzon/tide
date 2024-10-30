package color_test

import (
	"fmt"
	"testing"

	"github.com/watzon/tide/pkg/core/color"
	"github.com/watzon/tide/pkg/core/geometry"
)

func TestDitherMethods(t *testing.T) {
	// Common test palette
	palette := []color.Color{
		{R: 0, G: 0, B: 0, A: 255},       // Black
		{R: 255, G: 255, B: 255, A: 255}, // White
	}

	t.Run("DitherNone", func(t *testing.T) {
		c := color.Color{R: 128, G: 128, B: 128, A: 255}
		result := c.Dither(color.DitherNone, 0, 0, palette)

		// Should map to either black or white, whichever is closer
		if result != palette[0] && result != palette[1] {
			t.Errorf("DitherNone should map to nearest palette color")
		}
	})

	t.Run("Empty palette", func(t *testing.T) {
		c := color.Color{R: 128, G: 128, B: 128, A: 255}
		result := c.Dither(color.DitherNone, 0, 0, nil)

		// Should return original color when palette is empty
		if result != c {
			t.Errorf("Expected original color with empty palette, got %v", result)
		}
	})

	t.Run("Bayer dithering", func(t *testing.T) {
		c := color.Color{R: 128, G: 128, B: 128, A: 255}

		// Test 4x4 pattern
		results := make([]color.Color, 16)
		blackCount := 0
		whiteCount := 0

		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				result := c.Dither(color.DitherBayer, x, y, palette)
				results[y*4+x] = result

				if result == palette[0] {
					blackCount++
				} else if result == palette[1] {
					whiteCount++
				}
			}
		}

		// For 50% gray, should be roughly equal black and white
		if blackCount < 6 || whiteCount < 6 {
			t.Errorf("Bayer dithering should produce roughly equal black and white pixels, got %d black, %d white",
				blackCount, whiteCount)
		}
	})
}

func TestErrorBuffer(t *testing.T) {
	bounds := geometry.Rect{
		Min: geometry.Point{X: 0, Y: 0},
		Max: geometry.Point{X: 10, Y: 10},
	}

	t.Run("NewErrorBuffer", func(t *testing.T) {
		buffer := color.NewErrorBuffer(bounds)
		if buffer == nil {
			t.Error("NewErrorBuffer should not return nil")
		}
	})

	t.Run("Set and Get", func(t *testing.T) {
		buffer := color.NewErrorBuffer(bounds)
		p := geometry.Point{X: 5, Y: 5}
		err := [3]float64{1.0, 2.0, 3.0}

		buffer.Set(p, err)
		got := buffer.Get(p)

		if got != err {
			t.Errorf("Get returned %v, want %v", got, err)
		}
	})

	t.Run("Get nonexistent point", func(t *testing.T) {
		buffer := color.NewErrorBuffer(bounds)
		p := geometry.Point{X: 5, Y: 5}
		got := buffer.Get(p)

		if got != [3]float64{0, 0, 0} {
			t.Errorf("Get for nonexistent point should return zeros, got %v", got)
		}
	})

	t.Run("Set out of bounds", func(t *testing.T) {
		buffer := color.NewErrorBuffer(bounds)
		p := geometry.Point{X: 20, Y: 20}
		err := [3]float64{1.0, 2.0, 3.0}

		buffer.Set(p, err)
		got := buffer.Get(p)

		if got != [3]float64{0, 0, 0} {
			t.Error("Out of bounds Set should be ignored")
		}
	})

	t.Run("Clear", func(t *testing.T) {
		buffer := color.NewErrorBuffer(bounds)
		p := geometry.Point{X: 5, Y: 5}
		err := [3]float64{1.0, 2.0, 3.0}

		buffer.Set(p, err)
		buffer.Clear()
		got := buffer.Get(p)

		if got != [3]float64{0, 0, 0} {
			t.Error("Clear should remove all error terms")
		}
	})
}

func TestFloydSteinbergDither(t *testing.T) {
	palette := []color.Color{
		{R: 0, G: 0, B: 0, A: 255},       // Black
		{R: 255, G: 255, B: 255, A: 255}, // White
	}

	bounds := geometry.Rect{
		Min: geometry.Point{X: 0, Y: 0},
		Max: geometry.Point{X: 4, Y: 4},
	}

	t.Run("Error propagation", func(t *testing.T) {
		buffer := color.NewErrorBuffer(bounds)
		c := color.Color{R: 128, G: 128, B: 128, A: 255}

		// Dither first pixel, passing the buffer
		result := c.Dither(color.DitherFloydSteinberg, 0, 0, palette, buffer)

		// Check that error was propagated to right neighbor
		rightErr := buffer.Get(geometry.Point{X: 1, Y: 0})
		if rightErr == [3]float64{0, 0, 0} {
			t.Error("Error should be propagated to right neighbor")
		}

		// Check that error was propagated to bottom neighbors
		bottomErr := buffer.Get(geometry.Point{X: 0, Y: 1})
		if bottomErr == [3]float64{0, 0, 0} {
			t.Error("Error should be propagated to bottom neighbors")
		}

		// Result should be either black or white
		if result != palette[0] && result != palette[1] {
			t.Error("Result should be mapped to palette color")
		}
	})

	t.Run("Nil buffer fallback", func(t *testing.T) {
		c := color.Color{R: 128, G: 128, B: 128, A: 255}
		result := c.Dither(color.DitherFloydSteinberg, 0, 0, palette)

		// Should still work without buffer, falling back to nearest color
		if result != palette[0] && result != palette[1] {
			t.Error("Should fall back to nearest color with nil buffer")
		}
	})
}

func TestDitherWithEmptyPalette(t *testing.T) {
	testColor := color.Color{R: 128, G: 128, B: 128, A: 255}
	methods := []color.DitherMethod{
		color.DitherNone,
		color.DitherFloydSteinberg,
		color.DitherOrdered,
		color.DitherBayer,
	}

	for _, method := range methods {
		t.Run(fmt.Sprintf("method_%v", method), func(t *testing.T) {
			result := testColor.Dither(method, 0, 0, nil)
			if result != testColor {
				t.Errorf("Method %v with empty palette should return original color", method)
			}
		})
	}
}

func TestDitherWithSingleColorPalette(t *testing.T) {
	testColor := color.Color{R: 128, G: 128, B: 128, A: 255}
	singlePalette := []color.Color{{R: 100, G: 100, B: 100, A: 255}}
	methods := []color.DitherMethod{
		color.DitherNone,
		color.DitherFloydSteinberg,
		color.DitherOrdered,
		color.DitherBayer,
	}

	for _, method := range methods {
		t.Run(fmt.Sprintf("method_%v", method), func(t *testing.T) {
			result := testColor.Dither(method, 0, 0, singlePalette)
			if result != singlePalette[0] {
				t.Errorf("Method %v with single color palette should return that color", method)
			}
		})
	}
}

func TestFloydSteinbergWithoutBuffer(t *testing.T) {
	testColor := color.Color{R: 128, G: 128, B: 128, A: 255}
	palette := []color.Color{
		{R: 0, G: 0, B: 0, A: 255},       // Black
		{R: 255, G: 255, B: 255, A: 255}, // White
	}

	result := testColor.Dither(color.DitherFloydSteinberg, 0, 0, palette)
	if result != palette[0] && result != palette[1] {
		t.Error("Floyd-Steinberg without buffer should fall back to nearest color")
	}
}

func TestOrderedDitherPattern(t *testing.T) {
	testColor := color.Color{R: 128, G: 128, B: 128, A: 255}
	palette := []color.Color{
		{R: 0, G: 0, B: 0, A: 255},       // Black
		{R: 255, G: 255, B: 255, A: 255}, // White
	}

	results := make(map[string]color.Color)
	positions := [][2]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}}

	for _, pos := range positions {
		result := testColor.Dither(color.DitherOrdered, pos[0], pos[1], palette)
		key := fmt.Sprintf("%d,%d", pos[0], pos[1])
		results[key] = result
	}

	hasBlack := false
	hasWhite := false
	for _, c := range results {
		if c == palette[0] {
			hasBlack = true
		}
		if c == palette[1] {
			hasWhite = true
		}
	}

	if !hasBlack || !hasWhite {
		t.Error("Ordered dithering should produce both black and white pixels for 50% gray")
	}
}

func TestBayerDitherPattern2x2(t *testing.T) {
	testColor := color.Color{R: 128, G: 128, B: 128, A: 255}
	palette := []color.Color{
		{R: 0, G: 0, B: 0, A: 255},       // Black
		{R: 255, G: 255, B: 255, A: 255}, // White
	}

	var blackCount, whiteCount int
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			result := testColor.Dither(color.DitherBayer, x, y, palette)
			if result == palette[0] {
				blackCount++
			} else if result == palette[1] {
				whiteCount++
			}
		}
	}

	if blackCount == 0 || whiteCount == 0 {
		t.Errorf("Bayer dithering should produce both black and white pixels, got %d black and %d white",
			blackCount, whiteCount)
	}
}

func TestDitherWithEmptyMatrix(t *testing.T) {
	testColor := color.Color{R: 128, G: 128, B: 128, A: 255}
	palette := []color.Color{
		{R: 0, G: 0, B: 0, A: 255},
		{R: 255, G: 255, B: 255, A: 255},
	}

	originalBayer := color.Bayer4x4
	color.Bayer4x4 = color.DitherMatrix{}
	defer func() {
		color.Bayer4x4 = originalBayer
	}()

	result := testColor.Dither(color.DitherOrdered, 0, 0, palette)
	expectedNearest := testColor.Dither(color.DitherNone, 0, 0, palette)

	if result != expectedNearest {
		t.Errorf("Empty matrix should fall back to nearest color, got %v, want %v",
			result, expectedNearest)
	}
}

func TestFloydSteinbergWithNilBuffer(t *testing.T) {
	testColor := color.Color{R: 128, G: 128, B: 128, A: 255}
	palette := []color.Color{
		{R: 0, G: 0, B: 0, A: 255},
		{R: 255, G: 255, B: 255, A: 255},
	}

	result := testColor.Dither(color.DitherFloydSteinberg, 0, 0, palette, nil)
	expectedNearest := testColor.Dither(color.DitherNone, 0, 0, palette)

	if result != expectedNearest {
		t.Errorf("Nil buffer should fall back to nearest color, got %v, want %v",
			result, expectedNearest)
	}
}
