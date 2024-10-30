// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/watzon/tide/pkg/core/color"
	"github.com/watzon/tide/pkg/core/geometry"
	"github.com/watzon/tide/pkg/core/style"
	"github.com/watzon/tide/pkg/engine"
)

// Cell represents a single character cell for testing purposes
type Cell struct {
	Rune rune
	Fg   color.Color
	Bg   color.Color
}

// MockRenderContext implements engine.RenderContext for testing
type MockRenderContext struct {
	engine.RenderContext
	cells map[geometry.Point]Cell
}

func NewMockRenderContext() *MockRenderContext {
	return &MockRenderContext{
		cells: make(map[geometry.Point]Cell),
	}
}

func (m *MockRenderContext) DrawCell(x, y int, ch rune, fg, bg color.Color) {
	m.cells[geometry.Point{X: x, Y: y}] = Cell{
		Rune: ch,
		Fg:   fg,
		Bg:   bg,
	}
}

// MockChildRenderObject implements RenderObject for testing child paint calls
type MockChildRenderObject struct {
	BaseRenderObject
	painted bool
}

func NewMockChildRenderObject() *MockChildRenderObject {
	return &MockChildRenderObject{
		BaseRenderObject: BaseRenderObject{
			style: NewWidgetStyle(),
		},
	}
}

func (m *MockChildRenderObject) Paint(context engine.RenderContext) {
	m.painted = true
}

// BaseRenderObject tests
func TestBaseRenderObject_Layout(t *testing.T) {
	ro := NewBaseRenderObject(WidgetStyle{})
	constraints := NewConstraints(
		geometry.Size{Width: 10, Height: 20},
		geometry.Size{Width: 100, Height: 200},
	)

	size := ro.Layout(constraints)
	assert.Equal(t, constraints.MinSize, size)
	assert.Equal(t, constraints, ro.Constraints())
	assert.Equal(t, constraints.MinSize, ro.Size())
}

func TestBaseRenderObject_ParentChild(t *testing.T) {
	parent := NewBaseRenderObject(WidgetStyle{})
	child1 := NewBaseRenderObject(WidgetStyle{})
	child2 := NewBaseRenderObject(WidgetStyle{})

	// Test AppendChild
	parent.AppendChild(child1)
	parent.AppendChild(child2)
	assert.Equal(t, parent, child1.Parent())
	assert.Equal(t, parent, child2.Parent())
	assert.Equal(t, []RenderObject{child1, child2}, parent.Children())

	// Test RemoveChild
	parent.RemoveChild(child1)
	assert.Nil(t, child1.Parent())
	assert.Equal(t, []RenderObject{child2}, parent.Children())

	// Test ClearChildren
	parent.ClearChildren()
	assert.Empty(t, parent.Children())
	assert.Nil(t, child2.Parent())
}

func TestBaseRenderObject_Style(t *testing.T) {
	style := WidgetStyle{
		Style: style.Style{
			ForegroundColor: color.Red,
			BackgroundColor: color.Blue,
		},
	}
	ro := NewBaseRenderObject(style)
	assert.Equal(t, style, ro.Style())
}

func TestBaseRenderObject_Paint(t *testing.T) {
	ctx := NewMockRenderContext()
	style := WidgetStyle{
		Style: style.Style{
			BackgroundColor: color.Blue,
			ForegroundColor: color.Red,
		},
	}
	ro := NewBaseRenderObject(style)
	ro.size = geometry.Size{Width: 2, Height: 2}

	ro.Paint(ctx)

	// Verify background was painted
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			cell := ctx.cells[geometry.Point{X: x, Y: y}]
			assert.Equal(t, ' ', cell.Rune)
			assert.Equal(t, color.Red, cell.Fg)
			assert.Equal(t, color.Blue, cell.Bg)
		}
	}
}

// BaseRenderBox tests
func TestBaseRenderBox_Rects(t *testing.T) {
	style := WidgetStyle{
		Padding: EdgeInsetsAll(5),
		Margin:  EdgeInsetsAll(10),
	}
	box := &BaseRenderBox{
		BaseRenderObject: BaseRenderObject{
			style: style,
			size:  geometry.Size{Width: 100, Height: 100},
		},
	}

	// Test ContentRect
	contentRect := box.ContentRect()
	assert.Equal(t, geometry.Point{X: 5, Y: 5}, contentRect.Min)
	assert.Equal(t, geometry.Point{X: 95, Y: 95}, contentRect.Max)

	// Test PaddingRect
	paddingRect := box.PaddingRect()
	assert.Equal(t, geometry.Point{X: 0, Y: 0}, paddingRect.Min)
	assert.Equal(t, geometry.Point{X: 100, Y: 100}, paddingRect.Max)

	// Test BorderRect
	borderRect := box.BorderRect()
	assert.Equal(t, geometry.Point{X: -10, Y: -10}, borderRect.Min)
	assert.Equal(t, geometry.Point{X: 110, Y: 110}, borderRect.Max)

	// Test MarginRect
	marginRect := box.MarginRect()
	assert.Equal(t, geometry.Point{X: -20, Y: -20}, marginRect.Min)
	assert.Equal(t, geometry.Point{X: 120, Y: 120}, marginRect.Max)
}

func TestBaseRenderBox_Paint(t *testing.T) {
	ctx := NewMockRenderContext()
	style := WidgetStyle{
		Style: style.Style{
			BackgroundColor: color.Blue,
			ForegroundColor: color.Red,
		},
		Padding: EdgeInsetsAll(1),
	}
	box := &BaseRenderBox{
		BaseRenderObject: BaseRenderObject{
			style: style,
			size:  geometry.Size{Width: 3, Height: 3},
		},
	}

	box.Paint(ctx)

	// Verify background was painted within padding rect
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			cell := ctx.cells[geometry.Point{X: x, Y: y}]
			assert.Equal(t, ' ', cell.Rune)
			assert.Equal(t, color.Red, cell.Fg)
			assert.Equal(t, color.Blue, cell.Bg)
		}
	}
}

func TestPaintBackground(t *testing.T) {
	ctx := NewMockRenderContext()
	style := WidgetStyle{
		Style: style.Style{
			BackgroundColor: color.Blue,
			ForegroundColor: color.Red,
		},
	}
	rect := geometry.Rect{
		Min: geometry.Point{X: 1, Y: 1},
		Max: geometry.Point{X: 3, Y: 3},
	}

	paintBackground(ctx, style, rect)

	// Verify background was painted within rect
	for y := 1; y < 3; y++ {
		for x := 1; x < 3; x++ {
			cell := ctx.cells[geometry.Point{X: x, Y: y}]
			assert.Equal(t, ' ', cell.Rune)
			assert.Equal(t, color.Red, cell.Fg)
			assert.Equal(t, color.Blue, cell.Bg)
		}
	}
}

func TestBaseRenderObject_PaintWithChildren(t *testing.T) {
	ctx := NewMockRenderContext()
	parent := NewBaseRenderObject(WidgetStyle{
		Style: style.Style{
			BackgroundColor: color.Blue,
			ForegroundColor: color.Red,
		},
	})
	parent.size = geometry.Size{Width: 2, Height: 2}

	// Add two mock children
	child1 := NewMockChildRenderObject()
	child2 := NewMockChildRenderObject()
	parent.AppendChild(child1)
	parent.AppendChild(child2)

	parent.Paint(ctx)

	// Verify background was painted
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			cell := ctx.cells[geometry.Point{X: x, Y: y}]
			assert.Equal(t, ' ', cell.Rune)
			assert.Equal(t, color.Red, cell.Fg)
			assert.Equal(t, color.Blue, cell.Bg)
		}
	}

	// Verify children were painted
	assert.True(t, child1.painted, "First child should have been painted")
	assert.True(t, child2.painted, "Second child should have been painted")
}

func TestBaseRenderBox_PaintWithChildren(t *testing.T) {
	ctx := NewMockRenderContext()
	style := WidgetStyle{
		Style: style.Style{
			BackgroundColor: color.Blue,
			ForegroundColor: color.Red,
		},
		Padding: EdgeInsetsAll(1),
	}
	box := &BaseRenderBox{
		BaseRenderObject: BaseRenderObject{
			style: style,
			size:  geometry.Size{Width: 3, Height: 3},
		},
	}

	// Add two mock children
	child1 := NewMockChildRenderObject()
	child2 := NewMockChildRenderObject()
	box.AppendChild(child1)
	box.AppendChild(child2)

	box.Paint(ctx)

	// Verify background was painted within padding rect
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			cell := ctx.cells[geometry.Point{X: x, Y: y}]
			assert.Equal(t, ' ', cell.Rune)
			assert.Equal(t, color.Red, cell.Fg)
			assert.Equal(t, color.Blue, cell.Bg)
		}
	}

	// Verify children were painted
	assert.True(t, child1.painted, "First child should have been painted")
	assert.True(t, child2.painted, "Second child should have been painted")
}