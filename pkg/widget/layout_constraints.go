package widget

import (
	"math"

	"github.com/watzon/tide/pkg/core/geometry"
)

// AxisConstraints represents constraints along a single axis
type AxisConstraints struct {
	Min int
	Max int
}

// BoxConstraints represents constraints for box-model layouts
type BoxConstraints struct {
	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int
}

// ToConstraints converts BoxConstraints to regular Constraints
func (bc BoxConstraints) ToConstraints() Constraints {
	return Constraints{
		MinSize: geometry.Size{
			Width:  bc.MinWidth,
			Height: bc.MinHeight,
		},
		MaxSize: geometry.Size{
			Width:  bc.MaxWidth,
			Height: bc.MaxHeight,
		},
	}
}

// LooseConstraints creates constraints that must be at least a given size
func LooseConstraints(minSize geometry.Size) Constraints {
	return Constraints{
		MinSize: minSize,
		MaxSize: geometry.Size{
			Width:  math.MaxInt32,
			Height: math.MaxInt32,
		},
	}
}

// TightConstraints creates constraints for an exact size
func TightConstraints(size geometry.Size) Constraints {
	return Constraints{
		MinSize: size,
		MaxSize: size,
	}
}

// ExpandedConstraints creates constraints that fill available space
func ExpandedConstraints(maxSize geometry.Size) Constraints {
	return Constraints{
		MinSize: maxSize,
		MaxSize: maxSize,
	}
}
