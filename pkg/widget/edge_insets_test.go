// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEdgeInsetsZero(t *testing.T) {
	assert.Equal(t, 0, EdgeInsetsZero.Top)
	assert.Equal(t, 0, EdgeInsetsZero.Right)
	assert.Equal(t, 0, EdgeInsetsZero.Bottom)
	assert.Equal(t, 0, EdgeInsetsZero.Left)
}

func TestEdgeInsetsAll(t *testing.T) {
	insets := EdgeInsetsAll(10)
	assert.Equal(t, 10, insets.Top)
	assert.Equal(t, 10, insets.Right)
	assert.Equal(t, 10, insets.Bottom)
	assert.Equal(t, 10, insets.Left)
}

func TestEdgeInsetsSymmetric(t *testing.T) {
	insets := EdgeInsetsSymmetric(5, 10)
	assert.Equal(t, 5, insets.Top)
	assert.Equal(t, 10, insets.Right)
	assert.Equal(t, 5, insets.Bottom)
	assert.Equal(t, 10, insets.Left)
}

func TestNewEdgeInsets(t *testing.T) {
	insets := NewEdgeInsets(1, 2, 3, 4)
	assert.Equal(t, 1, insets.Top)
	assert.Equal(t, 2, insets.Right)
	assert.Equal(t, 3, insets.Bottom)
	assert.Equal(t, 4, insets.Left)
}

func TestEdgeInsets_Horizontal(t *testing.T) {
	insets := NewEdgeInsets(1, 2, 3, 4)
	assert.Equal(t, 6, insets.Horizontal()) // 2 + 4
}

func TestEdgeInsets_Vertical(t *testing.T) {
	insets := NewEdgeInsets(1, 2, 3, 4)
	assert.Equal(t, 4, insets.Vertical()) // 1 + 3
}

func TestEdgeInsets_Add(t *testing.T) {
	insets1 := NewEdgeInsets(1, 2, 3, 4)
	insets2 := NewEdgeInsets(2, 3, 4, 5)

	result := insets1.Add(insets2)
	assert.Equal(t, 3, result.Top)    // 1 + 2
	assert.Equal(t, 5, result.Right)  // 2 + 3
	assert.Equal(t, 7, result.Bottom) // 3 + 4
	assert.Equal(t, 9, result.Left)   // 4 + 5
}

func TestEdgeInsets_Scale(t *testing.T) {
	insets := NewEdgeInsets(1, 2, 3, 4)

	result := insets.Scale(2)
	assert.Equal(t, 2, result.Top)
	assert.Equal(t, 4, result.Right)
	assert.Equal(t, 6, result.Bottom)
	assert.Equal(t, 8, result.Left)

	// Test scale by 0
	result = insets.Scale(0)
	assert.True(t, result.IsZero())
}

func TestEdgeInsets_IsZero(t *testing.T) {
	// Test zero insets
	assert.True(t, EdgeInsetsZero.IsZero())

	// Test non-zero insets
	insets := NewEdgeInsets(1, 2, 3, 4)
	assert.False(t, insets.IsZero())

	// Test partially zero insets
	insets = NewEdgeInsets(0, 1, 0, 0)
	assert.False(t, insets.IsZero())
}

func TestEdgeInsets_Max(t *testing.T) {
	insets1 := NewEdgeInsets(1, 4, 3, 2)
	insets2 := NewEdgeInsets(2, 3, 4, 1)

	result := insets1.Max(insets2)
	assert.Equal(t, 2, result.Top)    // max(1, 2)
	assert.Equal(t, 4, result.Right)  // max(4, 3)
	assert.Equal(t, 4, result.Bottom) // max(3, 4)
	assert.Equal(t, 2, result.Left)   // max(2, 1)
}

func TestEdgeInsets_Min(t *testing.T) {
	insets1 := NewEdgeInsets(1, 4, 3, 2)
	insets2 := NewEdgeInsets(2, 3, 4, 1)

	result := insets1.Min(insets2)
	assert.Equal(t, 1, result.Top)    // min(1, 2)
	assert.Equal(t, 3, result.Right)  // min(4, 3)
	assert.Equal(t, 3, result.Bottom) // min(3, 4)
	assert.Equal(t, 1, result.Left)   // min(2, 1)
}

func TestEdgeInsets_Clamp(t *testing.T) {
	insets := NewEdgeInsets(1, 5, 3, 2)

	// Test normal case where min < max
	result := insets.Clamp(2, 4)
	assert.Equal(t, 2, result.Top)    // clamp(1, 2, 4)
	assert.Equal(t, 4, result.Right)  // clamp(5, 2, 4)
	assert.Equal(t, 3, result.Bottom) // clamp(3, 2, 4)
	assert.Equal(t, 2, result.Left)   // clamp(2, 2, 4)

	// Test case where min > max (4, 2)
	// ClampInt will swap the bounds to (2, 4)
	// So this behaves the same as the normal case
	result = insets.Clamp(4, 2)
	assert.Equal(t, 2, result.Top)    // clamp(1, 2, 4)
	assert.Equal(t, 4, result.Right)  // clamp(5, 2, 4)
	assert.Equal(t, 3, result.Bottom) // clamp(3, 2, 4)
	assert.Equal(t, 2, result.Left)   // clamp(2, 2, 4)
}
