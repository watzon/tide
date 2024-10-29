// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package engine_test

import (
	"testing"

	"github.com/watzon/tide/pkg/core"
	"github.com/watzon/tide/pkg/engine"
)

type mockBackend struct {
	cells [][]rune
	size  core.Size
}

func newMockBackend(size core.Size) *mockBackend {
	cells := make([][]rune, size.Height)
	for i := range cells {
		cells[i] = make([]rune, size.Width)
		// Initialize cells with space character
		for j := range cells[i] {
			cells[i][j] = ' '
		}
	}
	return &mockBackend{cells: cells, size: size}
}

func (m *mockBackend) Init() error     { return nil }
func (m *mockBackend) Shutdown() error { return nil }
func (m *mockBackend) Size() core.Size { return m.size }
func (m *mockBackend) Clear()          {}
func (m *mockBackend) DrawCell(x, y int, ch rune, fg, bg core.Color) {
	if x >= 0 && x < m.size.Width && y >= 0 && y < m.size.Height {
		m.cells[y][x] = ch
	}
}
func (m *mockBackend) Present() error { return nil }

func TestCompositor(t *testing.T) {
	t.Run("Layer ordering", func(t *testing.T) {
		comp := engine.NewCompositor()
		backend := newMockBackend(core.Size{Width: 80, Height: 24})

		// Add layers in reverse Z order
		comp.AddLayer(engine.Layer{
			Bounds: core.NewRect(0, 0, 10, 10),
			Z:      1, // Lower Z-index, drawn first
			Draw: func(b engine.Backend) {
				b.DrawCell(5, 5, 'A', core.Color{}, core.Color{})
			},
		})

		comp.AddLayer(engine.Layer{
			Bounds: core.NewRect(0, 0, 10, 10),
			Z:      2, // Higher Z-index, drawn last
			Draw: func(b engine.Backend) {
				b.DrawCell(5, 5, 'B', core.Color{}, core.Color{})
			},
		})

		comp.Compose(backend)

		if backend.cells[5][5] != 'B' {
			t.Errorf("expected cell to contain 'B', got %c", backend.cells[5][5])
		}
	})

	// Add more test cases
	t.Run("Empty compositor", func(t *testing.T) {
		comp := engine.NewCompositor()
		backend := newMockBackend(core.Size{Width: 80, Height: 24})

		// Should not panic
		comp.Compose(backend)
	})

	t.Run("Multiple layers same Z", func(t *testing.T) {
		comp := engine.NewCompositor()
		backend := newMockBackend(core.Size{Width: 80, Height: 24})

		comp.AddLayer(engine.Layer{
			Bounds: core.NewRect(0, 0, 10, 10),
			Z:      1,
			Draw: func(b engine.Backend) {
				b.DrawCell(5, 5, 'A', core.Color{}, core.Color{})
			},
		})

		comp.AddLayer(engine.Layer{
			Bounds: core.NewRect(0, 0, 10, 10),
			Z:      1,
			Draw: func(b engine.Backend) {
				b.DrawCell(5, 5, 'B', core.Color{}, core.Color{})
			},
		})

		comp.Compose(backend)

		if backend.cells[5][5] != 'B' {
			t.Errorf("expected cell to contain 'B', got %c", backend.cells[5][5])
		}
	})
}
