package terminal

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/watzon/tide/pkg/core/geometry"
)

// Cell represents a single character cell in the buffer
type Cell struct {
	Rune      rune
	Style     tcell.Style
	Combining []rune
	Width     int
}

// Buffer represents a screen buffer
type Buffer struct {
	lock   sync.RWMutex
	cells  map[geometry.Point]Cell
	size   geometry.Size
	cursor geometry.Point
	dirty  bool
}

// NewBuffer creates a new buffer with the given size
func NewBuffer(size geometry.Size) *Buffer {
	return &Buffer{
		cells: make(map[geometry.Point]Cell),
		size:  size,
	}
}

// SetCell sets a cell in the buffer
func (b *Buffer) SetCell(x, y int, ch rune, combining []rune, style tcell.Style) {
	b.lock.Lock()
	defer b.lock.Unlock()

	pos := geometry.Point{X: x, Y: y}
	b.cells[pos] = Cell{
		Rune:      ch,
		Style:     style,
		Combining: combining,
		Width:     runewidth.RuneWidth(ch),
	}
	b.dirty = true
}

// GetCell gets a cell from the buffer
func (b *Buffer) GetCell(x, y int) (Cell, bool) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	cell, ok := b.cells[geometry.Point{X: x, Y: y}]
	return cell, ok
}

// Clear clears the buffer
func (b *Buffer) Clear() {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.cells = make(map[geometry.Point]Cell)
	b.dirty = true
}

// Resize resizes the buffer
func (b *Buffer) Resize(size geometry.Size) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.size == size {
		return
	}

	newCells := make(map[geometry.Point]Cell)
	for pos, cell := range b.cells {
		if pos.X < size.Width && pos.Y < size.Height {
			newCells[pos] = cell
		}
	}

	b.cells = newCells
	b.size = size
	b.dirty = true
}

// SetCursor sets the cursor position
func (b *Buffer) SetCursor(x, y int) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.cursor = geometry.Point{X: x, Y: y}
	b.dirty = true
}

// GetCursor returns the current cursor position
func (b *Buffer) GetCursor() geometry.Point {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.cursor
}

// MoveCursor moves the cursor relative to its current position
func (b *Buffer) MoveCursor(dx, dy int) {
	b.cursor.X += dx
	b.cursor.Y += dy

	// Clamp cursor to buffer bounds
	if b.cursor.X < 0 {
		b.cursor.X = 0
	}
	if b.cursor.Y < 0 {
		b.cursor.Y = 0
	}
	if b.cursor.X >= b.size.Width {
		b.cursor.X = b.size.Width - 1
	}
	if b.cursor.Y >= b.size.Height {
		b.cursor.Y = b.size.Height - 1
	}

	b.dirty = true
}

// CopyFrom copies the contents of another buffer
func (b *Buffer) CopyFrom(other *Buffer) {
	b.lock.Lock()
	other.lock.RLock()
	defer b.lock.Unlock()
	defer other.lock.RUnlock()

	// Create new map to avoid sharing underlying data
	b.cells = make(map[geometry.Point]Cell, len(other.cells))
	for pos, cell := range other.cells {
		b.cells[pos] = cell
	}

	b.size = other.size
	b.cursor = other.cursor
	b.dirty = other.dirty
}
