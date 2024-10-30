package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/watzon/tide/pkg/core/geometry"
	"github.com/watzon/tide/pkg/engine"
)

// MockRenderContextProvider implements RenderContextProvider for testing
type MockRenderContextProvider struct {
	MockElement
	renderContext engine.RenderContext
}

func (p *MockRenderContextProvider) GetRenderContext() engine.RenderContext {
	return p.renderContext
}

func TestElementBuildContext_Parent(t *testing.T) {
	// Test with no parent
	widget := &MockWidget{}
	element := &MockElement{}
	element.widget = widget
	ctx := NewElementBuildContext(element)

	assert.Nil(t, ctx.Parent())

	// Test with parent
	parentWidget := &MockWidget{}
	parentElement := &MockElement{}
	parentElement.widget = parentWidget
	element.parent = parentElement
	parentElement.buildContext = NewElementBuildContext(parentElement)

	assert.NotNil(t, ctx.Parent())
	assert.Equal(t, parentElement.BuildContext(), ctx.Parent())
}

func TestElementBuildContext_Widget(t *testing.T) {
	widget := &MockWidget{}
	element := &MockElement{}
	element.widget = widget
	ctx := NewElementBuildContext(element)

	assert.Equal(t, widget, ctx.Widget())
}

func TestElementBuildContext_Element(t *testing.T) {
	widget := &MockWidget{}
	element := &MockElement{}
	element.widget = widget
	ctx := NewElementBuildContext(element)

	assert.Equal(t, element, ctx.Element())
}

func TestElementBuildContext_Size(t *testing.T) {
	expectedSize := geometry.Size{Width: 100, Height: 200}
	widget := &MockWidget{size: expectedSize}
	element := &MockElement{}
	element.widget = widget
	ctx := NewElementBuildContext(element)

	assert.Equal(t, expectedSize, ctx.Size())
}

func TestElementBuildContext_Constraints(t *testing.T) {
	expectedConstraints := Constraints{
		MinSize: geometry.Size{Width: 50, Height: 50},
		MaxSize: geometry.Size{Width: 200, Height: 200},
	}
	widget := &MockWidget{constraints: expectedConstraints}
	element := &MockElement{}
	element.widget = widget
	ctx := NewElementBuildContext(element)

	assert.Equal(t, expectedConstraints, ctx.Constraints())
}

func TestElementBuildContext_MarkNeedsBuild(t *testing.T) {
	widget := &MockWidget{}
	element := &MockElement{}
	element.widget = widget
	ctx := NewElementBuildContext(element)

	// Verify initial state
	assert.False(t, element.dirty)

	// Mark needs build
	ctx.MarkNeedsBuild()

	// Verify state changed
	assert.True(t, element.dirty)
}

func TestElementBuildContext_RenderContext(t *testing.T) {
	// Test with no render context provider
	widget := &MockWidget{}
	element := &MockElement{}
	element.widget = widget
	ctx := NewElementBuildContext(element)

	assert.Nil(t, ctx.RenderContext())

	// Test with render context provider
	mockRenderCtx := &engine.MockRenderContext{}
	provider := &MockRenderContextProvider{
		renderContext: mockRenderCtx,
	}
	element.parent = provider

	assert.Equal(t, mockRenderCtx, ctx.RenderContext())
}
