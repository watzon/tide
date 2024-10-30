// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package widget

import (
	"github.com/watzon/tide/pkg/core/geometry"
	"github.com/watzon/tide/pkg/engine"
)

// BuildContext provides context for widget building
type BuildContext interface {
	// Tree traversal
	Parent() BuildContext

	// Widget information
	Widget() Widget
	Element() Element

	// Layout information
	Size() geometry.Size
	Constraints() Constraints

	// Invalidation
	MarkNeedsBuild()

	// Rendering
	RenderContext() engine.RenderContext
}

// ElementBuildContext implements BuildContext for Elements
type ElementBuildContext struct {
	element Element
}

func NewElementBuildContext(element Element) BuildContext {
	return &ElementBuildContext{
		element: element,
	}
}

// Tree traversal
func (c *ElementBuildContext) Parent() BuildContext {
	if parent := c.element.Parent(); parent != nil {
		return parent.BuildContext()
	}
	return nil
}

// Widget information
func (c *ElementBuildContext) Widget() Widget {
	return c.element.Widget()
}

func (c *ElementBuildContext) Element() Element {
	return c.element
}

// Layout information
func (c *ElementBuildContext) Size() geometry.Size {
	return c.Widget().GetSize()
}

func (c *ElementBuildContext) Constraints() Constraints {
	return c.Widget().GetConstraints()
}

// Invalidation
func (c *ElementBuildContext) MarkNeedsBuild() {
	c.element.MarkNeedsBuild()
}

// Rendering
func (c *ElementBuildContext) RenderContext() engine.RenderContext {
	// Walk up the tree to find the nearest RenderContext
	current := c.element
	for current != nil {
		if ctx, ok := current.(RenderContextProvider); ok {
			return ctx.GetRenderContext()
		}
		current = current.Parent()
	}
	return nil
}

// RenderContextProvider allows elements to provide render context
type RenderContextProvider interface {
	GetRenderContext() engine.RenderContext
}
