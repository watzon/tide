// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package core_test

import (
	"testing"

	"github.com/watzon/tide/pkg/core"
)

func TestColor(t *testing.T) {
	t.Run("RGBA fully opaque", func(t *testing.T) {
		c := core.Color{R: 255, G: 128, B: 64, A: 255}
		r, g, b, a := c.RGBA()

		if r>>8 != 255 || g>>8 != 128 || b>>8 != 64 || a>>8 != 255 {
			t.Errorf("expected RGBA (255,128,64,255), got (%d,%d,%d,%d)", r>>8, g>>8, b>>8, a>>8)
		}
	})

	t.Run("RGBA with alpha", func(t *testing.T) {
		c := core.Color{R: 255, G: 128, B: 64, A: 128}
		r, g, b, a := c.RGBA()

		// Check that alpha is properly applied
		if r>>8 != 128 || g>>8 != 64 || b>>8 != 32 || a>>8 != 128 {
			t.Errorf("expected RGBA (128,64,32,128), got (%d,%d,%d,%d)", r>>8, g>>8, b>>8, a>>8)
		}
	})
}
