package utils

// equalRunes compares two rune slices for equality
func EqualRunes(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Clamp(f, low, high float32) float32 {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}

func ClampInt(i, low, high int) int {
	// If bounds are reversed, swap them
	if low > high {
		low, high = high, low
	}

	if i < low {
		return low
	}
	if i > high {
		return high
	}
	return i
}

// IsColorIntensityHigh returns true if the color component is above mid-range
func IsColorIntensityHigh(component uint8) bool {
	return component > 127
}

// ColorComponentToBasic converts a color component to its basic form (0 or 255)
func ColorComponentToBasic(component uint8) uint8 {
	if IsColorIntensityHigh(component) {
		return 255
	}
	return 0
}
