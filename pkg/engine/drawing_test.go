package engine

import (
	"fmt"
	"testing"

	"github.com/watzon/tide/pkg/core/capabilities"
	"github.com/watzon/tide/pkg/core/color"
	"github.com/watzon/tide/pkg/core/geometry"
	"github.com/watzon/tide/pkg/core/style"
)

// mockRenderContext implements RenderContext for testing
type mockRenderContext struct {
	*BaseRenderContext
	cells map[string]struct {
		ch rune
		fg color.Color
		bg color.Color
	}
}

func newMockRenderContext(width, height int) *mockRenderContext {
	caps := capabilities.Capabilities{
		ColorMode:        capabilities.ColorTrueColor,
		SupportsItalic:   true,
		SupportsBold:     true,
		SupportsKeyboard: true,
	}

	return &mockRenderContext{
		BaseRenderContext: NewBaseRenderContext(caps, geometry.Size{Width: width, Height: height}),
		cells: make(map[string]struct {
			ch rune
			fg color.Color
			bg color.Color
		}),
	}
}

func (m *mockRenderContext) Clear() {}

func (m *mockRenderContext) Present() error { return nil }

func (m *mockRenderContext) DrawCell(x, y int, ch rune, fg, bg color.Color) {
	key := fmt.Sprintf("%d,%d", x, y)
	m.cells[key] = struct {
		ch rune
		fg color.Color
		bg color.Color
	}{ch, fg, bg}
}

func (m *mockRenderContext) DrawStyledCell(x, y int, ch rune, fg, bg color.Color, s style.Style) {
	m.DrawCell(x, y, ch, fg, bg)
}

func (m *mockRenderContext) DrawText(pos geometry.Point, text string, s style.Style) {}

func TestFillRect(t *testing.T) {
	tests := []struct {
		name     string
		rect     geometry.Rect
		fg       color.Color
		bg       color.Color
		expected map[string]struct {
			ch rune
			fg color.Color
			bg color.Color
		}
	}{
		{
			name: "1x1 rectangle",
			rect: geometry.NewRect(0, 0, 1, 1),
			fg:   color.White,
			bg:   color.Blue,
			expected: map[string]struct {
				ch rune
				fg color.Color
				bg color.Color
			}{
				"0,0": {' ', color.White, color.Blue},
			},
		},
		{
			name: "2x2 rectangle",
			rect: geometry.NewRect(1, 1, 2, 2),
			fg:   color.Red,
			bg:   color.Green,
			expected: map[string]struct {
				ch rune
				fg color.Color
				bg color.Color
			}{
				"1,1": {' ', color.Red, color.Green},
				"2,1": {' ', color.Red, color.Green},
				"1,2": {' ', color.Red, color.Green},
				"2,2": {' ', color.Red, color.Green},
			},
		},
		{
			name: "empty rectangle",
			rect: geometry.NewRect(0, 0, 0, 0),
			fg:   color.White,
			bg:   color.Black,
			expected: map[string]struct {
				ch rune
				fg color.Color
				bg color.Color
			}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newMockRenderContext(10, 10)
			FillRect(ctx, tt.rect, tt.fg, tt.bg)

			if len(ctx.cells) != len(tt.expected) {
				t.Errorf("Expected %d cells to be drawn, got %d", len(tt.expected), len(ctx.cells))
			}

			for key, expected := range tt.expected {
				if cell, ok := ctx.cells[key]; !ok {
					t.Errorf("Expected cell at %s to be drawn", key)
				} else {
					if cell.ch != expected.ch {
						t.Errorf("Expected character %v at %s, got %v", expected.ch, key, cell.ch)
					}
					if cell.fg != expected.fg {
						t.Errorf("Expected foreground color %v at %s, got %v", expected.fg, key, cell.fg)
					}
					if cell.bg != expected.bg {
						t.Errorf("Expected background color %v at %s, got %v", expected.bg, key, cell.bg)
					}
				}
			}
		})
	}
}
