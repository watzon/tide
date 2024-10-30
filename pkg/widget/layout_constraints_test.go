// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package widget

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/watzon/tide/pkg/core/geometry"
)

func TestAxisConstraints(t *testing.T) {
	constraints := AxisConstraints{
		Min: 100,
		Max: 200,
	}

	assert.Equal(t, 100, constraints.Min)
	assert.Equal(t, 200, constraints.Max)
}

func TestBoxConstraints(t *testing.T) {
	constraints := BoxConstraints{
		MinWidth:  100,
		MaxWidth:  200,
		MinHeight: 150,
		MaxHeight: 250,
	}

	assert.Equal(t, 100, constraints.MinWidth)
	assert.Equal(t, 200, constraints.MaxWidth)
	assert.Equal(t, 150, constraints.MinHeight)
	assert.Equal(t, 250, constraints.MaxHeight)
}

func TestBoxConstraints_ToConstraints(t *testing.T) {
	boxConstraints := BoxConstraints{
		MinWidth:  100,
		MaxWidth:  200,
		MinHeight: 150,
		MaxHeight: 250,
	}

	constraints := boxConstraints.ToConstraints()

	// Test MinSize
	assert.Equal(t, geometry.Size{
		Width:  100,
		Height: 150,
	}, constraints.MinSize)

	// Test MaxSize
	assert.Equal(t, geometry.Size{
		Width:  200,
		Height: 250,
	}, constraints.MaxSize)
}

func TestLooseConstraints(t *testing.T) {
	minSize := geometry.Size{
		Width:  100,
		Height: 150,
	}

	constraints := LooseConstraints(minSize)

	// Test MinSize matches input
	assert.Equal(t, minSize, constraints.MinSize)

	// Test MaxSize is maximum possible
	assert.Equal(t, geometry.Size{
		Width:  math.MaxInt32,
		Height: math.MaxInt32,
	}, constraints.MaxSize)
}

func TestTightConstraints(t *testing.T) {
	size := geometry.Size{
		Width:  100,
		Height: 150,
	}

	constraints := TightConstraints(size)

	// Test both MinSize and MaxSize match input
	assert.Equal(t, size, constraints.MinSize)
	assert.Equal(t, size, constraints.MaxSize)

	// Verify constraints are actually tight
	assert.True(t, constraints.IsTight())
}

func TestExpandedConstraints(t *testing.T) {
	maxSize := geometry.Size{
		Width:  200,
		Height: 300,
	}

	constraints := ExpandedConstraints(maxSize)

	// Test both MinSize and MaxSize match input
	assert.Equal(t, maxSize, constraints.MinSize)
	assert.Equal(t, maxSize, constraints.MaxSize)

	// Verify constraints are tight (expanded constraints are always tight)
	assert.True(t, constraints.IsTight())

	// Verify they're different from regular tight constraints
	tightConstraints := TightConstraints(maxSize)
	assert.Equal(t, constraints, tightConstraints,
		"ExpandedConstraints should be equivalent to TightConstraints for the same size")
}
