// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package color

import (
	"math"
)

var (
	// Primary Colors
	Red   = Color{R: 255, G: 0, B: 0, A: 255}
	Green = Color{R: 0, G: 255, B: 0, A: 255}
	Blue  = Color{R: 0, G: 0, B: 255, A: 255}

	// Secondary Colors
	Yellow  = Color{R: 255, G: 255, B: 0, A: 255}
	Cyan    = Color{R: 0, G: 255, B: 255, A: 255}
	Magenta = Color{R: 255, G: 0, B: 255, A: 255}

	// Monochrome
	Black = Color{R: 0, G: 0, B: 0, A: 255}
	White = Color{R: 255, G: 255, B: 255, A: 255}
	Gray  = Color{R: 128, G: 128, B: 128, A: 255}

	// Common Web Colors
	Orange = Color{R: 255, G: 165, B: 0, A: 255}
	Purple = Color{R: 128, G: 0, B: 128, A: 255}
	Brown  = Color{R: 165, G: 42, B: 42, A: 255}
	Pink   = Color{R: 255, G: 192, B: 203, A: 255}

	// UI Colors
	Silver    = Color{R: 192, G: 192, B: 192, A: 255}
	LightGray = Color{R: 211, G: 211, B: 211, A: 255}
	DarkGray  = Color{R: 64, G: 64, B: 64, A: 255}
	Navy      = Color{R: 0, G: 0, B: 128, A: 255}
	Teal      = Color{R: 0, G: 128, B: 128, A: 255}
	Maroon    = Color{R: 128, G: 0, B: 0, A: 255}
	Olive     = Color{R: 128, G: 128, B: 0, A: 255}

	// Material Design-inspired Colors
	Primary = Color{R: 33, G: 150, B: 243, A: 255} // Light Blue
	Success = Color{R: 76, G: 175, B: 80, A: 255}  // Green
	Warning = Color{R: 255, G: 152, B: 0, A: 255}  // Orange
	Error   = Color{R: 244, G: 67, B: 54, A: 255}  // Red
	Info    = Color{R: 3, G: 169, B: 244, A: 255}  // Light Blue

	// Transparent
	Transparent = Color{R: 0, G: 0, B: 0, A: 0}

	// Reds
	DarkRed   = Color{R: 139, G: 0, B: 0, A: 255}
	IndianRed = Color{R: 205, G: 92, B: 92, A: 255}
	Crimson   = Color{R: 220, G: 20, B: 60, A: 255}

	// Greens
	ForestGreen = Color{R: 34, G: 139, B: 34, A: 255}
	LimeGreen   = Color{R: 50, G: 205, B: 50, A: 255}
	SeaGreen    = Color{R: 46, G: 139, B: 87, A: 255}

	// Blues
	RoyalBlue   = Color{R: 65, G: 105, B: 225, A: 255}
	SteelBlue   = Color{R: 70, G: 130, B: 180, A: 255}
	DeepSkyBlue = Color{R: 0, G: 191, B: 255, A: 255}

	// Yellows
	Gold      = Color{R: 255, G: 215, B: 0, A: 255}
	Goldenrod = Color{R: 218, G: 165, B: 32, A: 255}
	Khaki     = Color{R: 240, G: 230, B: 140, A: 255}

	// Purples
	Violet = Color{R: 238, G: 130, B: 238, A: 255}
	Orchid = Color{R: 218, G: 112, B: 214, A: 255}
	Plum   = Color{R: 221, G: 160, B: 221, A: 255}

	// Browns
	SaddleBrown = Color{R: 139, G: 69, B: 19, A: 255}
	Sienna      = Color{R: 160, G: 82, B: 45, A: 255}
	Peru        = Color{R: 205, G: 133, B: 63, A: 255}
)

type Color struct {
	R, G, B uint8
	A       uint8
}

func (c Color) RGBA() (uint32, uint32, uint32, uint32) {
	r := uint32(c.R)
	g := uint32(c.G)
	b := uint32(c.B)
	a := uint32(c.A)

	if a == 0xff {
		return r << 8, g << 8, b << 8, a << 8
	}

	r = (r * a) / 0xff
	g = (g * a) / 0xff
	b = (b * a) / 0xff

	return r << 8, g << 8, b << 8, a << 8
}

// RGBToHSL converts RGB values to HSL (Hue, Saturation, Lightness)
func RGBToHSL(r, g, b uint8) (h, s, l float64) {
	fr := float64(r) / 255.0
	fg := float64(g) / 255.0
	fb := float64(b) / 255.0

	max := math.Max(math.Max(fr, fg), fb)
	min := math.Min(math.Min(fr, fg), fb)

	l = (max + min) / 2.0

	if max == min {
		// achromatic
		return 0, 0, l
	}

	d := max - min
	if l > 0.5 {
		s = d / (2.0 - max - min)
	} else {
		s = d / (max + min)
	}

	switch max {
	case fr:
		h = (fg - fb) / d
		if fg < fb {
			h += 6
		}
	case fg:
		h = (fb-fr)/d + 2
	case fb:
		h = (fr-fg)/d + 4
	}
	h *= 60

	return h, s, l
}

// HSLToRGB converts HSL (Hue, Saturation, Lightness) to RGB values
func HSLToRGB(h, s, l float64) (r, g, b uint8) {
	if s == 0 {
		// achromatic - round to nearest integer
		v := uint8(math.Round(l * 255))
		return v, v, v
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	hue := h / 360 // normalize hue to 0-1
	tr := hue + 1.0/3.0
	tg := hue
	tb := hue - 1.0/3.0

	// Helper function to convert hue to RGB
	hueToRGB := func(t float64) float64 {
		if t < 0 {
			t += 1
		}
		if t > 1 {
			t -= 1
		}
		if t < 1.0/6.0 {
			return p + (q-p)*6*t
		}
		if t < 1.0/2.0 {
			return q
		}
		if t < 2.0/3.0 {
			return p + (q-p)*(2.0/3.0-t)*6
		}
		return p
	}

	return uint8(hueToRGB(tr) * 255),
		uint8(hueToRGB(tg) * 255),
		uint8(hueToRGB(tb) * 255)
}

// Lighten returns a new Color with increased lightness
func (c Color) Lighten(amount float64) Color {
	h, s, l := RGBToHSL(c.R, c.G, c.B)
	l = math.Min(1.0, l+amount)
	r, g, b := HSLToRGB(h, s, l)
	return Color{R: r, G: g, B: b, A: c.A}
}

// Darken returns a new Color with decreased lightness
func (c Color) Darken(amount float64) Color {
	h, s, l := RGBToHSL(c.R, c.G, c.B)
	l = math.Max(0.0, l-amount)
	r, g, b := HSLToRGB(h, s, l)
	return Color{R: r, G: g, B: b, A: c.A}
}

// WithAlpha returns a new Color with the specified alpha value
func (c Color) WithAlpha(alpha uint8) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: alpha}
}

// ColorDistance calculates the Euclidean distance between two colors in RGB space.
// Returns a value between 0 (identical colors) and ~441.67 (distance between black and white).
func ColorDistance(c1, c2 Color) float64 {
	dr := float64(c1.R) - float64(c2.R)
	dg := float64(c1.G) - float64(c2.G)
	db := float64(c1.B) - float64(c2.B)
	return math.Sqrt(dr*dr + dg*dg + db*db)
}
