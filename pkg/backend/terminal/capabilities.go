// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package terminal

import (
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// ColorMode represents the level of color support
type ColorMode int

const (
	ColorNone ColorMode = iota
	Color16
	Color256
	ColorTrueColor
)

// Capabilities represents what the terminal supports
type Capabilities struct {
	ColorMode      ColorMode
	Unicode        bool
	Italic         bool
	Strikethrough  bool
	Mouse          bool
	ModifiedKeys   bool
	BracketedPaste bool
	URLs           bool
	Title          bool
}

// DetectCapabilities returns the terminal's capabilities
func DetectCapabilities(screen tcell.Screen) Capabilities {
	term := strings.ToLower(os.Getenv("TERM"))
	colorTerm := strings.ToLower(os.Getenv("COLORTERM"))

	caps := Capabilities{}

	// Detect color support
	caps.ColorMode = detectColorMode(term, colorTerm)

	// Check for Unicode support based on TERM
	caps.Unicode = !strings.Contains(term, "ascii") && term != "dumb"

	// Check terminal features based on TERM type
	isXterm := strings.Contains(term, "xterm")
	isTmux := strings.Contains(term, "tmux")
	isScreen := strings.Contains(term, "screen")

	// Most modern terminals support these features
	caps.Italic = isXterm || isTmux
	caps.Strikethrough = isXterm || isTmux
	caps.Mouse = isXterm || isTmux || isScreen
	caps.ModifiedKeys = isXterm || isTmux || isScreen
	caps.BracketedPaste = isXterm || isTmux

	// Check for URL support
	caps.URLs = detectURLSupport(term)

	// Check for title support
	caps.Title = detectTitleSupport(term)

	return caps
}

func detectColorMode(term, colorTerm string) ColorMode {
	// Check explicit COLORTERM setting
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		return ColorTrueColor
	}

	// Check based on TERM
	if strings.Contains(term, "256color") {
		return Color256
	}

	if strings.Contains(term, "color") || strings.Contains(term, "ansi") {
		return Color16
	}

	return ColorNone
}

func detectURLSupport(term string) bool {
	urlCapableTerms := map[string]bool{
		"iterm":        true,
		"iterm2":       true,
		"wezterm":      true,
		"konsole":      true,
		"vte":          true,
		"terminator":   true,
		"gnome":        true,
		"gnome-256":    true,
		"gnome-direct": true,
	}

	for supported := range urlCapableTerms {
		if strings.Contains(term, supported) {
			return true
		}
	}

	return false
}

func detectTitleSupport(term string) bool {
	noTitleTerms := map[string]bool{
		"dumb":   true,
		"cons25": true,
		"emacs":  true,
		"linux":  true,
		"sun":    true,
		"vt52":   true,
		"vt100":  true,
		"ansi":   true,
	}

	return !noTitleTerms[term]
}

// Add these methods to Terminal struct

func (t *Terminal) Capabilities() Capabilities {
	return DetectCapabilities(t.screen)
}

func (t *Terminal) ColorMode() ColorMode {
	return t.Capabilities().ColorMode
}

func (t *Terminal) SupportsColor() bool {
	return t.ColorMode() != ColorNone
}

func (t *Terminal) SupportsTrueColor() bool {
	return t.ColorMode() == ColorTrueColor
}

func (t *Terminal) SupportsUnicode() bool {
	return t.Capabilities().Unicode
}
