package terminal_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/watzon/tide/pkg/backend/terminal"
	"github.com/watzon/tide/pkg/core/geometry"
)

func cellsEqual(a, b terminal.Cell) bool {
	if a.Rune != b.Rune || a.Style != b.Style || a.Width != b.Width {
		return false
	}
	if len(a.Combining) != len(b.Combining) {
		return false
	}
	for i := range a.Combining {
		if a.Combining[i] != b.Combining[i] {
			return false
		}
	}
	return true
}

func TestBuffer(t *testing.T) {
	t.Run("NewBuffer", func(t *testing.T) {
		size := geometry.Size{Width: 80, Height: 24}
		buf := terminal.NewBuffer(size)

		if buf == nil {
			t.Fatal("NewBuffer returned nil")
		}

		// Test initial state
		if buf.Size() != size {
			t.Errorf("expected size %v, got %v", size, buf.Size())
		}

		cursor := buf.GetCursor()
		if cursor != (geometry.Point{X: 0, Y: 0}) {
			t.Errorf("expected cursor at origin, got %v", cursor)
		}
	})

	t.Run("SetCell and GetCell", func(t *testing.T) {
		buf := terminal.NewBuffer(geometry.Size{Width: 80, Height: 24})
		style := tcell.StyleDefault.Foreground(tcell.ColorRed)

		// Test setting and getting a cell
		buf.SetCell(5, 5, 'A', nil, style)
		cell, exists := buf.GetCell(5, 5)

		if !exists {
			t.Fatal("cell should exist after SetCell")
		}
		if cell.Rune != 'A' {
			t.Errorf("expected rune 'A', got %c", cell.Rune)
		}
		if cell.Style != style {
			t.Error("cell style mismatch")
		}

		// Test getting non-existent cell
		_, exists = buf.GetCell(100, 100)
		if exists {
			t.Error("GetCell should return false for non-existent cell")
		}
	})

	t.Run("Clear", func(t *testing.T) {
		buf := terminal.NewBuffer(geometry.Size{Width: 80, Height: 24})
		style := tcell.StyleDefault

		// Add some cells
		buf.SetCell(0, 0, 'A', nil, style)
		buf.SetCell(1, 1, 'B', nil, style)

		// Clear the buffer
		buf.Clear()

		// Verify cells are gone
		for y := 0; y < 24; y++ {
			for x := 0; x < 80; x++ {
				if _, exists := buf.GetCell(x, y); exists {
					t.Errorf("cell at (%d,%d) should not exist after Clear", x, y)
				}
			}
		}
	})

	t.Run("Resize", func(t *testing.T) {
		buf := terminal.NewBuffer(geometry.Size{Width: 80, Height: 24})
		style := tcell.StyleDefault

		// Add cells near the edges
		buf.SetCell(79, 23, 'A', nil, style) // Will be preserved
		buf.SetCell(0, 0, 'B', nil, style)   // Will be preserved
		buf.SetCell(85, 25, 'C', nil, style) // Will be removed

		// Resize smaller
		newSize := geometry.Size{Width: 40, Height: 20}
		buf.Resize(newSize)

		// Check size updated
		if buf.Size() != newSize {
			t.Errorf("expected size %v after resize, got %v", newSize, buf.Size())
		}

		// Check cells
		if _, exists := buf.GetCell(0, 0); !exists {
			t.Error("cell at (0,0) should still exist")
		}
		if _, exists := buf.GetCell(85, 25); exists {
			t.Error("cell outside new bounds should not exist")
		}
	})

	t.Run("MoveCursor", func(t *testing.T) {
		size := geometry.Size{Width: 80, Height: 24}
		buf := terminal.NewBuffer(size)

		tests := []struct {
			name     string
			dx, dy   int
			expected geometry.Point
		}{
			{"move right", 5, 0, geometry.Point{X: 5, Y: 0}},
			{"move down", 0, 5, geometry.Point{X: 5, Y: 5}},
			{"move left", -2, 0, geometry.Point{X: 3, Y: 5}},
			{"move up", 0, -2, geometry.Point{X: 3, Y: 3}},
			{"clamp negative", -10, -10, geometry.Point{X: 0, Y: 0}},
			{"clamp positive", 100, 100, geometry.Point{X: 79, Y: 23}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				buf.MoveCursor(tt.dx, tt.dy)
				got := buf.GetCursor()
				if got != tt.expected {
					t.Errorf("MoveCursor(%d,%d) = %v, want %v", tt.dx, tt.dy, got, tt.expected)
				}
			})
		}
	})

	t.Run("CopyFrom", func(t *testing.T) {
		src := terminal.NewBuffer(geometry.Size{Width: 80, Height: 24})
		dst := terminal.NewBuffer(geometry.Size{Width: 40, Height: 20})

		// Setup source buffer
		style := tcell.StyleDefault
		src.SetCell(5, 5, 'A', nil, style)
		src.SetCursor(10, 10)

		// Copy to destination
		dst.CopyFrom(src)

		// Verify size copied
		if dst.Size() != src.Size() {
			t.Errorf("size mismatch after copy: got %v, want %v", dst.Size(), src.Size())
		}

		// Verify cells copied
		srcCell, _ := src.GetCell(5, 5)
		dstCell, exists := dst.GetCell(5, 5)
		if !exists {
			t.Fatal("cell should exist in destination after copy")
		}
		if !cellsEqual(dstCell, srcCell) {
			t.Error("cell content mismatch after copy")
		}

		// Verify cursor position copied
		if dst.GetCursor() != src.GetCursor() {
			t.Errorf("cursor position mismatch: got %v, want %v", dst.GetCursor(), src.GetCursor())
		}

		// Verify modifying source doesn't affect destination
		src.SetCell(5, 5, 'B', nil, style)
		dstCell, _ = dst.GetCell(5, 5)
		if dstCell.Rune != 'A' {
			t.Error("modifying source buffer should not affect destination")
		}
	})
}
