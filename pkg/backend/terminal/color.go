// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// pkg/backend/terminal/color.go
package terminal

import (
	"fmt"
	"math"
	"sync"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/watzon/tide/pkg/core"
)

// colorCache provides thread-safe caching of color conversions
type colorCache struct {
	sync.RWMutex
	trueColors map[core.Color]tcell.Color
	palette256 map[core.Color]tcell.Color
	palette16  map[core.Color]tcell.Color
}

func newColorCache() *colorCache {
	return &colorCache{
		trueColors: make(map[core.Color]tcell.Color),
		palette256: make(map[core.Color]tcell.Color),
		palette16:  make(map[core.Color]tcell.Color),
	}
}

// ColorOptimizer handles color optimization and caching
type ColorOptimizer struct {
	cache *colorCache
	mode  ColorMode
}

func NewColorOptimizer(mode ColorMode) *ColorOptimizer {
	return &ColorOptimizer{
		cache: newColorCache(),
		mode:  mode,
	}
}

// GetColor returns an optimized tcell.Color for the given core.Color
func (co *ColorOptimizer) GetColor(c core.Color) tcell.Color {
	// Handle transparent/nil colors
	if c.A == 0 {
		return tcell.ColorDefault
	}

	// Check cache first
	co.cache.RLock()
	var cached tcell.Color
	var ok bool

	switch co.mode {
	case ColorTrueColor:
		cached, ok = co.cache.trueColors[c]
	case Color256:
		cached, ok = co.cache.palette256[c]
	case Color16:
		cached, ok = co.cache.palette16[c]
	default:
		co.cache.RUnlock()
		return tcell.ColorDefault
	}

	if ok {
		co.cache.RUnlock()
		return cached
	}
	co.cache.RUnlock()

	// Convert color based on mode
	var result tcell.Color
	switch co.mode {
	case ColorTrueColor:
		result = co.convertTrueColor(c)
	case Color256:
		result = co.convert256Color(c)
	case Color16:
		result = co.convert16Color(c)
	default:
		return tcell.ColorDefault
	}

	// Cache the result
	co.cache.Lock()
	switch co.mode {
	case ColorTrueColor:
		co.cache.trueColors[c] = result
	case Color256:
		co.cache.palette256[c] = result
	case Color16:
		co.cache.palette16[c] = result
	}
	co.cache.Unlock()

	return result
}

func (co *ColorOptimizer) convertTrueColor(c core.Color) tcell.Color {
	return tcell.NewRGBColor(int32(c.R), int32(c.G), int32(c.B))
}

func (co *ColorOptimizer) convert256Color(c core.Color) tcell.Color {
	// Standard 216 color cube (6x6x6)
	if c.R == c.G && c.G == c.B {
		// Grayscale (24 levels)
		if c.R < 8 {
			return tcell.PaletteColor(16) // black
		}
		if c.R > 238 {
			return tcell.PaletteColor(231) // white
		}
		return tcell.PaletteColor(232 + int((c.R-8)/10))
	}

	// Convert to 6x6x6 color cube
	r := int(math.Round(float64(c.R) / 51.0))
	g := int(math.Round(float64(c.G) / 51.0))
	b := int(math.Round(float64(c.B) / 51.0))

	// Calculate the color index in the 6x6x6 cube
	return tcell.PaletteColor(16 + (36 * r) + (6 * g) + b)
}

func (co *ColorOptimizer) convert16Color(c core.Color) tcell.Color {
	// For debugging, let's log the intensity decision for pure colors
	if testing.Verbose() {
		maxChannel := max(max(c.R, c.G), c.B)
		minChannel := min(min(c.R, c.G), c.B)
		_, s, l := rgbToHsl(c.R, c.G, c.B)
		fmt.Printf("Color RGB(%d,%d,%d) - max: %d, min: %d, HSL(s: %.2f, l: %.2f) - intense: %v\n",
			c.R, c.G, c.B, maxChannel, minChannel, s, l, isIntenseColor(c))
	}

	h, s, l := rgbToHsl(c.R, c.G, c.B)

	// Handle grayscale colors first
	if s < 0.2 {
		if l < 0.2 {
			return tcell.ColorBlack
		}
		if l > 0.8 {
			return tcell.ColorWhite
		}
		return tcell.ColorGray
	}

	bright := isIntenseColor(c)

	switch {
	case h < 30 || h >= 330:
		return pickColor(bright, tcell.ColorMaroon, tcell.ColorRed)
	case h < 90:
		return pickColor(bright, tcell.ColorOlive, tcell.ColorYellow)
	case h < 150:
		return pickColor(bright, tcell.ColorGreen, tcell.ColorLime)
	case h < 210:
		return pickColor(bright, tcell.ColorTeal, tcell.ColorAqua)
	case h < 270:
		return pickColor(bright, tcell.ColorNavy, tcell.ColorBlue)
	default:
		return pickColor(bright, tcell.ColorPurple, tcell.ColorFuchsia)
	}
}

// Helper for determining relative color intensity
func isIntenseColor(c core.Color) bool {
	maxChannel := max(max(c.R, c.G), c.B)
	minChannel := min(min(c.R, c.G), c.B)

	// Pure colors (like 255,0,0) should NOT be considered intense
	if maxChannel == 255 && minChannel == 0 {
		return false
	}

	// Bright variants (like 255,128,128) should be considered intense
	if maxChannel > 128 && minChannel > 64 {
		return true
	}

	// For other cases, use HSL
	_, s, l := rgbToHsl(c.R, c.G, c.B)
	return l > 0.6 && s < 0.8
}

// Helper functions for finding min/max
func max(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}

func min(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

// Helper functions

func pickColor(bright bool, dark, light tcell.Color) tcell.Color {
	if bright {
		return light
	}
	return dark
}

func rgbToHsl(r, g, b uint8) (h, s, l float64) {
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

// Add color optimizer to Terminal struct
func (t *Terminal) initColorOptimizer() {
	if t.colorOptimizer == nil {
		t.colorOptimizer = NewColorOptimizer(t.ColorMode())
	}
}
