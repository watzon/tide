// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package core

type Color struct {
	R, G, B uint8
	A       uint8
}

func (c Color) RGBA() (uint32, uint32, uint32, uint32) {
	r := uint32(c.R)
	g := uint32(c.G)
	b := uint32(c.B)
	a := uint32(c.A)

	if a == 0xff {
		return r << 8, g << 8, b << 8, a << 8
	}

	r = (r * a) / 0xff
	g = (g * a) / 0xff
	b = (b * a) / 0xff

	return r << 8, g << 8, b << 8, a << 8
}
