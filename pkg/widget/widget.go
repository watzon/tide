package widget

import (
	"fmt"

	"github.com/watzon/tide/pkg/core/geometry"
)

// Key uniquely identifies a widget in the tree
type Key interface {
	String() string
}

// Widget is the base interface for all widgets
type Widget interface {
	// Core identity
	GetKey() Key
	GetType() string

	// Layout
	GetConstraints() Constraints
	GetSize() geometry.Size

	// Building
	Build(context BuildContext) Widget

	// Rendering
	CreateRenderObject() RenderObject
	UpdateRenderObject(RenderObject)
}

// BaseWidget provides common functionality for all widgets
type BaseWidget struct {
	key         Key
	constraints Constraints
	size        geometry.Size
	style       WidgetStyle
}

// Identity methods
func (w *BaseWidget) GetKey() Key {
	return w.key
}

func (w *BaseWidget) GetType() string {
	return fmt.Sprintf("%T", w)
}

// Layout methods
func (w *BaseWidget) GetConstraints() Constraints {
	return w.constraints
}

func (w *BaseWidget) GetSize() geometry.Size {
	return w.size
}

// Style methods
func (w *BaseWidget) GetStyle() WidgetStyle {
	return w.style
}

func (w *BaseWidget) WithStyle(style WidgetStyle) *BaseWidget {
	w.style = style
	return w
}

// Builder methods - these should be overridden by implementing widgets
func (w *BaseWidget) Build(context BuildContext) Widget {
	return w // Base widgets are leaves by default
}

func (w *BaseWidget) CreateRenderObject() RenderObject {
	return &BaseRenderObject{
		style: w.style,
	}
}

func (w *BaseWidget) UpdateRenderObject(renderObject RenderObject) {
	if baseRenderObject, ok := renderObject.(*BaseRenderObject); ok {
		baseRenderObject.style = w.style
	}
}

// WithKey sets a key for the widget
func (w *BaseWidget) WithKey(key Key) *BaseWidget {
	w.key = key
	return w
}

// WithConstraints sets constraints for the widget
func (w *BaseWidget) WithConstraints(constraints Constraints) *BaseWidget {
	w.constraints = constraints
	return w
}
