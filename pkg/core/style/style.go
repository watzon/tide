// Copyright (c) 2024 Chris Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package style

import (
	"github.com/watzon/tide/pkg/core/capabilities"
	"github.com/watzon/tide/pkg/core/color"
)

// Style represents visual properties in a backend-agnostic way
type Style struct {
	// Colors
	ForegroundColor color.Color
	BackgroundColor color.Color

	// Text properties
	Bold          bool
	Italic        bool
	Underline     bool
	StrikeThrough bool
}

// AdaptStyle adapts the style for specific backend capabilities
func (s Style) AdaptStyle(caps capabilities.Capabilities) Style {
	adapted := s

	// Adapt colors based on backend capabilities
	if caps.ColorMode < capabilities.ColorTrueColor {
		adapted.ForegroundColor = adapted.ForegroundColor.QuantizeTo(color.ColorMode(caps.ColorMode))
		adapted.BackgroundColor = adapted.BackgroundColor.QuantizeTo(color.ColorMode(caps.ColorMode))
	}

	// Remove unsupported text styles
	if !caps.SupportsItalic {
		adapted.Italic = false
	}
	if !caps.SupportsBold {
		adapted.Bold = false
	}
	if !caps.SupportsUnderline {
		adapted.Underline = false
	}
	if !caps.SupportsStrikethrough {
		adapted.StrikeThrough = false
	}

	return adapted
}
