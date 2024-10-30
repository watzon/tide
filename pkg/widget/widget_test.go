package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/watzon/tide/pkg/core/color"
	"github.com/watzon/tide/pkg/core/geometry"
)

// StringKey implements Key interface for testing
type StringKey string

func (k StringKey) String() string {
	return string(k)
}

// MockBuildContext implements BuildContext for testing
type MockBuildContext struct {
	BuildContext
}

func TestBaseWidget_Identity(t *testing.T) {
	key := StringKey("test-key")
	widget := &BaseWidget{
		key: key,
	}

	t.Run("GetKey", func(t *testing.T) {
		assert.Equal(t, key, widget.GetKey())
	})

	t.Run("GetType", func(t *testing.T) {
		expectedType := "*widget.BaseWidget"
		assert.Equal(t, expectedType, widget.GetType())
	})
}

func TestBaseWidget_Layout(t *testing.T) {
	constraints := NewConstraints(
		geometry.Size{Width: 100, Height: 100},
		geometry.Size{Width: 200, Height: 200},
	)
	size := geometry.Size{Width: 150, Height: 150}

	widget := &BaseWidget{
		constraints: constraints,
		size:        size,
	}

	t.Run("GetConstraints", func(t *testing.T) {
		assert.Equal(t, constraints, widget.GetConstraints())
	})

	t.Run("GetSize", func(t *testing.T) {
		assert.Equal(t, size, widget.GetSize())
	})
}

func TestBaseWidget_Style(t *testing.T) {
	widget := &BaseWidget{}
	style := NewWidgetStyle().
		WithForeground(color.Red).
		WithBackground(color.Blue)

	t.Run("GetStyle", func(t *testing.T) {
		widget.style = style
		assert.Equal(t, style, widget.GetStyle())
	})

	t.Run("WithStyle", func(t *testing.T) {
		result := widget.WithStyle(style)
		assert.Equal(t, style, result.style)
		assert.Same(t, widget, result, "WithStyle should return the same widget instance")
	})
}

func TestBaseWidget_Build(t *testing.T) {
	widget := &BaseWidget{}
	ctx := &MockBuildContext{}

	built := widget.Build(ctx)
	assert.Same(t, widget, built, "Base widget should return itself from Build")
}

func TestBaseWidget_RenderObject(t *testing.T) {
	widget := &BaseWidget{
		style: NewWidgetStyle().WithForeground(color.Red),
	}

	t.Run("CreateRenderObject", func(t *testing.T) {
		renderObj := widget.CreateRenderObject()
		baseRenderObj, ok := renderObj.(*BaseRenderObject)
		assert.True(t, ok, "Should create a BaseRenderObject")
		assert.Equal(t, widget.style, baseRenderObj.style)
	})

	t.Run("UpdateRenderObject", func(t *testing.T) {
		renderObj := &BaseRenderObject{
			style: NewWidgetStyle(),
		}
		widget.UpdateRenderObject(renderObj)
		assert.Equal(t, widget.style, renderObj.style)

		// Test with non-BaseRenderObject (should not panic)
		widget.UpdateRenderObject(&MockRenderObject{})
	})
}

func TestBaseWidget_Fluent(t *testing.T) {
	widget := &BaseWidget{}
	key := StringKey("test-key")
	constraints := NewConstraints(
		geometry.Size{Width: 100, Height: 100},
		geometry.Size{Width: 200, Height: 200},
	)

	t.Run("WithKey", func(t *testing.T) {
		result := widget.WithKey(key)
		assert.Equal(t, key, result.key)
		assert.Same(t, widget, result)
	})

	t.Run("WithConstraints", func(t *testing.T) {
		result := widget.WithConstraints(constraints)
		assert.Equal(t, constraints, result.constraints)
		assert.Same(t, widget, result)
	})
}

func TestMockWidget(t *testing.T) {
	size := geometry.Size{Width: 100, Height: 100}
	constraints := NewConstraints(
		geometry.Size{Width: 50, Height: 50},
		geometry.Size{Width: 150, Height: 150},
	)
	buildResult := &BaseWidget{}
	ctx := &MockBuildContext{}

	t.Run("Default behavior", func(t *testing.T) {
		widget := &MockWidget{
			size:        size,
			constraints: constraints,
		}

		assert.Equal(t, size, widget.GetSize())
		assert.Equal(t, constraints, widget.GetConstraints())
		assert.Same(t, widget, widget.Build(ctx))
	})

	t.Run("With build result", func(t *testing.T) {
		widget := &MockWidget{
			buildResult: buildResult,
			shouldBuild: true,
		}

		assert.Same(t, buildResult, widget.Build(ctx))
	})

	t.Run("With shouldBuild false", func(t *testing.T) {
		widget := &MockWidget{
			buildResult: buildResult,
			shouldBuild: false,
		}

		assert.Same(t, widget, widget.Build(ctx))
	})
}

// MockRenderObject for testing UpdateRenderObject with non-BaseRenderObject
type MockRenderObject struct {
	RenderObject
}
