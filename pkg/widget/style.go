package widget

import (
	"github.com/watzon/tide/pkg/core/capabilities"
	"github.com/watzon/tide/pkg/core/color"
	"github.com/watzon/tide/pkg/core/geometry"
	"github.com/watzon/tide/pkg/core/style"
)

// WidgetStyle extends core.style.Style with widget-specific styling
type WidgetStyle struct {
	style.Style // Embed core style

	// Layout properties
	Padding EdgeInsets
	Margin  EdgeInsets
	MinSize geometry.Size
	MaxSize geometry.Size

	// Border properties
	BorderStyle BorderStyle
	BorderColor color.Color
	BorderWidth EdgeInsets
}

// BorderStyle represents different border types
type BorderStyle int

const (
	BorderNone BorderStyle = iota
	BorderSingle
	BorderDouble
	BorderRounded
	BorderHeavy
	BorderDashed
	BorderDotted
)

// NewWidgetStyle creates a new style with default values
func NewWidgetStyle() WidgetStyle {
	return WidgetStyle{
		Style: style.Style{
			ForegroundColor: color.White,
			BackgroundColor: color.Transparent,
		},
		Padding:     EdgeInsets{},
		Margin:      EdgeInsets{},
		BorderStyle: BorderNone,
		BorderColor: color.Transparent,
		BorderWidth: EdgeInsets{},
	}
}

// Style modification methods (fluent interface)
func (s WidgetStyle) WithForeground(c color.Color) WidgetStyle {
	s.ForegroundColor = c
	return s
}

func (s WidgetStyle) WithBackground(c color.Color) WidgetStyle {
	s.BackgroundColor = c
	return s
}

func (s WidgetStyle) WithBold(bold bool) WidgetStyle {
	s.Bold = bold
	return s
}

func (s WidgetStyle) WithItalic(italic bool) WidgetStyle {
	s.Italic = italic
	return s
}

func (s WidgetStyle) WithUnderline(underline bool) WidgetStyle {
	s.Underline = underline
	return s
}

func (s WidgetStyle) WithStrikeThrough(strikeThrough bool) WidgetStyle {
	s.StrikeThrough = strikeThrough
	return s
}

func (s WidgetStyle) WithPadding(insets EdgeInsets) WidgetStyle {
	s.Padding = insets
	return s
}

func (s WidgetStyle) WithMargin(insets EdgeInsets) WidgetStyle {
	s.Margin = insets
	return s
}

func (s WidgetStyle) WithBorder(style BorderStyle, color color.Color, width EdgeInsets) WidgetStyle {
	s.BorderStyle = style
	s.BorderColor = color
	s.BorderWidth = width
	return s
}

// Merge combines two styles, with the other style taking precedence
func (s WidgetStyle) Merge(other WidgetStyle) WidgetStyle {
	result := s

	// Only override colors if they're not transparent
	if other.ForegroundColor.A > 0 {
		result.ForegroundColor = other.ForegroundColor
	}
	if other.BackgroundColor.A > 0 {
		result.BackgroundColor = other.BackgroundColor
	}

	// Text properties
	result.Bold = result.Bold || other.Bold
	result.Italic = result.Italic || other.Italic
	result.Underline = result.Underline || other.Underline
	result.StrikeThrough = result.StrikeThrough || other.StrikeThrough

	// Layout properties (other takes precedence)
	result.Padding = other.Padding
	result.Margin = other.Margin

	// Border properties
	if other.BorderStyle != BorderNone {
		result.BorderStyle = other.BorderStyle
		result.BorderColor = other.BorderColor
		result.BorderWidth = other.BorderWidth
	}

	return result
}

// AdaptStyle adapts the style for specific backend capabilities
func (s WidgetStyle) AdaptStyle(caps capabilities.Capabilities) WidgetStyle {
	adapted := s

	// Adapt colors based on backend capabilities
	if caps.ColorMode < capabilities.ColorTrueColor {
		adapted.ForegroundColor = adapted.ForegroundColor.QuantizeTo(color.ColorMode(caps.ColorMode))
		adapted.BackgroundColor = adapted.BackgroundColor.QuantizeTo(color.ColorMode(caps.ColorMode))
		if adapted.BorderColor.A > 0 {
			adapted.BorderColor = adapted.BorderColor.QuantizeTo(color.ColorMode(caps.ColorMode))
		}
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

// Helper functions for common style combinations
func (s WidgetStyle) Disabled() WidgetStyle {
	return s.WithForeground(s.ForegroundColor.WithAlpha(128))
}

func (s WidgetStyle) Selected() WidgetStyle {
	return s.WithBackground(color.Color{R: 0, G: 0, B: 128, A: 255})
}

func (s WidgetStyle) Focused() WidgetStyle {
	return s.WithBorder(BorderSingle, color.Color{R: 0, G: 128, B: 255, A: 255}, EdgeInsets{
		Top: 1, Right: 1, Bottom: 1, Left: 1,
	})
}

// Common style presets
var (
	DefaultStyle = NewWidgetStyle()

	PrimaryStyle = NewWidgetStyle().
			WithForeground(color.Color{R: 0, G: 122, B: 255, A: 255})

	SuccessStyle = NewWidgetStyle().
			WithForeground(color.Color{R: 40, G: 167, B: 69, A: 255})

	WarningStyle = NewWidgetStyle().
			WithForeground(color.Color{R: 255, G: 193, B: 7, A: 255})

	ErrorStyle = NewWidgetStyle().
			WithForeground(color.Color{R: 220, G: 53, B: 69, A: 255})
)
