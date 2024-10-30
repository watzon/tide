// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package widget

import (
	"github.com/watzon/tide/pkg/core/geometry"
	"github.com/watzon/tide/pkg/engine"
)

// RenderObject is the base interface for all render objects
type RenderObject interface {
	// Layout and sizing
	Layout(constraints Constraints) geometry.Size
	Size() geometry.Size
	Constraints() Constraints

	// Parent/child relationships
	Parent() RenderObject
	Children() []RenderObject

	// Style and appearance
	Style() WidgetStyle

	// Backend-specific rendering
	Paint(context engine.RenderContext)
}

// BaseRenderObject provides default implementation for RenderObjects
type BaseRenderObject struct {
	size        geometry.Size
	constraints Constraints
	style       WidgetStyle
	parent      RenderObject
	children    []RenderObject
}

// Layout and sizing
func (r *BaseRenderObject) Layout(constraints Constraints) geometry.Size {
	r.constraints = constraints
	// Default layout just uses minimum size
	r.size = constraints.Constrain(constraints.MinSize)
	return r.size
}

func (r *BaseRenderObject) Size() geometry.Size {
	return r.size
}

func (r *BaseRenderObject) Constraints() Constraints {
	return r.constraints
}

// Parent/child relationships
func (r *BaseRenderObject) Parent() RenderObject {
	return r.parent
}

func (r *BaseRenderObject) Children() []RenderObject {
	return r.children
}

// Style access
func (r *BaseRenderObject) Style() WidgetStyle {
	return r.style
}

// Child management
func (r *BaseRenderObject) AppendChild(child RenderObject) {
	if baseChild, ok := child.(*BaseRenderObject); ok {
		baseChild.parent = r
	}
	r.children = append(r.children, child)
}

func (r *BaseRenderObject) RemoveChild(child RenderObject) {
	for i, c := range r.children {
		if c == child {
			if baseChild, ok := child.(*BaseRenderObject); ok {
				baseChild.parent = nil
			}
			r.children = append(r.children[:i], r.children[i+1:]...)
			return
		}
	}
}

func (r *BaseRenderObject) ClearChildren() {
	for _, child := range r.children {
		if baseChild, ok := child.(*BaseRenderObject); ok {
			baseChild.parent = nil
		}
	}
	r.children = nil
}

// Paint provides a default implementation that paints children
func (r *BaseRenderObject) Paint(context engine.RenderContext) {
	// Paint background if style specifies it
	if r.style.BackgroundColor.A > 0 {
		engine.FillRect(
			context,
			geometry.Rect{
				Min: geometry.Point{X: 0, Y: 0},
				Max: geometry.Point{X: r.size.Width, Y: r.size.Height},
			},
			r.style.ForegroundColor,
			r.style.BackgroundColor,
		)
	}

	// Paint children
	for _, child := range r.children {
		child.Paint(context)
	}
}

// Helper functions

// paintBackground fills a rectangle with the background color
func paintBackground(ctx engine.RenderContext, style WidgetStyle, rect geometry.Rect) {
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			ctx.DrawCell(x, y, ' ', style.ForegroundColor, style.BackgroundColor)
		}
	}
}

// NewBaseRenderObject creates a new BaseRenderObject with the given style
func NewBaseRenderObject(style WidgetStyle) *BaseRenderObject {
	return &BaseRenderObject{
		style:    style,
		children: make([]RenderObject, 0),
	}
}

// RenderBox adds box-model functionality to RenderObject
type RenderBox interface {
	RenderObject

	// Box model
	PaintBorder(context engine.RenderContext)
	PaintBackground(context engine.RenderContext)
	PaintContent(context engine.RenderContext)

	// Layout helpers
	ContentRect() geometry.Rect
	PaddingRect() geometry.Rect
	BorderRect() geometry.Rect
	MarginRect() geometry.Rect
}

// BaseRenderBox provides box model implementation
type BaseRenderBox struct {
	BaseRenderObject
}

func (r *BaseRenderBox) Paint(context engine.RenderContext) {
	r.PaintBackground(context)
	r.PaintBorder(context)
	r.PaintContent(context)
}

func (r *BaseRenderBox) PaintBackground(context engine.RenderContext) {
	if r.style.BackgroundColor.A > 0 {
		paintBackground(context, r.style, r.PaddingRect())
	}
}

func (r *BaseRenderBox) PaintBorder(context engine.RenderContext) {
	// Implementation depends on border style
	// Will be implemented when we add border styles
}

func (r *BaseRenderBox) PaintContent(context engine.RenderContext) {
	for _, child := range r.children {
		child.Paint(context)
	}
}

func (r *BaseRenderBox) ContentRect() geometry.Rect {
	insets := r.style.Padding
	return geometry.Rect{
		Min: geometry.Point{X: insets.Left, Y: insets.Top},
		Max: geometry.Point{
			X: r.size.Width - insets.Right,
			Y: r.size.Height - insets.Bottom,
		},
	}
}

func (r *BaseRenderBox) PaddingRect() geometry.Rect {
	return geometry.Rect{
		Min: geometry.Point{X: 0, Y: 0},
		Max: geometry.Point{X: r.size.Width, Y: r.size.Height},
	}
}

func (r *BaseRenderBox) BorderRect() geometry.Rect {
	insets := r.style.Margin
	return geometry.Rect{
		Min: geometry.Point{X: -insets.Left, Y: -insets.Top},
		Max: geometry.Point{
			X: r.size.Width + insets.Right,
			Y: r.size.Height + insets.Bottom,
		},
	}
}

func (r *BaseRenderBox) MarginRect() geometry.Rect {
	border := r.BorderRect()
	insets := r.style.Margin
	return geometry.Rect{
		Min: geometry.Point{
			X: border.Min.X - insets.Left,
			Y: border.Min.Y - insets.Top,
		},
		Max: geometry.Point{
			X: border.Max.X + insets.Right,
			Y: border.Max.Y + insets.Bottom,
		},
	}
}
