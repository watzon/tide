// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package engine

import (
	"github.com/watzon/tide/pkg/backend/terminal"
	"github.com/watzon/tide/pkg/core/capabilities"
	"github.com/watzon/tide/pkg/core/color"
	"github.com/watzon/tide/pkg/core/geometry"
	"github.com/watzon/tide/pkg/core/style"
)

// TerminalContext adapts the terminal backend to the RenderContext interface
type TerminalContext struct {
	*BaseRenderContext
	term *terminal.Terminal
}

func NewTerminalContext(term *terminal.Terminal) *TerminalContext {
	ctx := &TerminalContext{
		BaseRenderContext: NewBaseRenderContext(capabilities.Capabilities{
			ColorMode:        capabilities.ColorTrueColor,
			SupportsItalic:   true,
			SupportsBold:     true,
			SupportsKeyboard: true,
		}, term.Size()),
		term: term,
	}

	// Set initial clip rect to full terminal size
	size := term.Size()
	ctx.PushClipRect(geometry.NewRect(0, 0, size.Width, size.Height))

	return ctx
}

// Basic drawing operations
func (t *TerminalContext) Clear() {
	t.term.Clear()
}

func (t *TerminalContext) Present() error {
	return t.term.Present()
}

// Cell operations
func (t *TerminalContext) DrawCell(x, y int, ch rune, fg, bg color.Color) {
	if !t.IsInBounds(x, y) || !t.IsInClipRect(x, y) {
		return
	}
	tx, ty := t.TransformPoint(x, y)
	t.term.DrawCell(tx, ty, ch, fg, bg)
}

func (t *TerminalContext) DrawStyledCell(x, y int, ch rune, fg, bg color.Color, s style.Style) {
	if !t.IsInBounds(x, y) || !t.IsInClipRect(x, y) {
		return
	}
	tx, ty := t.TransformPoint(x, y)

	// Convert style.Style to terminal.StyleMask
	var mask terminal.StyleMask
	if s.Bold {
		mask |= terminal.StyleBold
	}
	if s.Italic {
		mask |= terminal.StyleItalic
	}
	if s.Underline {
		mask |= terminal.StyleUnderline
	}

	t.term.DrawStyledCell(tx, ty, ch, fg, bg, mask)
}

// Text operations
func (t *TerminalContext) DrawText(pos geometry.Point, text string, s style.Style) {
	if !t.IsInClipRect(pos.X, pos.Y) {
		return
	}
	tx, ty := t.TransformPoint(pos.X, pos.Y)
	t.term.DrawText(tx, ty, text, s.ForegroundColor, s.BackgroundColor, terminal.StyleMask(0))
}

// Box model operations
func (t *TerminalContext) PaintBorder(rect geometry.Rect, s style.Style) {
	// Transform rect according to current offset and clipping
	tRect := t.TransformRect(rect)
	if tRect.IsEmpty() {
		return
	}
	t.term.DrawBorder(tRect, s)
}

// Helper methods
func (t *TerminalContext) TransformPoint(x, y int) (int, int) {
	return x + t.offset.X, y + t.offset.Y
}

func (t *TerminalContext) TransformRect(r geometry.Rect) geometry.Rect {
	return geometry.Rect{
		Min: geometry.Point{
			X: r.Min.X + t.offset.X,
			Y: r.Min.Y + t.offset.Y,
		},
		Max: geometry.Point{
			X: r.Max.X + t.offset.X,
			Y: r.Max.Y + t.offset.Y,
		},
	}
}

func (t *TerminalContext) IsInBounds(x, y int) bool {
	size := t.term.Size()
	return x >= 0 && x < size.Width && y >= 0 && y < size.Height
}
