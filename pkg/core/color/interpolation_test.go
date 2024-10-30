package color_test

import (
	"testing"

	"github.com/watzon/tide/pkg/core/color"
)

func TestLerp(t *testing.T) {
	tests := []struct {
		name   string
		c1     color.Color
		c2     color.Color
		t      float64
		expect color.Color
	}{
		{
			name:   "Midpoint black to white",
			c1:     color.Color{R: 0, G: 0, B: 0, A: 255},
			c2:     color.Color{R: 255, G: 255, B: 255, A: 255},
			t:      0.5,
			expect: color.Color{R: 127, G: 127, B: 127, A: 255},
		},
		{
			name:   "Start point",
			c1:     color.Color{R: 100, G: 150, B: 200, A: 255},
			c2:     color.Color{R: 200, G: 50, B: 100, A: 255},
			t:      0.0,
			expect: color.Color{R: 100, G: 150, B: 200, A: 255},
		},
		{
			name:   "End point",
			c1:     color.Color{R: 100, G: 150, B: 200, A: 255},
			c2:     color.Color{R: 200, G: 50, B: 100, A: 255},
			t:      1.0,
			expect: color.Color{R: 200, G: 50, B: 100, A: 255},
		},
		{
			name:   "Alpha interpolation",
			c1:     color.Color{R: 100, G: 100, B: 100, A: 0},
			c2:     color.Color{R: 100, G: 100, B: 100, A: 255},
			t:      0.5,
			expect: color.Color{R: 100, G: 100, B: 100, A: 127},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := color.Lerp(tt.c1, tt.c2, tt.t)
			if result != tt.expect {
				t.Errorf("Lerp() = %v, want %v", result, tt.expect)
			}
		})
	}
}

func TestGradientNormalCase(t *testing.T) {
	start := color.Color{R: 0, G: 0, B: 0, A: 255}
	end := color.Color{R: 255, G: 255, B: 255, A: 255}
	steps := 5

	colors := color.Gradient(start, end, steps)

	if len(colors) != steps {
		t.Errorf("expected %d colors, got %d", steps, len(colors))
	}

	validateGradient(t, colors, start, end)
}

func TestGradientSameColors(t *testing.T) {
	start := color.Color{R: 100, G: 100, B: 100, A: 255}
	end := color.Color{R: 100, G: 100, B: 100, A: 255}
	steps := 5

	colors := color.Gradient(start, end, steps)

	if len(colors) != steps {
		t.Errorf("expected %d colors, got %d", steps, len(colors))
	}

	validateGradient(t, colors, start, end)
}

func TestGradientSingleStep(t *testing.T) {
	start := color.Color{R: 100, G: 100, B: 100, A: 255}
	end := color.Color{R: 200, G: 200, B: 200, A: 255}

	colors := color.Gradient(start, end, 1)

	if len(colors) != 1 {
		t.Errorf("expected 1 color, got %d", len(colors))
	}
	if colors[0] != start {
		t.Errorf("expected %v, got %v", start, colors[0])
	}
}

// Helper function to validate gradient properties
func validateGradient(t *testing.T, colors []color.Color, start, end color.Color) {
	t.Helper()

	if colors[0] != start {
		t.Error("First color should match start color")
	}
	if colors[len(colors)-1] != end {
		t.Error("Last color should match end color")
	}

	for i := 1; i < len(colors); i++ {
		if colors[i].R < colors[i-1].R ||
			colors[i].G < colors[i-1].G ||
			colors[i].B < colors[i-1].B {
			t.Error("Colors should increase monotonically")
		}
	}
}

func TestMix(t *testing.T) {
	tests := []struct {
		name   string
		c1     color.Color
		c2     color.Color
		weight float64
		expect color.Color
	}{
		{
			name:   "Equal mix",
			c1:     color.Color{R: 0, G: 0, B: 0, A: 255},
			c2:     color.Color{R: 255, G: 255, B: 255, A: 255},
			weight: 0.5,
			expect: color.Color{R: 127, G: 127, B: 127, A: 255},
		},
		{
			name:   "Full first color",
			c1:     color.Color{R: 100, G: 100, B: 100, A: 255},
			c2:     color.Color{R: 200, G: 200, B: 200, A: 255},
			weight: 0.0,
			expect: color.Color{R: 100, G: 100, B: 100, A: 255},
		},
		{
			name:   "Full second color",
			c1:     color.Color{R: 100, G: 100, B: 100, A: 255},
			c2:     color.Color{R: 200, G: 200, B: 200, A: 255},
			weight: 1.0,
			expect: color.Color{R: 200, G: 200, B: 200, A: 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := color.Mix(tt.c1, tt.c2, tt.weight)
			if result != tt.expect {
				t.Errorf("Mix() = %v, want %v", result, tt.expect)
			}
		})
	}
}
