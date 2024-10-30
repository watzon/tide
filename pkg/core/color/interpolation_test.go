package color_test

import (
	"testing"

	"github.com/watzon/tide/pkg/core/color"
)

func TestInterpolation(t *testing.T) {
	t.Run("Lerp", func(t *testing.T) {
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
	})

	t.Run("Gradient", func(t *testing.T) {
		t.Run("Normal case", func(t *testing.T) {
			start := color.Color{R: 0, G: 0, B: 0, A: 255}
			end := color.Color{R: 255, G: 255, B: 255, A: 255}
			steps := 5

			colors := color.Gradient(start, end, steps)

			// Check length
			if len(colors) != steps {
				t.Errorf("expected %d colors, got %d", steps, len(colors))
			}

			// Check start and end colors
			if colors[0] != start {
				t.Error("First color should match start color")
			}
			if colors[len(colors)-1] != end {
				t.Error("Last color should match end color")
			}

			// Check that colors are monotonically increasing
			for i := 1; i < len(colors); i++ {
				if colors[i].R < colors[i-1].R ||
					colors[i].G < colors[i-1].G ||
					colors[i].B < colors[i-1].B {
					t.Error("Colors should increase monotonically")
				}
			}
		})

		t.Run("Edge cases", func(t *testing.T) {
			start := color.Color{R: 100, G: 100, B: 100, A: 255}
			end := color.Color{R: 100, G: 100, B: 100, A: 255}
			steps := 5

			colors := color.Gradient(start, end, steps)

			// Check length
			if len(colors) != steps {
				t.Errorf("expected %d colors, got %d", steps, len(colors))
			}

			// Check start and end colors
			if colors[0] != start {
				t.Error("First color should match start color")
			}
			if colors[len(colors)-1] != end {
				t.Error("Last color should match end color")
			}

			// Check that colors are monotonically increasing
			for i := 1; i < len(colors); i++ {
				if colors[i].R < colors[i-1].R ||
					colors[i].G < colors[i-1].G ||
					colors[i].B < colors[i-1].B {
					t.Error("Colors should increase monotonically")
				}
			}
		})

		t.Run("Single step", func(t *testing.T) {
			start := color.Color{R: 100, G: 100, B: 100, A: 255}
			end := color.Color{R: 200, G: 200, B: 200, A: 255}
			steps := 1

			colors := color.Gradient(start, end, steps)

			if len(colors) != 1 {
				t.Errorf("expected 1 color, got %d", len(colors))
			}
			if colors[0] != start {
				t.Errorf("expected %v, got %v", start, colors[0])
			}
		})
	})

	t.Run("Mix", func(t *testing.T) {
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
	})
}
