// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package render

import (
	"github.com/watzon/tide/pkg/core"
)

// Cell represents a single character cell in the buffer
type Cell struct {
	Rune rune
	Fg   core.Color
	Bg   core.Color
}

// Buffer provides a drawing surface that can be rendered to a backend
type Buffer struct {
	cells [][]Cell
	size  core.Size
}

func NewBuffer(size core.Size) *Buffer {
	cells := make([][]Cell, size.Height)
	for i := range cells {
		cells[i] = make([]Cell, size.Width)
	}

	return &Buffer{
		cells: cells,
		size:  size,
	}
}

func (b *Buffer) GetCell(x, y int) Cell {
	if x < 0 || x >= b.size.Width || y < 0 || y >= b.size.Height {
		return Cell{}
	}

	return b.cells[y][x]
}

func (b *Buffer) SetCell(x, y int, ch rune, fg, bg core.Color) {
	if x < 0 || x >= b.size.Width || y < 0 || y >= b.size.Height {
		return
	}

	b.cells[y][x] = Cell{
		Rune: ch,
		Fg:   fg,
		Bg:   bg,
	}
}
