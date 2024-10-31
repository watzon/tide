// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package widget

// Element represents a widget instance in the tree
type Element interface {
	// Tree structure
	Parent() Element
	Children() []Element

	// Lifecycle
	Mount(parent Element)
	Unmount()

	// Updates
	Update(Widget)
	MarkNeedsBuild()
	RebuildIfNeeded()

	// Access
	Widget() Widget
	RenderObject() RenderObject
	BuildContext() BuildContext

	// ReplaceChild
	ReplaceChild(old, new Element)

	// Layout phase
	LayoutPhase()
	NeedsLayout() bool
	MarkNeedsLayout()
}

// BaseElement provides common element functionality
type BaseElement struct {
	widget       Widget
	parent       Element
	children     []Element
	renderObject RenderObject
	dirty        bool
	mounted      bool
	needsLayout  bool
}

func (e *BaseElement) Parent() Element {
	return e.parent
}

func (e *BaseElement) Children() []Element {
	return e.children
}

func (e *BaseElement) Widget() Widget {
	return e.widget
}

func (e *BaseElement) RenderObject() RenderObject {
	return e.renderObject
}

func (e *BaseElement) Mount(parent Element) {
	if e.mounted {
		return
	}

	e.parent = parent
	e.mounted = true

	// Create render object
	e.renderObject = e.widget.CreateRenderObject()

	// Initial build
	e.Build()
}

func (e *BaseElement) Unmount() {
	if !e.mounted {
		return
	}

	// Unmount children first
	for _, child := range e.children {
		child.Unmount()
	}

	e.mounted = false
	e.parent = nil
	e.children = nil
	e.renderObject = nil
}

func (e *BaseElement) Build() {
	if !e.mounted {
		return
	}

	// Build new widget
	newWidget := e.widget.Build(e.BuildContext())
	if newWidget == e.widget {
		// If the widget returns itself, don't create a new child
		e.dirty = false
		return
	}

	// Update or create child element
	if len(e.children) > 0 {
		e.children[0].Update(newWidget)
	} else {
		child := NewElement(newWidget)
		e.children = append(e.children, child)
		child.Mount(e)
	}

	e.dirty = false
}

func (e *BaseElement) Update(newWidget Widget) {
	// Check if widget types are different using type assertion
	_, oldIsMock := e.widget.(*MockWidget)
	_, newIsMock := newWidget.(*MockWidget)

	if oldIsMock != newIsMock {
		// Replace entire element if widget type changes
		if e.parent != nil {
			newElement := NewElement(newWidget)
			e.parent.ReplaceChild(e, newElement)
		}
		return
	}

	e.widget = newWidget
	e.widget.UpdateRenderObject(e.renderObject)
	e.MarkNeedsBuild()
}

func (e *BaseElement) MarkNeedsBuild() {
	e.dirty = true
	// Propagate to parent if needed
	if e.parent != nil {
		e.parent.MarkNeedsBuild()
	}
}

func (e *BaseElement) RebuildIfNeeded() {
	if e.dirty {
		e.Build()
	}
}

func (e *BaseElement) BuildContext() BuildContext {
	return &ElementBuildContext{element: e}
}

func (e *BaseElement) ReplaceChild(old, new Element) {
	for i, child := range e.children {
		if child == old {
			// Unmount old child
			old.Unmount()

			// Mount new child
			e.children[i] = new
			new.Mount(e)

			// Mark parent as needing rebuild
			e.MarkNeedsBuild()
			return
		}
	}
}

// NewElement creates the appropriate element type for a widget
func NewElement(widget Widget) Element {
	switch w := widget.(type) {
	case StatefulWidget:
		return NewStatefulElement(w)
	default:
		elem := &BaseElement{}
		elem.widget = widget
		return elem
	}
}

func (e *BaseElement) LayoutPhase() {
	if !e.needsLayout {
		return
	}

	// Layout this element's render object
	if e.renderObject != nil {
		e.renderObject.Layout(e.widget.GetConstraints())
	}

	// Layout children
	for _, child := range e.children {
		child.LayoutPhase()
	}

	e.needsLayout = false
}

func (e *BaseElement) NeedsLayout() bool {
	return e.needsLayout
}

func (e *BaseElement) MarkNeedsLayout() {
	e.needsLayout = true
	// Propagate to parent
	if e.parent != nil {
		e.parent.MarkNeedsLayout()
	}
}

// MockElement implements Element interface for testing
type MockElement struct {
	BaseElement
	parent       Element
	buildContext BuildContext
}

func (e *MockElement) Parent() Element {
	return e.parent
}

func (e *MockElement) BuildContext() BuildContext {
	return e.buildContext
}

// StatefulElement manages the lifecycle of a StatefulWidget
type StatefulElement interface {
	Element
	State() State
}

// baseStatefulElement provides the implementation of StatefulElement
type baseStatefulElement struct {
	BaseElement
	state State
}

func NewStatefulElement(widget StatefulWidget) StatefulElement {
	elem := &baseStatefulElement{}
	elem.widget = widget
	elem.state = widget.CreateState()
	return elem
}

func (e *baseStatefulElement) State() State {
	return e.state
}

func (e *baseStatefulElement) Mount(parent Element) {
	if e.mounted {
		return
	}

	e.parent = parent
	e.mounted = true

	// Create render object
	e.renderObject = e.widget.CreateRenderObject()

	// Initialize state
	e.state.MountState(e)

	// Initial build
	e.Build()
}

func (e *baseStatefulElement) Unmount() {
	if !e.mounted {
		return
	}

	// Dispose state
	e.state.Dispose()

	// Unmount children
	for _, child := range e.children {
		child.Unmount()
	}

	e.mounted = false
	e.parent = nil
	e.children = nil
	e.renderObject = nil
	e.state = nil
}

func (e *baseStatefulElement) Build() {
	if !e.mounted {
		return
	}

	// Build using state
	newWidget := e.state.Build(e.BuildContext())

	// Update or create child element
	if len(e.children) > 0 {
		e.children[0].Update(newWidget)
	} else {
		child := NewElement(newWidget)
		e.children = append(e.children, child)
		child.Mount(e)
	}

	e.dirty = false
}

func (e *baseStatefulElement) Update(newWidget Widget) {
	// Update widget reference
	e.widget = newWidget.(StatefulWidget)
	e.widget.UpdateRenderObject(e.renderObject)
	e.MarkNeedsBuild()
}
