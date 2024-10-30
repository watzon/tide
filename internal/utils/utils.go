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

// Helper functions for finding min/max
func Max(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}

func Min(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
