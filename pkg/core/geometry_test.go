// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package core_test

import (
	"testing"

	"github.com/watzon/tide/pkg/core"
)

func TestRect(t *testing.T) {
	t.Run("NewRect", func(t *testing.T) {
		rect := core.NewRect(10, 20, 30, 40)

		if rect.Min.X != 10 || rect.Min.Y != 20 {
			t.Errorf("expected Min point (10,20), got (%d,%d)", rect.Min.X, rect.Min.Y)
		}

		if rect.Max.X != 40 || rect.Max.Y != 60 {
			t.Errorf("expected Max point (40,60), got (%d,%d)", rect.Max.X, rect.Max.Y)
		}
	})

	t.Run("Size", func(t *testing.T) {
		rect := core.NewRect(10, 20, 30, 40)
		size := rect.Size()

		if size.Width != 30 || size.Height != 40 {
			t.Errorf("expected size (30,40), got (%d,%d)", size.Width, size.Height)
		}
	})

	t.Run("Contains", func(t *testing.T) {
		rect := core.NewRect(10, 20, 30, 40)
		tests := []struct {
			point    core.Point
			expected bool
			name     string
		}{
			{core.Point{X: 15, Y: 25}, true, "point inside"},
			{core.Point{X: 5, Y: 25}, false, "point left"},
			{core.Point{X: 45, Y: 25}, false, "point right"},
			{core.Point{X: 15, Y: 15}, false, "point above"},
			{core.Point{X: 15, Y: 65}, false, "point below"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := rect.Contains(tt.point); got != tt.expected {
					t.Errorf("Contains() = %v, want %v", got, tt.expected)
				}
			})
		}
	})
}
