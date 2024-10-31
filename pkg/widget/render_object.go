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

// Layout implements the box model layout algorithm
func (r *BaseRenderBox) Layout(constraints Constraints) geometry.Size {
	// 1. Calculate available content space by subtracting padding and border
	horizontalInsets := r.style.Padding.Left + r.style.Padding.Right +
		r.style.BorderWidth.Left + r.style.BorderWidth.Right
	verticalInsets := r.style.Padding.Top + r.style.Padding.Bottom +
		r.style.BorderWidth.Top + r.style.BorderWidth.Bottom

	// 2. Create content constraints
	contentConstraints := Constraints{
		MinSize: geometry.Size{
			Width:  max(0, constraints.MinSize.Width-horizontalInsets),
			Height: max(0, constraints.MinSize.Height-verticalInsets),
		},
		MaxSize: geometry.Size{
			Width:  max(0, constraints.MaxSize.Width-horizontalInsets),
			Height: max(0, constraints.MaxSize.Height-verticalInsets),
		},
	}

	// 3. Layout children within content constraints
	contentSize := r.layoutChildren(contentConstraints)

	// 4. Calculate final size including padding and border
	r.size = geometry.Size{
		Width:  contentSize.Width + horizontalInsets,
		Height: contentSize.Height + verticalInsets,
	}

	// 5. Ensure size satisfies original constraints
	r.size = constraints.Constrain(r.size)
	return r.size
}

func (r *BaseRenderBox) layoutChildren(constraints Constraints) geometry.Size {
	if len(r.children) == 0 {
		return constraints.MinSize
	}

	// For now, just layout the first child with full constraints
	// This will be expanded when we add different layout behaviors (row, column, etc.)
	child := r.children[0]
	return child.Layout(constraints)
}

func (r *BaseRenderBox) PaintBackground(context engine.RenderContext) {
	if r.style.BackgroundColor.A > 0 {
		paintBackground(context, r.style, r.PaddingRect())
	}
}

func (r *BaseRenderBox) PaintBorder(context engine.RenderContext) {
	if r.style.BorderWidth.IsZero() {
		return
	}
	// Let the backend handle the border painting
	context.PaintBorder(r.BorderRect(), r.style.Style)
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

// NewBaseRenderBox creates a new BaseRenderBox with default style
func NewBaseRenderBox() *BaseRenderBox {
	return &BaseRenderBox{
		BaseRenderObject: BaseRenderObject{
			style: NewWidgetStyle(),
		},
	}
}

func (r *BaseRenderBox) WithStyle(style WidgetStyle) {
	r.style = style
}
