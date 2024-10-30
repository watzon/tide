package widget

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/watzon/tide/pkg/core/geometry"
)

func TestConstraintsUnbounded(t *testing.T) {
	assert.Equal(t, geometry.Size{Width: 0, Height: 0}, ConstraintsUnbounded.MinSize)
	assert.Equal(t, geometry.Size{Width: math.MaxInt32, Height: math.MaxInt32}, ConstraintsUnbounded.MaxSize)
}

func TestConstraintsTight(t *testing.T) {
	size := geometry.Size{Width: 100, Height: 200}
	constraints := ConstraintsTight(size)

	assert.Equal(t, size, constraints.MinSize)
	assert.Equal(t, size, constraints.MaxSize)
}

func TestNewConstraints(t *testing.T) {
	minSize := geometry.Size{Width: 50, Height: 60}
	maxSize := geometry.Size{Width: 200, Height: 300}

	constraints := NewConstraints(minSize, maxSize)

	assert.Equal(t, minSize, constraints.MinSize)
	assert.Equal(t, maxSize, constraints.MaxSize)
}

func TestConstraints_WithMinSize(t *testing.T) {
	original := NewConstraints(
		geometry.Size{Width: 50, Height: 60},
		geometry.Size{Width: 200, Height: 300},
	)

	// Test normal case
	newMin := geometry.Size{Width: 75, Height: 80}
	modified := original.WithMinSize(newMin)
	assert.Equal(t, newMin, modified.MinSize)

	// Test normalization case - when new min exceeds max, it should keep original max
	newMin = geometry.Size{Width: 250, Height: 350}
	modified = original.WithMinSize(newMin)
	assert.Equal(t, newMin, modified.MinSize)
	assert.Equal(t, newMin, modified.MaxSize)
}

func TestConstraints_WithMaxSize(t *testing.T) {
	original := NewConstraints(
		geometry.Size{Width: 50, Height: 60},
		geometry.Size{Width: 200, Height: 300},
	)

	// Test normal case
	newMax := geometry.Size{Width: 150, Height: 250}
	modified := original.WithMaxSize(newMax)
	assert.Equal(t, newMax, modified.MaxSize)

	// Test normalization case - when new max is below min, min should adjust down
	newMax = geometry.Size{Width: 25, Height: 30}
	modified = original.WithMaxSize(newMax)
	assert.Equal(t, newMax, modified.MinSize)
	assert.Equal(t, newMax, modified.MaxSize)
}

func TestConstraints_Normalize(t *testing.T) {
	// Test case where min > max
	constraints := Constraints{
		MinSize: geometry.Size{Width: 200, Height: 300},
		MaxSize: geometry.Size{Width: 100, Height: 200},
	}

	normalized := constraints.Normalize()
	assert.True(t, normalized.IsNormalized())
	assert.Equal(t, geometry.Size{Width: 100, Height: 200}, normalized.MinSize)
	assert.Equal(t, geometry.Size{Width: 200, Height: 300}, normalized.MaxSize)
}

func TestConstraints_Constrain(t *testing.T) {
	constraints := NewConstraints(
		geometry.Size{Width: 50, Height: 60},
		geometry.Size{Width: 200, Height: 300},
	)

	// Test within bounds
	size := geometry.Size{Width: 100, Height: 150}
	constrained := constraints.Constrain(size)
	assert.Equal(t, size, constrained)

	// Test below minimum
	size = geometry.Size{Width: 25, Height: 30}
	constrained = constraints.Constrain(size)
	assert.Equal(t, constraints.MinSize, constrained)

	// Test above maximum
	size = geometry.Size{Width: 250, Height: 350}
	constrained = constraints.Constrain(size)
	assert.Equal(t, constraints.MaxSize, constrained)
}

func TestConstraints_IsSatisfiedBy(t *testing.T) {
	constraints := NewConstraints(
		geometry.Size{Width: 50, Height: 60},
		geometry.Size{Width: 200, Height: 300},
	)

	// Test satisfied
	assert.True(t, constraints.IsSatisfiedBy(geometry.Size{Width: 100, Height: 150}))

	// Test not satisfied - too small
	assert.False(t, constraints.IsSatisfiedBy(geometry.Size{Width: 25, Height: 30}))

	// Test not satisfied - too large
	assert.False(t, constraints.IsSatisfiedBy(geometry.Size{Width: 250, Height: 350}))
}

func TestConstraints_Tightness(t *testing.T) {
	// Test loose constraints
	loose := NewConstraints(
		geometry.Size{Width: 50, Height: 60},
		geometry.Size{Width: 200, Height: 300},
	)
	assert.False(t, loose.HasTightWidth())
	assert.False(t, loose.HasTightHeight())
	assert.False(t, loose.IsTight())

	// Test tight width only
	tightWidth := NewConstraints(
		geometry.Size{Width: 100, Height: 60},
		geometry.Size{Width: 100, Height: 300},
	)
	assert.True(t, tightWidth.HasTightWidth())
	assert.False(t, tightWidth.HasTightHeight())
	assert.False(t, tightWidth.IsTight())

	// Test tight height only
	tightHeight := NewConstraints(
		geometry.Size{Width: 50, Height: 200},
		geometry.Size{Width: 200, Height: 200},
	)
	assert.False(t, tightHeight.HasTightWidth())
	assert.True(t, tightHeight.HasTightHeight())
	assert.False(t, tightHeight.IsTight())

	// Test fully tight
	tight := ConstraintsTight(geometry.Size{Width: 100, Height: 100})
	assert.True(t, tight.HasTightWidth())
	assert.True(t, tight.HasTightHeight())
	assert.True(t, tight.IsTight())
}

func TestConstraints_IsNormalized(t *testing.T) {
	// Test normalized constraints
	normalized := NewConstraints(
		geometry.Size{Width: 50, Height: 60},
		geometry.Size{Width: 200, Height: 300},
	)
	assert.True(t, normalized.IsNormalized())

	// Test non-normalized constraints
	nonNormalized := Constraints{
		MinSize: geometry.Size{Width: 200, Height: 300},
		MaxSize: geometry.Size{Width: 50, Height: 60},
	}
	assert.False(t, nonNormalized.IsNormalized())
}

func TestConstraints_String(t *testing.T) {
	constraints := NewConstraints(
		geometry.Size{Width: 50, Height: 60},
		geometry.Size{Width: 200, Height: 300},
	)

	expected := "Constraints(min: {50 60}, max: {200 300})"
	assert.Equal(t, expected, constraints.String())
}
