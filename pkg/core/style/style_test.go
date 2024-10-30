// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package style

import (
	"testing"

	"github.com/watzon/tide/pkg/core/capabilities"
	"github.com/watzon/tide/pkg/core/color"
)

func TestAdaptStyle(t *testing.T) {
	tests := []struct {
		name string
		s    Style
		caps capabilities.Capabilities
		want Style
	}{
		{
			name: "No capability restrictions",
			s: Style{
				ForegroundColor: color.Color{R: 255, G: 128, B: 64},
				BackgroundColor: color.Color{R: 64, G: 128, B: 255},
				Bold:            true,
				Italic:          true,
				Underline:       true,
				StrikeThrough:   true,
			},
			caps: capabilities.Capabilities{
				ColorMode:             capabilities.ColorTrueColor,
				SupportsItalic:        true,
				SupportsBold:          true,
				SupportsUnderline:     true,
				SupportsStrikethrough: true,
			},
			want: Style{
				ForegroundColor: color.Color{R: 255, G: 128, B: 64},
				BackgroundColor: color.Color{R: 64, G: 128, B: 255},
				Bold:            true,
				Italic:          true,
				Underline:       true,
				StrikeThrough:   true,
			},
		},
		{
			name: "Limited color mode",
			s: Style{
				ForegroundColor: color.Color{R: 255, G: 128, B: 64},
				BackgroundColor: color.Color{R: 64, G: 128, B: 255},
				Bold:            true,
				Italic:          true,
				Underline:       true,
				StrikeThrough:   true,
			},
			caps: capabilities.Capabilities{
				ColorMode:             capabilities.Color256,
				SupportsItalic:        true,
				SupportsBold:          true,
				SupportsUnderline:     true,
				SupportsStrikethrough: true,
			},
			want: Style{
				ForegroundColor: color.Color{R: 255, G: 128, B: 64}.QuantizeTo(color.ColorMode(capabilities.Color256)),
				BackgroundColor: color.Color{R: 64, G: 128, B: 255}.QuantizeTo(color.ColorMode(capabilities.Color256)),
				Bold:            true,
				Italic:          true,
				Underline:       true,
				StrikeThrough:   true,
			},
		},
		{
			name: "No style support",
			s: Style{
				ForegroundColor: color.Color{R: 255, G: 128, B: 64},
				BackgroundColor: color.Color{R: 64, G: 128, B: 255},
				Bold:            true,
				Italic:          true,
				Underline:       true,
				StrikeThrough:   true,
			},
			caps: capabilities.Capabilities{
				ColorMode:             capabilities.ColorTrueColor,
				SupportsItalic:        false,
				SupportsBold:          false,
				SupportsUnderline:     false,
				SupportsStrikethrough: false,
			},
			want: Style{
				ForegroundColor: color.Color{R: 255, G: 128, B: 64},
				BackgroundColor: color.Color{R: 64, G: 128, B: 255},
				Bold:            false,
				Italic:          false,
				Underline:       false,
				StrikeThrough:   false,
			},
		},
		{
			name: "Mixed support",
			s: Style{
				ForegroundColor: color.Color{R: 255, G: 128, B: 64},
				BackgroundColor: color.Color{R: 64, G: 128, B: 255},
				Bold:            true,
				Italic:          true,
				Underline:       true,
				StrikeThrough:   true,
			},
			caps: capabilities.Capabilities{
				ColorMode:             capabilities.Color256,
				SupportsItalic:        true,
				SupportsBold:          false,
				SupportsUnderline:     true,
				SupportsStrikethrough: false,
			},
			want: Style{
				ForegroundColor: color.Color{R: 255, G: 128, B: 64}.QuantizeTo(color.ColorMode(capabilities.Color256)),
				BackgroundColor: color.Color{R: 64, G: 128, B: 255}.QuantizeTo(color.ColorMode(capabilities.Color256)),
				Bold:            false,
				Italic:          true,
				Underline:       true,
				StrikeThrough:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.AdaptStyle(tt.caps)
			if got != tt.want {
				t.Errorf("AdaptStyle() = %v, want %v", got, tt.want)
			}
		})
	}
}
