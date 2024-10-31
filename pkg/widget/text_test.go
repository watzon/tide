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
)

func TestText_Paint(t *testing.T) {
	ctx := NewMockRenderContext()
	text := NewText("Hello")

	// Set constraints to establish size
	constraints := NewConstraints(
		geometry.Size{Width: 5, Height: 1},
		geometry.Size{Width: 5, Height: 1},
	)
	text.WithConstraints(constraints)

	// Set style
	style := NewWidgetStyle().
		WithForeground(color.White).
		WithBackground(color.Blue)
	text.WithStyle(style)

	// Create and paint render object
	renderObj := text.CreateRenderObject()
	renderObj.Layout(constraints) // Need to layout before painting
	renderObj.Paint(ctx)

	// Verify text was painted
	for i, ch := range "Hello" {
		cell := ctx.cells[geometry.Point{X: i, Y: 0}]
		assert.Equal(t, ch, cell.Rune)
		assert.Equal(t, color.White, cell.Fg)
		assert.Equal(t, color.Blue, cell.Bg)
	}
}
