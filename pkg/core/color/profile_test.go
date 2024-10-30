package color_test

import (
	"math"
	"testing"

	"github.com/watzon/tide/pkg/core/color"
)

func TestStandardProfiles(t *testing.T) {
	tests := []struct {
		name    string
		profile color.Profile
		gamma   float64
	}{
		{"Default (sRGB)", color.DefaultProfile, 2.2},
		{"Linear RGB", color.LinearProfile, 1.0},
		{"Display P3", color.DisplayP3Profile, 2.2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.profile.Gamma() != tt.gamma {
				t.Errorf("Expected gamma %v, got %v", tt.gamma, tt.profile.Gamma())
			}
		})
	}
}

func TestProfileConversion(t *testing.T) {
	tests := []struct {
		name     string
		color    color.Color
		from     color.Profile
		to       color.Profile
		expected color.Color
	}{
		{
			name:     "Identity conversion",
			color:    color.Color{R: 128, G: 128, B: 128, A: 255},
			from:     color.DefaultProfile,
			to:       color.DefaultProfile,
			expected: color.Color{R: 128, G: 128, B: 128, A: 255},
		},
		{
			name:     "sRGB to Linear",
			color:    color.Color{R: 255, G: 255, B: 255, A: 255},
			from:     color.DefaultProfile,
			to:       color.LinearProfile,
			expected: color.Color{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:     "Linear to sRGB",
			color:    color.Color{R: 128, G: 128, B: 128, A: 255},
			from:     color.LinearProfile,
			to:       color.DefaultProfile,
			expected: color.Color{R: 186, G: 186, B: 186, A: 255},
		},
		{
			name:     "Preserve alpha",
			color:    color.Color{R: 128, G: 128, B: 128, A: 128},
			from:     color.DefaultProfile,
			to:       color.LinearProfile,
			expected: color.Color{R: 56, G: 56, B: 56, A: 128},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.color.ConvertToProfile(tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("ConvertToProfile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLinearRGBConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    color.Color
		gamma    float64
		expected color.Color
	}{
		{
			name:     "Gamma 1.0 (identity)",
			input:    color.Color{R: 128, G: 128, B: 128, A: 255},
			gamma:    1.0,
			expected: color.Color{R: 128, G: 128, B: 128, A: 255},
		},
		{
			name:     "Gamma 2.2 (sRGB)",
			input:    color.Color{R: 128, G: 128, B: 128, A: 255},
			gamma:    2.2,
			expected: color.Color{R: 56, G: 56, B: 56, A: 255},
		},
		{
			name:     "Black",
			input:    color.Color{R: 0, G: 0, B: 0, A: 255},
			gamma:    2.2,
			expected: color.Color{R: 0, G: 0, B: 0, A: 255},
		},
		{
			name:     "White",
			input:    color.Color{R: 255, G: 255, B: 255, A: 255},
			gamma:    2.2,
			expected: color.Color{R: 255, G: 255, B: 255, A: 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToLinearRGB(tt.gamma)
			if result != tt.expected {
				t.Errorf("ToLinearRGB(%v) = %v, want %v", tt.gamma, result, tt.expected)
			}

			roundtrip := result.FromLinearRGB(tt.gamma)
			if !colorsNearlyEqual(roundtrip, tt.input) {
				t.Errorf("Roundtrip conversion failed: got %v, want %v (±1)", roundtrip, tt.input)
			}
		})
	}
}

func TestProfileEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		color color.Color
		gamma float64
	}{
		{"Zero gamma", color.Color{R: 128, G: 128, B: 128, A: 255}, 0.0},
		{"Negative gamma", color.Color{R: 128, G: 128, B: 128, A: 255}, -1.0},
		{"Very high gamma", color.Color{R: 128, G: 128, B: 128, A: 255}, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result1 := tt.color.ToLinearRGB(tt.gamma)
			result2 := result1.FromLinearRGB(tt.gamma)

			validateEdgeCaseResults(t, tt.color, result1, result2)
		})
	}
}

func TestProfileGetters(t *testing.T) {
	tests := []struct {
		name        string
		profile     color.Profile
		expectSpace color.ColorSpace
		expectWhite [3]float64
	}{
		{
			name:        "Default profile",
			profile:     color.DefaultProfile,
			expectSpace: color.ColorSpaceSRGB,
			expectWhite: [3]float64{0.9505, 1.0, 1.0890},
		},
		{
			name:        "Linear profile",
			profile:     color.LinearProfile,
			expectSpace: color.ColorSpaceLinearRGB,
			expectWhite: [3]float64{0.9505, 1.0, 1.0890},
		},
		{
			name:        "Display P3 profile",
			profile:     color.DisplayP3Profile,
			expectSpace: color.ColorSpaceDisplayP3,
			expectWhite: [3]float64{0.9505, 1.0, 1.0890},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validateProfileGetters(t, tt.profile, tt.expectSpace, tt.expectWhite)
		})
	}
}

// Helper functions
func colorsNearlyEqual(c1, c2 color.Color) bool {
	return math.Abs(float64(c1.R)-float64(c2.R)) <= 1 &&
		math.Abs(float64(c1.G)-float64(c2.G)) <= 1 &&
		math.Abs(float64(c1.B)-float64(c2.B)) <= 1 &&
		c1.A == c2.A
}

func validateEdgeCaseResults(t *testing.T, original, result1, result2 color.Color) {
	t.Helper()
	if result1 == original {
		t.Error("Expected color to be modified with extreme gamma value")
	}
	if result1.A != original.A || result2.A != original.A {
		t.Error("Alpha channel should be preserved through all conversions")
	}
	if result2 == original {
		t.Error("Expected roundtrip conversion with extreme gamma to modify color")
	}
}

func validateProfileGetters(t *testing.T, profile color.Profile, expectSpace color.ColorSpace, expectWhite [3]float64) {
	t.Helper()
	if space := profile.Space(); space != expectSpace {
		t.Errorf("Space() = %v, want %v", space, expectSpace)
	}

	whitePoint := profile.WhitePoint()
	for i, v := range whitePoint {
		if v != expectWhite[i] {
			t.Errorf("WhitePoint()[%d] = %v, want %v", i, v, expectWhite[i])
		}
	}
}
