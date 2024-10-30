package widget

import (
	"fmt"
	"math"

	"github.com/watzon/tide/internal/utils"
	"github.com/watzon/tide/pkg/core/geometry"
)

// Constraints defines the rules for laying out widgets
type Constraints struct {
	// MinSize defines the minimum allowed size
	MinSize geometry.Size
	// MaxSize defines the maximum allowed size
	MaxSize geometry.Size
}

// Common constraint configurations
var (
	// ConstraintsUnbounded represents constraints with no limits
	ConstraintsUnbounded = Constraints{
		MinSize: geometry.Size{Width: 0, Height: 0},
		MaxSize: geometry.Size{
			Width:  math.MaxInt32,
			Height: math.MaxInt32,
		},
	}

	// ConstraintsTight forces a specific size
	ConstraintsTight = func(size geometry.Size) Constraints {
		return Constraints{
			MinSize: size,
			MaxSize: size,
		}
	}
)

// NewConstraints creates constraints with the given bounds
func NewConstraints(minSize, maxSize geometry.Size) Constraints {
	return Constraints{
		MinSize: minSize,
		MaxSize: maxSize,
	}
}

// WithMinSize returns new constraints with the given minimum size
func (c Constraints) WithMinSize(size geometry.Size) Constraints {
	c.MinSize = size
	return c.Normalize()
}

// WithMaxSize returns new constraints with the given maximum size
func (c Constraints) WithMaxSize(size geometry.Size) Constraints {
	c.MaxSize = size
	return c.Normalize()
}

// Normalize ensures constraints are valid
func (c Constraints) Normalize() Constraints {
	// Ensure min <= max
	return Constraints{
		MinSize: geometry.Size{
			Width:  min(c.MinSize.Width, c.MaxSize.Width),
			Height: min(c.MinSize.Height, c.MaxSize.Height),
		},
		MaxSize: geometry.Size{
			Width:  max(c.MinSize.Width, c.MaxSize.Width),
			Height: max(c.MinSize.Height, c.MaxSize.Height),
		},
	}
}

// Constrain forces a size to fit within the constraints
func (c Constraints) Constrain(size geometry.Size) geometry.Size {
	return geometry.Size{
		Width:  utils.ClampInt(size.Width, c.MinSize.Width, c.MaxSize.Width),
		Height: utils.ClampInt(size.Height, c.MinSize.Height, c.MaxSize.Height),
	}
}

// IsSatisfiedBy checks if a size satisfies the constraints
func (c Constraints) IsSatisfiedBy(size geometry.Size) bool {
	return size.Width >= c.MinSize.Width &&
		size.Width <= c.MaxSize.Width &&
		size.Height >= c.MinSize.Height &&
		size.Height <= c.MaxSize.Height
}

// HasTightWidth returns true if the width is tightly constrained
func (c Constraints) HasTightWidth() bool {
	return c.MinSize.Width == c.MaxSize.Width
}

// HasTightHeight returns true if the height is tightly constrained
func (c Constraints) HasTightHeight() bool {
	return c.MinSize.Height == c.MaxSize.Height
}

// IsTight returns true if both dimensions are tightly constrained
func (c Constraints) IsTight() bool {
	return c.HasTightWidth() && c.HasTightHeight()
}

// IsNormalized returns true if the constraints are valid
func (c Constraints) IsNormalized() bool {
	return c.MinSize.Width <= c.MaxSize.Width &&
		c.MinSize.Height <= c.MaxSize.Height
}

// String provides a readable representation of constraints
func (c Constraints) String() string {
	return fmt.Sprintf("Constraints(min: %v, max: %v)", c.MinSize, c.MaxSize)
}
