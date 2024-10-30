// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseElement_Parent(t *testing.T) {
	parent := &BaseElement{
		widget: &MockWidget{},
	}
	child := &BaseElement{
		widget: &MockWidget{},
	}
	child.Mount(parent)

	assert.Equal(t, parent, child.Parent())
}

func TestBaseElement_Children(t *testing.T) {
	parent := &BaseElement{
		widget: &MockWidget{},
	}
	child := &BaseElement{
		widget: &MockWidget{},
	}
	child.Mount(parent)
	parent.children = append(parent.children, child)

	assert.Equal(t, []Element{child}, parent.Children())
}

func TestBaseElement_Mount(t *testing.T) {
	parent := &BaseElement{}
	child := &BaseElement{
		widget: &MockWidget{},
	}

	// Test initial mount
	child.Mount(parent)
	assert.True(t, child.mounted)
	assert.Equal(t, parent, child.parent)
	assert.NotNil(t, child.renderObject)

	// Test mounting an already mounted element
	child.Mount(parent)
	assert.True(t, child.mounted)
}

func TestBaseElement_Unmount(t *testing.T) {
	parent := &BaseElement{
		widget: &MockWidget{},
	}
	child := &BaseElement{
		widget: &MockWidget{},
	}
	grandChild := &BaseElement{
		widget: &MockWidget{},
	}

	// Setup hierarchy
	child.Mount(parent)
	grandChild.Mount(child)
	child.children = append(child.children, grandChild)

	// Test unmount
	child.Unmount()
	assert.False(t, child.mounted)
	assert.Nil(t, child.parent)
	assert.Empty(t, child.children)
	assert.Nil(t, child.renderObject)

	// Test unmounting an already unmounted element
	child.Unmount()
	assert.False(t, child.mounted)
}

func TestBaseElement_Build(t *testing.T) {
	parent := &BaseElement{
		widget: &MockWidget{},
	}
	child := &BaseElement{
		widget: &MockWidget{
			shouldBuild: true,
			buildResult: &MockWidget{},
		},
	}

	// Test build before mount
	child.Build()
	assert.Empty(t, child.children)

	// Test build after mount
	child.Mount(parent)
	child.Build()
	assert.NotEmpty(t, child.children)
	assert.False(t, child.dirty)

	// Test build with existing children
	existingChild := child.children[0]
	child.Build()
	assert.Equal(t, 1, len(child.children))
	assert.Equal(t, existingChild, child.children[0])
}

func TestBaseElement_Update(t *testing.T) {
	parent := &BaseElement{}
	child := &BaseElement{
		widget: &MockWidget{},
	}
	child.Mount(parent)

	// Test update with same widget type
	newWidget := &MockWidget{}
	child.Update(newWidget)
	assert.Equal(t, newWidget, child.widget)
	assert.True(t, child.dirty)

	// Test update with different widget type
	differentWidget := &BaseWidget{}
	child.Update(differentWidget)
	assert.NotEqual(t, differentWidget, child.widget)
}

func TestBaseElement_MarkNeedsBuild(t *testing.T) {
	parent := &BaseElement{
		widget: &MockWidget{},
	}
	child := &BaseElement{
		widget: &MockWidget{},
	}
	child.Mount(parent)

	child.MarkNeedsBuild()
	assert.True(t, child.dirty)
	assert.True(t, parent.dirty)
}

func TestBaseElement_RebuildIfNeeded(t *testing.T) {
	element := &BaseElement{
		widget: &MockWidget{},
		dirty:  true,
	}
	element.mounted = true

	element.RebuildIfNeeded()
	assert.False(t, element.dirty)
}

func TestBaseElement_ReplaceChild(t *testing.T) {
	parent := &BaseElement{
		widget: &MockWidget{},
	}
	oldChild := &BaseElement{
		widget: &MockWidget{},
	}
	newChild := &BaseElement{
		widget: &MockWidget{},
	}

	parent.children = append(parent.children, oldChild)
	oldChild.Mount(parent)

	parent.ReplaceChild(oldChild, newChild)
	assert.Equal(t, newChild, parent.children[0])
	assert.True(t, parent.dirty)
}

func TestBaseElement_LayoutPhase(t *testing.T) {
	element := &BaseElement{
		widget:      &MockWidget{},
		needsLayout: true,
	}
	element.mounted = true
	element.renderObject = &BaseRenderObject{}

	element.LayoutPhase()
	assert.False(t, element.needsLayout)
}

func TestBaseElement_MarkNeedsLayout(t *testing.T) {
	parent := &BaseElement{
		widget: &MockWidget{},
	}
	child := &BaseElement{
		widget: &MockWidget{},
	}
	child.Mount(parent)

	child.MarkNeedsLayout()
	assert.True(t, child.needsLayout)
	assert.True(t, parent.needsLayout)
}

func TestNewElement(t *testing.T) {
	widget := &MockWidget{}
	element := NewElement(widget)

	assert.NotNil(t, element)
	assert.Equal(t, widget, element.Widget())
}
