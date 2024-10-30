package engine

import (
	"testing"

	"github.com/watzon/tide/pkg/core/capabilities"
	"github.com/watzon/tide/pkg/core/geometry"
)

func TestNewBaseRenderContext(t *testing.T) {
	caps := capabilities.Capabilities{
		ColorMode:        capabilities.ColorTrueColor,
		SupportsItalic:   true,
		SupportsBold:     true,
		SupportsKeyboard: true,
	}
	size := geometry.Size{Width: 80, Height: 24}

	ctx := NewBaseRenderContext(caps, size)

	if ctx == nil {
		t.Fatal("Expected non-nil context")
	}

	// Test initial state
	if ctx.capabilities != caps {
		t.Error("Capabilities not set correctly")
	}
	if ctx.size != size {
		t.Error("Size not set correctly")
	}
	if ctx.clipRect != nil {
		t.Error("Expected nil initial clipRect")
	}
	if ctx.offset != (geometry.Point{}) {
		t.Error("Expected zero initial offset")
	}
}

func TestBaseRenderContextCapabilities(t *testing.T) {
	caps := capabilities.Capabilities{ColorMode: capabilities.ColorTrueColor}
	ctx := NewBaseRenderContext(caps, geometry.Size{})

	if got := ctx.Capabilities(); got != caps {
		t.Errorf("Capabilities() = %v, want %v", got, caps)
	}
}

func TestBaseRenderContextSize(t *testing.T) {
	size := geometry.Size{Width: 100, Height: 50}
	ctx := NewBaseRenderContext(capabilities.Capabilities{}, size)

	if got := ctx.Size(); got != size {
		t.Errorf("Size() = %v, want %v", got, size)
	}
}

func TestBaseRenderContextClipRect(t *testing.T) {
	ctx := NewBaseRenderContext(capabilities.Capabilities{}, geometry.Size{})
	rect1 := geometry.NewRect(0, 0, 10, 10)
	rect2 := geometry.NewRect(5, 5, 15, 15)

	// Test pushing clip rects
	ctx.PushClipRect(rect1)
	if ctx.clipRect == nil || ctx.clipRect.Rect != rect1 {
		t.Error("First clip rect not set correctly")
	}

	ctx.PushClipRect(rect2)
	if ctx.clipRect == nil || ctx.clipRect.Rect != rect2 {
		t.Error("Second clip rect not set correctly")
	}
	if ctx.clipRect.Next == nil || ctx.clipRect.Next.Rect != rect1 {
		t.Error("Clip rect stack not maintained correctly")
	}

	// Test popping clip rects
	ctx.PopClipRect()
	if ctx.clipRect == nil || ctx.clipRect.Rect != rect1 {
		t.Error("First clip rect not restored correctly after pop")
	}

	ctx.PopClipRect()
	if ctx.clipRect != nil {
		t.Error("Clip rect not cleared after final pop")
	}

	// Test popping empty stack
	ctx.PopClipRect()
	if ctx.clipRect != nil {
		t.Error("PopClipRect should handle empty stack gracefully")
	}
}

func TestBaseRenderContextOffset(t *testing.T) {
	ctx := NewBaseRenderContext(capabilities.Capabilities{}, geometry.Size{})

	// Test pushing offsets
	offset1 := geometry.Point{X: 10, Y: 5}
	ctx.PushOffset(offset1)
	if ctx.offset != offset1 {
		t.Errorf("First offset not set correctly, got %v, want %v", ctx.offset, offset1)
	}

	offset2 := geometry.Point{X: 5, Y: 10}
	ctx.PushOffset(offset2)
	expected := geometry.Point{X: 15, Y: 15}
	if ctx.offset != expected {
		t.Errorf("Combined offset not set correctly, got %v, want %v", ctx.offset, expected)
	}

	// Test popping offsets
	ctx.PopOffset()
	if ctx.offset != (geometry.Point{}) {
		t.Error("Offset not cleared after pop")
	}

	// Test popping empty offset
	ctx.PopOffset()
	if ctx.offset != (geometry.Point{}) {
		t.Error("PopOffset should handle empty stack gracefully")
	}
}

func TestBaseRenderContextIsInClipRect(t *testing.T) {
	ctx := NewBaseRenderContext(capabilities.Capabilities{}, geometry.Size{})

	// Test with no clip rect
	if !ctx.IsInClipRect(0, 0) {
		t.Error("Expected point to be in bounds when no clip rect is set")
	}

	// Test with clip rect
	clipRect := geometry.NewRect(10, 10, 20, 20)
	ctx.PushClipRect(clipRect)

	tests := []struct {
		x, y     int
		expected bool
		name     string
	}{
		{15, 15, true, "point inside clip rect"},
		{5, 5, false, "point outside clip rect"},
		{30, 30, false, "point far outside clip rect"},
		{10, 10, true, "point on clip rect boundary"},
		{29, 29, true, "point on far clip rect boundary"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ctx.IsInClipRect(tt.x, tt.y); got != tt.expected {
				t.Errorf("IsInClipRect(%d, %d) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}

	// Test with offset
	ctx.PushOffset(geometry.Point{X: 5, Y: 5})
	if !ctx.IsInClipRect(10, 10) { // This becomes (15,15) after offset
		t.Error("Expected point to be in bounds with offset")
	}
}

func TestBaseRenderContextTransformPoint(t *testing.T) {
	ctx := NewBaseRenderContext(capabilities.Capabilities{}, geometry.Size{})

	// Test with no offset
	x, y := ctx.TransformPoint(10, 20)
	if x != 10 || y != 20 {
		t.Errorf("TransformPoint(10, 20) = (%d, %d), want (10, 20)", x, y)
	}

	// Test with offset
	ctx.PushOffset(geometry.Point{X: 5, Y: 10})
	x, y = ctx.TransformPoint(10, 20)
	if x != 15 || y != 30 {
		t.Errorf("TransformPoint(10, 20) with offset = (%d, %d), want (15, 30)", x, y)
	}
}

func TestBaseRenderContextIsInBounds(t *testing.T) {
	size := geometry.Size{Width: 80, Height: 24}
	ctx := NewBaseRenderContext(capabilities.Capabilities{}, size)

	tests := []struct {
		x, y     int
		expected bool
		name     string
	}{
		{0, 0, true, "origin"},
		{79, 23, true, "max bounds"},
		{-1, 0, false, "left out of bounds"},
		{80, 0, false, "right out of bounds"},
		{0, -1, false, "top out of bounds"},
		{0, 24, false, "bottom out of bounds"},
		{40, 12, true, "middle point"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ctx.IsInBounds(tt.x, tt.y); got != tt.expected {
				t.Errorf("IsInBounds(%d, %d) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}
