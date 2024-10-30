// Copyright (c) 2024 Chris Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package widget

import "github.com/watzon/tide/internal/utils"

// EdgeInsets represents spacing measurements for all four edges
type EdgeInsets struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// Common EdgeInsets configurations
var (
	// EdgeInsetsZero represents zero insets on all sides
	EdgeInsetsZero = EdgeInsets{}

	// EdgeInsetsAll creates equal insets on all sides
	EdgeInsetsAll = func(value int) EdgeInsets {
		return EdgeInsets{
			Top:    value,
			Right:  value,
			Bottom: value,
			Left:   value,
		}
	}

	// EdgeInsetsSymmetric creates symmetric horizontal and vertical insets
	EdgeInsetsSymmetric = func(vertical, horizontal int) EdgeInsets {
		return EdgeInsets{
			Top:    vertical,
			Right:  horizontal,
			Bottom: vertical,
			Left:   horizontal,
		}
	}
)

// Constructor functions

// NewEdgeInsets creates EdgeInsets with specific values for each side
func NewEdgeInsets(top, right, bottom, left int) EdgeInsets {
	return EdgeInsets{
		Top:    top,
		Right:  right,
		Bottom: bottom,
		Left:   left,
	}
}

// Helper methods

// Horizontal returns the total horizontal insets (left + right)
func (e EdgeInsets) Horizontal() int {
	return e.Left + e.Right
}

// Vertical returns the total vertical insets (top + bottom)
func (e EdgeInsets) Vertical() int {
	return e.Top + e.Bottom
}

// Add combines two EdgeInsets
func (e EdgeInsets) Add(other EdgeInsets) EdgeInsets {
	return EdgeInsets{
		Top:    e.Top + other.Top,
		Right:  e.Right + other.Right,
		Bottom: e.Bottom + other.Bottom,
		Left:   e.Left + other.Left,
	}
}

// Scale multiplies all insets by a factor
func (e EdgeInsets) Scale(factor int) EdgeInsets {
	return EdgeInsets{
		Top:    e.Top * factor,
		Right:  e.Right * factor,
		Bottom: e.Bottom * factor,
		Left:   e.Left * factor,
	}
}

// IsZero returns true if all insets are zero
func (e EdgeInsets) IsZero() bool {
	return e.Top == 0 && e.Right == 0 && e.Bottom == 0 && e.Left == 0
}

// Max returns EdgeInsets with the maximum value between two EdgeInsets for each side
func (e EdgeInsets) Max(other EdgeInsets) EdgeInsets {
	return EdgeInsets{
		Top:    utils.Max(e.Top, other.Top),
		Right:  utils.Max(e.Right, other.Right),
		Bottom: utils.Max(e.Bottom, other.Bottom),
		Left:   utils.Max(e.Left, other.Left),
	}
}

// Min returns EdgeInsets with the minimum value between two EdgeInsets for each side
func (e EdgeInsets) Min(other EdgeInsets) EdgeInsets {
	return EdgeInsets{
		Top:    utils.Min(e.Top, other.Top),
		Right:  utils.Min(e.Right, other.Right),
		Bottom: utils.Min(e.Bottom, other.Bottom),
		Left:   utils.Min(e.Left, other.Left),
	}
}

// Clamp ensures all insets are within a range
func (e EdgeInsets) Clamp(min, max int) EdgeInsets {
	return EdgeInsets{
		Top:    utils.Clamp(e.Top, min, max),
		Right:  utils.Clamp(e.Right, min, max),
		Bottom: utils.Clamp(e.Bottom, min, max),
		Left:   utils.Clamp(e.Left, min, max),
	}
}
