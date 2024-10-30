package engine

import (
	"github.com/watzon/tide/pkg/core/capabilities"
	"github.com/watzon/tide/pkg/core/color"
	"github.com/watzon/tide/pkg/core/geometry"
	"github.com/watzon/tide/pkg/core/style"
)

// RenderContext provides backend-specific rendering capabilities
type RenderContext interface {
	// Backend capabilities
	Capabilities() capabilities.Capabilities

	// Drawing surface information
	Size() geometry.Size

	// Basic drawing operations
	Clear()
	Present() error

	// Cell operations
	DrawCell(x, y int, ch rune, fg, bg color.Color)
	DrawStyledCell(x, y int, ch rune, fg, bg color.Color, style style.Style)

	// Text operations
	DrawText(pos geometry.Point, text string, style style.Style)

	// Clipping
	PushClipRect(rect geometry.Rect)
	PopClipRect()

	// Transformation
	PushOffset(offset geometry.Point)
	PopOffset()
}

// ClipRect represents a clipping rectangle
type ClipRect struct {
	Rect geometry.Rect
	Next *ClipRect
}

// BaseRenderContext provides common functionality for render contexts
type BaseRenderContext struct {
	capabilities capabilities.Capabilities
	size         geometry.Size
	clipRect     *ClipRect
	offset       geometry.Point
}

func NewBaseRenderContext(caps capabilities.Capabilities, size geometry.Size) *BaseRenderContext {
	return &BaseRenderContext{
		capabilities: caps,
		size:         size,
	}
}

func (c *BaseRenderContext) Capabilities() capabilities.Capabilities {
	return c.capabilities
}

func (c *BaseRenderContext) Size() geometry.Size {
	return c.size
}

func (c *BaseRenderContext) PushClipRect(rect geometry.Rect) {
	c.clipRect = &ClipRect{
		Rect: rect,
		Next: c.clipRect,
	}
}

func (c *BaseRenderContext) PopClipRect() {
	if c.clipRect != nil {
		c.clipRect = c.clipRect.Next
	}
}

func (c *BaseRenderContext) PushOffset(offset geometry.Point) {
	c.offset = geometry.Point{
		X: c.offset.X + offset.X,
		Y: c.offset.Y + offset.Y,
	}
}

func (c *BaseRenderContext) PopOffset() {
	// Offsets should be balanced with pushes
	if c.offset.X != 0 || c.offset.Y != 0 {
		c.offset = geometry.Point{}
	}
}

// Helper methods for implementations

// IsInClipRect checks if a point is within the current clip rect
func (c *BaseRenderContext) IsInClipRect(x, y int) bool {
	if c.clipRect == nil {
		return true
	}

	// Apply offset
	x += c.offset.X
	y += c.offset.Y

	return x >= c.clipRect.Rect.Min.X &&
		x < c.clipRect.Rect.Max.X &&
		y >= c.clipRect.Rect.Min.Y &&
		y < c.clipRect.Rect.Max.Y
}

// TransformPoint applies the current offset to a point
func (c *BaseRenderContext) TransformPoint(x, y int) (int, int) {
	return x + c.offset.X, y + c.offset.Y
}

// IsInBounds checks if a point is within the render context bounds
func (c *BaseRenderContext) IsInBounds(x, y int) bool {
	return x >= 0 && x < c.size.Width &&
		y >= 0 && y < c.size.Height
}

// MockRenderContext provides a test implementation of RenderContext
type MockRenderContext struct {
	*BaseRenderContext
	DrawCellCalls  []DrawCellCall
	DrawTextCalls  []DrawTextCall
	ClearCalled    bool
	PresentCalled  bool
	ClipRectPushes []geometry.Rect
	ClipRectPops   int
	OffsetPushes   []geometry.Point
	OffsetPops     int
}

type DrawCellCall struct {
	X, Y   int
	Char   rune
	Fg, Bg color.Color
	Style  style.Style
}

type DrawTextCall struct {
	Pos   geometry.Point
	Text  string
	Style style.Style
}

func NewMockRenderContext(size geometry.Size) *MockRenderContext {
	return &MockRenderContext{
		BaseRenderContext: NewBaseRenderContext(capabilities.Capabilities{
			ColorMode: capabilities.ColorTrueColor,
		}, size),
		DrawCellCalls: make([]DrawCellCall, 0),
		DrawTextCalls: make([]DrawTextCall, 0),
	}
}

func (c *MockRenderContext) Clear() {
	c.ClearCalled = true
}

func (c *MockRenderContext) Present() error {
	c.PresentCalled = true
	return nil
}

func (c *MockRenderContext) DrawCell(x, y int, ch rune, fg, bg color.Color) {
	c.DrawCellCalls = append(c.DrawCellCalls, DrawCellCall{
		X: x, Y: y,
		Char: ch,
		Fg:   fg, Bg: bg,
	})
}

func (c *MockRenderContext) DrawStyledCell(x, y int, ch rune, fg, bg color.Color, s style.Style) {
	c.DrawCellCalls = append(c.DrawCellCalls, DrawCellCall{
		X: x, Y: y,
		Char: ch,
		Fg:   fg, Bg: bg,
		Style: s,
	})
}

func (c *MockRenderContext) DrawText(pos geometry.Point, text string, s style.Style) {
	c.DrawTextCalls = append(c.DrawTextCalls, DrawTextCall{
		Pos:   pos,
		Text:  text,
		Style: s,
	})
}

func (c *MockRenderContext) PushClipRect(rect geometry.Rect) {
	c.ClipRectPushes = append(c.ClipRectPushes, rect)
	c.BaseRenderContext.PushClipRect(rect)
}

func (c *MockRenderContext) PopClipRect() {
	c.ClipRectPops++
	c.BaseRenderContext.PopClipRect()
}

func (c *MockRenderContext) PushOffset(offset geometry.Point) {
	c.OffsetPushes = append(c.OffsetPushes, offset)
	c.BaseRenderContext.PushOffset(offset)
}

func (c *MockRenderContext) PopOffset() {
	c.OffsetPops++
	c.BaseRenderContext.PopOffset()
}
