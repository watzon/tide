// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package color_test

import (
	"math"
	"testing"

	"github.com/watzon/tide/pkg/core/color"
)

func TestColor(t *testing.T) {
	t.Run("RGBA fully opaque", func(t *testing.T) {
		c := color.Color{R: 255, G: 128, B: 64, A: 255}
		r, g, b, a := c.RGBA()

		if r>>8 != 255 || g>>8 != 128 || b>>8 != 64 || a>>8 != 255 {
			t.Errorf("expected RGBA (255,128,64,255), got (%d,%d,%d,%d)", r>>8, g>>8, b>>8, a>>8)
		}
	})

	t.Run("RGBA with alpha", func(t *testing.T) {
		c := color.Color{R: 255, G: 128, B: 64, A: 128}
		r, g, b, a := c.RGBA()

		// Check that alpha is properly applied
		if r>>8 != 128 || g>>8 != 64 || b>>8 != 32 || a>>8 != 128 {
			t.Errorf("expected RGBA (128,64,32,128), got (%d,%d,%d,%d)", r>>8, g>>8, b>>8, a>>8)
		}
	})
}

func TestHSLToRGB(t *testing.T) {
	tests := []struct {
		name    string
		h, s, l float64
		wantR   uint8
		wantG   uint8
		wantB   uint8
	}{
		{
			name:  "Pure red",
			h:     0,
			s:     1.0,
			l:     0.5,
			wantR: 255,
			wantG: 0,
			wantB: 0,
		},
		{
			name:  "Pure green",
			h:     120,
			s:     1.0,
			l:     0.5,
			wantR: 0,
			wantG: 255,
			wantB: 0,
		},
		{
			name:  "Pure blue",
			h:     240,
			s:     1.0,
			l:     0.5,
			wantR: 0,
			wantG: 0,
			wantB: 255,
		},
		{
			name:  "Gray (no saturation)",
			h:     0,
			s:     0,
			l:     0.5,
			wantR: 128,
			wantG: 128,
			wantB: 128,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b := color.HSLToRGB(tt.h, tt.s, tt.l)
			if r != tt.wantR || g != tt.wantG || b != tt.wantB {
				t.Errorf("HSLToRGB(%v, %v, %v) = (%v, %v, %v), want (%v, %v, %v)",
					tt.h, tt.s, tt.l, r, g, b, tt.wantR, tt.wantG, tt.wantB)
			}
		})
	}
}

func TestColorModifications(t *testing.T) {
	t.Run("Lighten", func(t *testing.T) {
		c := color.Color{R: 100, G: 100, B: 100, A: 255}
		lighter := c.Lighten(0.2)

		// Lightened color should have higher RGB values
		if lighter.R <= c.R || lighter.G <= c.G || lighter.B <= c.B {
			t.Errorf("Lightened color should have higher RGB values, got R:%d G:%d B:%d",
				lighter.R, lighter.G, lighter.B)
		}

		// Alpha should remain unchanged
		if lighter.A != c.A {
			t.Errorf("Alpha should remain unchanged, got %d, want %d", lighter.A, c.A)
		}
	})

	t.Run("Darken", func(t *testing.T) {
		c := color.Color{R: 200, G: 200, B: 200, A: 255}
		darker := c.Darken(0.2)

		// Darkened color should have lower RGB values
		if darker.R >= c.R || darker.G >= c.G || darker.B >= c.B {
			t.Errorf("Darkened color should have lower RGB values, got R:%d G:%d B:%d",
				darker.R, darker.G, darker.B)
		}

		// Alpha should remain unchanged
		if darker.A != c.A {
			t.Errorf("Alpha should remain unchanged, got %d, want %d", darker.A, c.A)
		}
	})

	t.Run("WithAlpha", func(t *testing.T) {
		c := color.Color{R: 100, G: 150, B: 200, A: 255}
		newAlpha := uint8(128)
		modified := c.WithAlpha(newAlpha)

		// RGB values should remain unchanged
		if modified.R != c.R || modified.G != c.G || modified.B != c.B {
			t.Errorf("RGB values should remain unchanged")
		}

		// Alpha should be updated
		if modified.A != newAlpha {
			t.Errorf("Alpha should be %d, got %d", newAlpha, modified.A)
		}
	})

	t.Run("Color conversion roundtrip", func(t *testing.T) {
		original := color.Color{R: 123, G: 45, B: 67}
		h, s, l := color.RGBToHSL(original.R, original.G, original.B)
		r, g, b := color.HSLToRGB(h, s, l)

		// Allow for small rounding differences (±1)
		if math.Abs(float64(r)-float64(original.R)) > 1 ||
			math.Abs(float64(g)-float64(original.G)) > 1 ||
			math.Abs(float64(b)-float64(original.B)) > 1 {
			t.Errorf("Color conversion roundtrip failed: original(%d,%d,%d) got(%d,%d,%d)",
				original.R, original.G, original.B, r, g, b)
		}
	})
}

func TestColorDistance(t *testing.T) {
	tests := []struct {
		name     string
		c1, c2   color.Color
		wantDist float64
	}{
		{
			name:     "Same color",
			c1:       color.Color{R: 100, G: 100, B: 100, A: 255},
			c2:       color.Color{R: 100, G: 100, B: 100, A: 255},
			wantDist: 0,
		},
		{
			name:     "Black to white",
			c1:       color.Color{R: 0, G: 0, B: 0, A: 255},
			c2:       color.Color{R: 255, G: 255, B: 255, A: 255},
			wantDist: 441.67295593006372, // sqrt(255^2 * 3)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dist := color.ColorDistance(tt.c1, tt.c2)
			if dist != tt.wantDist {
				t.Errorf("ColorDistance() = %v, want %v", dist, tt.wantDist)
			}
		})
	}
}
