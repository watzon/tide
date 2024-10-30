// Copyright (c) 2024 Chris Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package capabilities

// Capabilities describes what the rendering backend can do
type Capabilities struct {
	// Color support
	ColorMode ColorMode

	// Text styling support
	SupportsItalic        bool
	SupportsBold          bool
	SupportsUnderline     bool
	SupportsStrikethrough bool

	// Input capabilities
	SupportsMouse    bool
	SupportsKeyboard bool
}

// ColorMode represents different levels of color support
type ColorMode int

const (
	ColorNone ColorMode = iota
	Color16
	Color256
	ColorTrueColor
)
