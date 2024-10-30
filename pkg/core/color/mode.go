package color

// ColorMode represents color support levels
type ColorMode int

const (
	ColorNone ColorMode = iota
	Color16
	Color256
	ColorTrueColor
)
