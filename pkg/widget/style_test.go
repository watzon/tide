// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/watzon/tide/pkg/core/capabilities"
	"github.com/watzon/tide/pkg/core/color"
)

func TestNewWidgetStyle(t *testing.T) {
	s := NewWidgetStyle()

	// Test default values
	assert.Equal(t, color.White, s.ForegroundColor)
	assert.Equal(t, color.Transparent, s.BackgroundColor)
	assert.Equal(t, EdgeInsets{}, s.Padding)
	assert.Equal(t, EdgeInsets{}, s.Margin)
	assert.Equal(t, BorderNone, s.BorderStyle)
	assert.Equal(t, color.Transparent, s.BorderColor)
	assert.Equal(t, EdgeInsets{}, s.BorderWidth)
}

func TestWidgetStyle_WithMethods(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(WidgetStyle) WidgetStyle
		verify   func(*testing.T, WidgetStyle)
	}{
		{
			name: "WithForeground",
			modifier: func(s WidgetStyle) WidgetStyle {
				return s.WithForeground(color.Red)
			},
			verify: func(t *testing.T, s WidgetStyle) {
				assert.Equal(t, color.Red, s.ForegroundColor)
			},
		},
		{
			name: "WithBackground",
			modifier: func(s WidgetStyle) WidgetStyle {
				return s.WithBackground(color.Blue)
			},
			verify: func(t *testing.T, s WidgetStyle) {
				assert.Equal(t, color.Blue, s.BackgroundColor)
			},
		},
		{
			name: "WithBold",
			modifier: func(s WidgetStyle) WidgetStyle {
				return s.WithBold(true)
			},
			verify: func(t *testing.T, s WidgetStyle) {
				assert.True(t, s.Bold)
			},
		},
		{
			name: "WithItalic",
			modifier: func(s WidgetStyle) WidgetStyle {
				return s.WithItalic(true)
			},
			verify: func(t *testing.T, s WidgetStyle) {
				assert.True(t, s.Italic)
			},
		},
		{
			name: "WithUnderline",
			modifier: func(s WidgetStyle) WidgetStyle {
				return s.WithUnderline(true)
			},
			verify: func(t *testing.T, s WidgetStyle) {
				assert.True(t, s.Underline)
			},
		},
		{
			name: "WithStrikeThrough",
			modifier: func(s WidgetStyle) WidgetStyle {
				return s.WithStrikeThrough(true)
			},
			verify: func(t *testing.T, s WidgetStyle) {
				assert.True(t, s.StrikeThrough)
			},
		},
		{
			name: "WithPadding",
			modifier: func(s WidgetStyle) WidgetStyle {
				return s.WithPadding(EdgeInsetsAll(5))
			},
			verify: func(t *testing.T, s WidgetStyle) {
				assert.Equal(t, EdgeInsetsAll(5), s.Padding)
			},
		},
		{
			name: "WithMargin",
			modifier: func(s WidgetStyle) WidgetStyle {
				return s.WithMargin(EdgeInsetsAll(10))
			},
			verify: func(t *testing.T, s WidgetStyle) {
				assert.Equal(t, EdgeInsetsAll(10), s.Margin)
			},
		},
		{
			name: "WithBorder",
			modifier: func(s WidgetStyle) WidgetStyle {
				return s.WithBorder(BorderSingle, color.Red, EdgeInsetsAll(1))
			},
			verify: func(t *testing.T, s WidgetStyle) {
				assert.Equal(t, BorderSingle, s.BorderStyle)
				assert.Equal(t, color.Red, s.BorderColor)
				assert.Equal(t, EdgeInsetsAll(1), s.BorderWidth)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := NewWidgetStyle()
			modified := tt.modifier(base)
			tt.verify(t, modified)
		})
	}
}

func TestWidgetStyle_Merge(t *testing.T) {
	base := NewWidgetStyle().
		WithForeground(color.Red).
		WithBackground(color.Blue).
		WithBold(true).
		WithPadding(EdgeInsetsAll(5)).
		WithBorder(BorderSingle, color.Green, EdgeInsetsAll(1))

	other := NewWidgetStyle().
		WithForeground(color.Green).
		WithBackground(color.Transparent). // Should not override
		WithItalic(true).
		WithPadding(EdgeInsetsAll(10)).
		WithBorder(BorderDouble, color.Blue, EdgeInsetsAll(2))

	merged := base.Merge(other)

	// Test color merging
	assert.Equal(t, color.Green, merged.ForegroundColor)
	assert.Equal(t, color.Blue, merged.BackgroundColor)

	// Test text properties
	assert.True(t, merged.Bold)
	assert.True(t, merged.Italic)

	// Test layout properties
	assert.Equal(t, EdgeInsetsAll(10), merged.Padding)

	// Test border properties
	assert.Equal(t, BorderDouble, merged.BorderStyle)
	assert.Equal(t, color.Blue, merged.BorderColor)
	assert.Equal(t, EdgeInsetsAll(2), merged.BorderWidth)
}

func TestWidgetStyle_AdaptStyle(t *testing.T) {
	style := NewWidgetStyle().
		WithForeground(color.Red).
		WithBackground(color.Blue).
		WithBold(true).
		WithItalic(true).
		WithUnderline(true).
		WithStrikeThrough(true).
		WithBorder(BorderSingle, color.Green, EdgeInsetsAll(1))

	tests := []struct {
		name   string
		caps   capabilities.Capabilities
		verify func(*testing.T, WidgetStyle)
	}{
		{
			name: "Limited colors",
			caps: capabilities.Capabilities{
				ColorMode: capabilities.Color16,
			},
			verify: func(t *testing.T, s WidgetStyle) {
				// For 16-color mode, colors should be quantized to the basic palette
				quantizedRed := color.Red.QuantizeTo(color.ColorMode(capabilities.Color16))
				quantizedBlue := color.Blue.QuantizeTo(color.ColorMode(capabilities.Color16))
				quantizedGreen := color.Green.QuantizeTo(color.ColorMode(capabilities.Color16))

				assert.Equal(t, quantizedRed, s.ForegroundColor)
				assert.Equal(t, quantizedBlue, s.BackgroundColor)
				assert.Equal(t, quantizedGreen, s.BorderColor)
			},
		},
		{
			name: "No formatting",
			caps: capabilities.Capabilities{
				ColorMode: capabilities.ColorTrueColor,
			},
			verify: func(t *testing.T, s WidgetStyle) {
				assert.False(t, s.Bold)
				assert.False(t, s.Italic)
				assert.False(t, s.Underline)
				assert.False(t, s.StrikeThrough)
			},
		},
		{
			name: "Full capabilities",
			caps: capabilities.Capabilities{
				ColorMode:             capabilities.ColorTrueColor,
				SupportsBold:          true,
				SupportsItalic:        true,
				SupportsUnderline:     true,
				SupportsStrikethrough: true,
			},
			verify: func(t *testing.T, s WidgetStyle) {
				assert.Equal(t, color.Red, s.ForegroundColor)
				assert.Equal(t, color.Blue, s.BackgroundColor)
				assert.Equal(t, color.Green, s.BorderColor)
				assert.True(t, s.Bold)
				assert.True(t, s.Italic)
				assert.True(t, s.Underline)
				assert.True(t, s.StrikeThrough)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapted := style.AdaptStyle(tt.caps)
			tt.verify(t, adapted)
		})
	}
}

func TestWidgetStyle_CommonStyles(t *testing.T) {
	base := NewWidgetStyle().
		WithForeground(color.Red).
		WithBackground(color.Blue)

	t.Run("Disabled", func(t *testing.T) {
		disabled := base.Disabled()
		assert.Equal(t, uint8(128), disabled.ForegroundColor.A)
	})

	t.Run("Selected", func(t *testing.T) {
		selected := base.Selected()
		assert.Equal(t, color.Color{R: 0, G: 0, B: 128, A: 255}, selected.BackgroundColor)
	})

	t.Run("Focused", func(t *testing.T) {
		focused := base.Focused()
		assert.Equal(t, BorderSingle, focused.BorderStyle)
		assert.Equal(t, color.Color{R: 0, G: 128, B: 255, A: 255}, focused.BorderColor)
		assert.Equal(t, EdgeInsets{Top: 1, Right: 1, Bottom: 1, Left: 1}, focused.BorderWidth)
	})
}

func TestWidgetStyle_Presets(t *testing.T) {
	assert.NotNil(t, DefaultStyle)
	assert.NotNil(t, PrimaryStyle)
	assert.NotNil(t, SuccessStyle)
	assert.NotNil(t, WarningStyle)
	assert.NotNil(t, ErrorStyle)

	// Test specific preset colors
	assert.Equal(t, color.Color{R: 0, G: 122, B: 255, A: 255}, PrimaryStyle.ForegroundColor)
	assert.Equal(t, color.Color{R: 40, G: 167, B: 69, A: 255}, SuccessStyle.ForegroundColor)
	assert.Equal(t, color.Color{R: 255, G: 193, B: 7, A: 255}, WarningStyle.ForegroundColor)
	assert.Equal(t, color.Color{R: 220, G: 53, B: 69, A: 255}, ErrorStyle.ForegroundColor)
}
