package color

import "math"

// ColorSpace represents different color spaces
type ColorSpace int

const (
	ColorSpaceSRGB ColorSpace = iota
	ColorSpaceLinearRGB
	ColorSpaceDisplayP3
)

// Profile represents a color profile with specific gamut and transfer characteristics
type Profile struct {
	space      ColorSpace
	gamma      float64
	whitePoint [3]float64
}

// Standard profiles
var (
	DefaultProfile = Profile{
		space:      ColorSpaceSRGB,
		gamma:      2.2,
		whitePoint: [3]float64{0.9505, 1.0, 1.0890}, // D65 white point
	}

	LinearProfile = Profile{
		space:      ColorSpaceLinearRGB,
		gamma:      1.0,
		whitePoint: [3]float64{0.9505, 1.0, 1.0890},
	}

	DisplayP3Profile = Profile{
		space:      ColorSpaceDisplayP3,
		gamma:      2.2,
		whitePoint: [3]float64{0.9505, 1.0, 1.0890},
	}
)

// Getter methods for Profile
func (p Profile) Space() ColorSpace {
	return p.space
}

func (p Profile) Gamma() float64 {
	return p.gamma
}

func (p Profile) WhitePoint() [3]float64 {
	return p.whitePoint
}

// ToLinearRGB converts a color to linear RGB space
func (c Color) ToLinearRGB(gamma float64) Color {
	if gamma == 1.0 {
		return c
	}

	toLinear := func(v uint8) uint8 {
		normalized := float64(v) / 255.0
		linear := math.Pow(normalized, gamma)
		return uint8(math.Round(linear * 255.0))
	}

	return Color{
		R: toLinear(c.R),
		G: toLinear(c.G),
		B: toLinear(c.B),
		A: c.A,
	}
}

// FromLinearRGB converts a color from linear RGB space
func (c Color) FromLinearRGB(gamma float64) Color {
	if gamma == 1.0 {
		return c
	}

	fromLinear := func(v uint8) uint8 {
		normalized := float64(v) / 255.0
		nonlinear := math.Pow(normalized, 1.0/gamma)
		return uint8(math.Round(nonlinear * 255.0))
	}

	return Color{
		R: fromLinear(c.R),
		G: fromLinear(c.G),
		B: fromLinear(c.B),
		A: c.A,
	}
}

// ConvertToProfile converts a color from one profile to another
func (c Color) ConvertToProfile(from, to Profile) Color {
	if from == to {
		return c
	}

	// Convert to linear RGB first
	linear := c.ToLinearRGB(from.gamma)

	// Convert to target space
	return linear.FromLinearRGB(to.gamma)
}
