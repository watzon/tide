// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// pkg/backend/terminal/capabilities_test.go
package terminal_test

import (
	"os"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/watzon/tide/pkg/backend/terminal"
	"github.com/watzon/tide/pkg/core"
)

func withEnv(env map[string]string, f func()) {
	// Save original env
	oldEnv := make(map[string]string)
	for k := range env {
		oldEnv[k] = os.Getenv(k)
	}

	// Set new env
	for k, v := range env {
		os.Setenv(k, v)
	}

	// Run function
	f()

	// Restore original env
	for k, v := range oldEnv {
		if v == "" {
			os.Unsetenv(k)
		} else {
			os.Setenv(k, v)
		}
	}
}

func TestColorModeDetection(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		wantMode terminal.ColorMode
	}{
		{
			name: "true color via COLORTERM",
			env: map[string]string{
				"TERM":      "xterm",
				"COLORTERM": "truecolor",
			},
			wantMode: terminal.ColorTrueColor,
		},
		{
			name: "256 colors via TERM",
			env: map[string]string{
				"TERM":      "xterm-256color",
				"COLORTERM": "",
			},
			wantMode: terminal.Color256,
		},
		{
			name: "16 colors",
			env: map[string]string{
				"TERM":      "xterm-color",
				"COLORTERM": "",
			},
			wantMode: terminal.Color16,
		},
		{
			name: "no color",
			env: map[string]string{
				"TERM":      "dumb",
				"COLORTERM": "",
			},
			wantMode: terminal.ColorNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withEnv(tt.env, func() {
				ctx := setupTest(t)
				defer ctx.term.Shutdown()

				caps := ctx.term.Capabilities()
				if caps.ColorMode != tt.wantMode {
					t.Errorf("got color mode %v, want %v", caps.ColorMode, tt.wantMode)
				}
			})
		})
	}
}

func TestColorDrawing(t *testing.T) {
	tests := []struct {
		name   string
		fg     core.Color
		bg     core.Color
		env    map[string]string
		verify func(*testing.T, tcell.SimulationScreen, tcell.Style)
	}{
		{
			name: "true color support",
			fg:   core.Color{R: 123, G: 45, B: 67, A: 255},
			bg:   core.Color{R: 89, G: 156, B: 234, A: 255},
			env: map[string]string{
				"TERM":      "xterm-direct",
				"COLORTERM": "truecolor",
			},
			verify: func(t *testing.T, screen tcell.SimulationScreen, style tcell.Style) {
				fg, bg, _ := style.Decompose()
				fgr, fgg, fgb := fg.RGB()
				bgr, bgg, bgb := bg.RGB()

				if fgr != 123 || fgg != 45 || fgb != 67 {
					t.Errorf("foreground color mismatch, got RGB(%d,%d,%d)", fgr, fgg, fgb)
				}
				if bgr != 89 || bgg != 156 || bgb != 234 {
					t.Errorf("background color mismatch, got RGB(%d,%d,%d)", bgr, bgg, bgb)
				}
			},
		},
		{
			name: "basic color fallback",
			fg:   core.Color{R: 255, G: 0, B: 0, A: 255}, // Pure red
			bg:   core.Color{R: 0, G: 0, B: 255, A: 255}, // Pure blue
			env: map[string]string{
				"TERM":      "xterm-color",
				"COLORTERM": "",
			},
			verify: func(t *testing.T, screen tcell.SimulationScreen, style tcell.Style) {
				fg, bg, _ := style.Decompose()

				t.Logf("Expected fg: %v (%T), got: %v (%T)", tcell.ColorMaroon, tcell.ColorMaroon, fg, fg)
				t.Logf("Expected bg: %v (%T), got: %v (%T)", tcell.ColorNavy, tcell.ColorNavy, bg, bg)

				if fg != tcell.ColorMaroon {
					r, g, b := fg.RGB()
					t.Errorf("expected foreground color to be maroon, got %v (RGB: %d,%d,%d)", fg, r, g, b)
				}
				if bg != tcell.ColorNavy {
					r, g, b := bg.RGB()
					t.Errorf("expected background color to be navy, got %v (RGB: %d,%d,%d)", bg, r, g, b)
				}
			},
		},
		{
			name: "bright basic colors",
			fg:   core.Color{R: 255, G: 128, B: 128, A: 255}, // Bright red
			bg:   core.Color{R: 128, G: 128, B: 255, A: 255}, // Bright blue
			env: map[string]string{
				"TERM":      "xterm-color",
				"COLORTERM": "",
			},
			verify: func(t *testing.T, screen tcell.SimulationScreen, style tcell.Style) {
				fg, bg, _ := style.Decompose()

				// These should map to bright colors due to higher brightness
				if fg != tcell.ColorRed {
					t.Errorf("expected foreground color to be bright red, got %v", fg)
				}
				if bg != tcell.ColorBlue {
					t.Errorf("expected background color to be bright blue, got %v", bg)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withEnv(tt.env, func() {
				ctx := setupTest(t)
				defer ctx.term.Shutdown()

				// Draw a character with our test colors
				ctx.term.DrawCell(0, 0, 'X', tt.fg, tt.bg)
				ctx.term.Present()

				// Get the style from the simulation screen
				simScreen := ctx.screen.(tcell.SimulationScreen)
				_, _, style, _ := simScreen.GetContent(0, 0)

				tt.verify(t, simScreen, style)
			})
		})
	}
}

func TestStyleAttributes(t *testing.T) {
	tests := []struct {
		name       string
		style      terminal.StyleMask
		verifyAttr func(tcell.Style) bool
	}{
		{
			name:  "bold",
			style: terminal.StyleBold,
			verifyAttr: func(s tcell.Style) bool {
				_, _, attrs := s.Decompose()
				return attrs&tcell.AttrBold != 0
			},
		},
		{
			name:  "italic",
			style: terminal.StyleItalic,
			verifyAttr: func(s tcell.Style) bool {
				_, _, attrs := s.Decompose()
				return attrs&tcell.AttrItalic != 0
			},
		},
		{
			name:  "underline",
			style: terminal.StyleUnderline,
			verifyAttr: func(s tcell.Style) bool {
				_, _, attrs := s.Decompose()
				return attrs&tcell.AttrUnderline != 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := setupTest(t)
			defer ctx.term.Shutdown()

			// Draw with style
			ctx.term.DrawStyledCell(0, 0, 'X',
				core.Color{R: 255, G: 255, B: 255, A: 255},
				core.Color{R: 0, G: 0, B: 0, A: 255},
				tt.style,
			)
			ctx.term.Present()

			// Verify style was applied
			simScreen := ctx.screen.(tcell.SimulationScreen)
			_, _, style, _ := simScreen.GetContent(0, 0)

			if !tt.verifyAttr(style) {
				t.Errorf("style %s was not applied correctly", tt.name)
			}
		})
	}
}

func TestCapabilityDetection(t *testing.T) {
	tests := []struct {
		name   string
		env    map[string]string
		verify func(*testing.T, terminal.Capabilities)
	}{
		{
			name: "modern terminal",
			env: map[string]string{
				"TERM": "xterm-256color",
			},
			verify: func(t *testing.T, caps terminal.Capabilities) {
				if !caps.Mouse {
					t.Error("mouse support should be enabled for xterm")
				}
				if !caps.Unicode {
					t.Error("unicode should be supported in xterm")
				}
				if !caps.BracketedPaste {
					t.Error("bracketed paste should be supported in xterm")
				}
			},
		},
		{
			name: "basic terminal",
			env: map[string]string{
				"TERM": "dumb",
			},
			verify: func(t *testing.T, caps terminal.Capabilities) {
				if caps.Mouse {
					t.Error("mouse support should be disabled for dumb terminal")
				}
				if caps.Unicode {
					t.Error("unicode should not be supported in dumb terminal")
				}
				if caps.Title {
					t.Error("title support should be disabled for dumb terminal")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withEnv(tt.env, func() {
				ctx := setupTest(t)
				defer ctx.term.Shutdown()

				caps := ctx.term.Capabilities()
				tt.verify(t, caps)
			})
		})
	}
}
