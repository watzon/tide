// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package engine

import (
	"github.com/watzon/tide/pkg/core"
)

// Backend defines the interface for different rendering backends
type Backend interface {
	// Initialize the backend
	Init() error

	// Clean up resources
	Shutdown() error

	// Get the current size of the rendering surface
	Size() core.Size

	// Clear the entire surface
	Clear()

	// Draw a single cell with the given rune and style
	DrawCell(x, y int, ch rune, fg, bg core.Color)

	// Present the current frame
	Present() error
}
