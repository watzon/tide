// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package render_test

import (
	"testing"

	"github.com/watzon/tide/internal/render"
	"github.com/watzon/tide/pkg/core"
)

func TestBuffer(t *testing.T) {
	t.Run("NewBuffer", func(t *testing.T) {
		size := core.Size{Width: 80, Height: 24}
		buffer := render.NewBuffer(size)

		if buffer == nil {
			t.Error("expected non-nil buffer")
		}
	})

	t.Run("SetCell within bounds", func(t *testing.T) {
		size := core.Size{Width: 80, Height: 24}
		buffer := render.NewBuffer(size)

		fg := core.Color{R: 255, G: 255, B: 255, A: 255}
		bg := core.Color{R: 0, G: 0, B: 0, A: 255}

		buffer.SetCell(10, 10, 'A', fg, bg)

		// Add a method to Buffer to get cell for testing
		cell := buffer.GetCell(10, 10)
		if cell.Rune != 'A' {
			t.Errorf("expected rune 'A', got %c", cell.Rune)
		}
	})

	t.Run("SetCell out of bounds", func(t *testing.T) {
		size := core.Size{Width: 80, Height: 24}
		buffer := render.NewBuffer(size)

		fg := core.Color{R: 255, G: 255, B: 255, A: 255}
		bg := core.Color{R: 0, G: 0, B: 0, A: 255}

		// These should not panic
		buffer.SetCell(-1, 10, 'A', fg, bg)
		buffer.SetCell(80, 10, 'A', fg, bg)
		buffer.SetCell(10, -1, 'A', fg, bg)
		buffer.SetCell(10, 24, 'A', fg, bg)
	})
}
