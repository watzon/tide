// Copyright (c) 2024 Chris Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package engine

import (
	"github.com/watzon/tide/pkg/core/color"
	"github.com/watzon/tide/pkg/core/geometry"
)

// FillRect fills a rectangle with the given background color
func FillRect(ctx RenderContext, rect geometry.Rect, fg, bg color.Color) {
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			ctx.DrawCell(x, y, ' ', fg, bg)
		}
	}
}
