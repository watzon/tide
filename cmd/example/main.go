// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// cmd/example/main.go
// cmd/example/main.go
package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/watzon/tide/pkg/backend/terminal"
	"github.com/watzon/tide/pkg/core"
)

func main() {
	// Initialize terminal
	term, err := terminal.New()
	if err != nil {
		fmt.Printf("Error initializing terminal: %v\n", err)
		os.Exit(1)
	}
	defer term.Shutdown()

	// Enable features
	term.EnableUnicode()
	term.EnableCombiningChars()
	term.EnableMouse()
	term.SetMouseMode(terminal.MouseClick)

	// Set window title
	term.SetTitle("Tide Terminal Demo")

	// Create channels for control
	quit := make(chan struct{})
	done := make(chan struct{})

	// Start event handling in a separate goroutine
	go func() {
		defer close(done)
		term.HandleEvents(func(ev terminal.Event) bool {
			switch ev := ev.(type) {
			case terminal.KeyEvent:
				if ev.Key == tcell.KeyCtrlC || ev.Key == tcell.KeyEsc {
					close(quit)
					return true
				}
				if ev.Key == tcell.KeyCtrlV {
					// Demo clipboard paste
					if content, err := term.GetClipboard(); err == nil {
						term.SetClipboard("Pasted: " + content)
					}
				}
			case terminal.MouseEvent:
				handleMouseClick(term, ev)
			}
			return false
		})
	}()

	// Get terminal capabilities
	caps := term.Capabilities()

	// Get terminal size
	size := term.Size()

	// Calculate center position for our box
	boxWidth := 60
	boxHeight := 20
	startX := (size.Width - boxWidth) / 2
	startY := (size.Height - boxHeight) / 2

	// Colors
	border := core.Color{R: 75, G: 0, B: 130, A: 255}  // Indigo
	title := core.Color{R: 255, G: 215, B: 0, A: 255}  // Gold
	text := core.Color{R: 200, G: 200, B: 200, A: 255} // Light gray
	bg := core.Color{R: 0, G: 0, B: 0, A: 255}         // Pure black
	highlight := core.Color{R: 0, G: 255, B: 127}      // Spring green

	// Animation ticker
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	// Draw loop
drawLoop:
	for {
		select {
		case <-quit:
			break drawLoop
		case <-ticker.C:
			// Clear screen with background color
			term.Clear()

			// Draw box border with double-line characters
			drawBox(term, startX, startY, boxWidth, boxHeight, border, bg, text)

			// Draw title with combining characters
			titleText := "✨ Tide Terminal Demo ♥\u0308" // Heart with diaeresis
			titleX := startX + (boxWidth-term.StringWidth(titleText))/2
			term.DrawText(titleX, startY, titleText, title, bg, terminal.StyleBold)

			// Draw capability information
			type menuItem struct {
				text      string
				style     terminal.StyleMask
				highlight bool
			}

			menuItems := []menuItem{
				{fmt.Sprintf("Color Mode: %v", caps.ColorMode), terminal.StyleBold, false},
				{fmt.Sprintf("Unicode: %v", caps.Unicode), terminal.StyleItalic, false},
				{fmt.Sprintf("Mouse: %v", caps.Mouse), terminal.StyleUnderline, false},
				{"", 0, false},
				{"Interactive Features:", terminal.StyleBold, false},
				{"  • Click anywhere to draw", 0, true},
				{"  • Press Ctrl+V to test clipboard", 0, true},
				{"  • Press Ctrl+C or ESC to exit", 0, true},
				{"", 0, false},
				{"Unicode Examples:", terminal.StyleBold, false},
				{"  • Boxes: ┌─┐│└┘", 0, false},
				{"  • Blocks: █▀▄▌▐", 0, false},
				{"  • Symbols: ★✦✧✪✫✬✭", 0, false},
				{"  • Combining: e\u0301 a\u0308 n\u0303 u\u0308", 0, false},
				{"", 0, false},
				{time.Now().Format("Current Time: 15:04:05"), 0, false},
			}

			pulseValue := float64(time.Now().UnixNano()/int64(time.Millisecond)) * 0.01
			pulseIntensity := (math.Sin(pulseValue) + 1) / 2

			for i, item := range menuItems {
				x := startX + 2
				y := startY + 2 + i

				itemFg := text
				if item.highlight {
					// Create a pulsing highlight effect
					r := uint8(float64(highlight.R) * pulseIntensity)
					g := uint8(float64(highlight.G) * pulseIntensity)
					b := uint8(float64(highlight.B) * pulseIntensity)
					itemFg = core.Color{R: r, G: g, B: b, A: 255}
				}

				drawStyledText(term, x, y, item.text, itemFg, bg, item.style)
			}

			// Draw color spectrum demo
			drawColorSpectrum(term, startX+2, startY+boxHeight-3, boxWidth-4)

			// Present the frame
			term.Present()
		}
	}

	// Wait for the event handler to clean up
	<-done
}

func drawBox(term *terminal.Terminal, x, y, width, height int, borderColor, bgColor, textColor core.Color) {
	// Ensure alpha channels are set
	borderColor.A = 255
	bgColor.A = 255
	textColor.A = 255

	// Box drawing characters
	const (
		topLeft     = '┌'
		topRight    = '┐'
		bottomLeft  = '└'
		bottomRight = '┘'
		horizontal  = '─'
		vertical    = '│'
	)

	// Draw corners with full opacity
	term.DrawCell(x, y, topLeft, borderColor, bgColor)
	term.DrawCell(x+width-1, y, topRight, borderColor, bgColor)
	term.DrawCell(x, y+height-1, bottomLeft, borderColor, bgColor)
	term.DrawCell(x+width-1, y+height-1, bottomRight, borderColor, bgColor)

	// Draw horizontal borders
	for i := 1; i < width-1; i++ {
		term.DrawCell(x+i, y, horizontal, borderColor, bgColor)
		term.DrawCell(x+i, y+height-1, horizontal, borderColor, bgColor)
	}

	// Draw vertical borders
	for i := 1; i < height-1; i++ {
		term.DrawCell(x, y+i, vertical, borderColor, bgColor)
		term.DrawCell(x+width-1, y+i, vertical, borderColor, bgColor)
	}

	// Fill background
	for i := 1; i < width-1; i++ {
		for j := 1; j < height-1; j++ {
			term.DrawCell(x+i, y+j, ' ', textColor, bgColor)
		}
	}
}

func drawStyledText(term *terminal.Terminal, x, y int, text string, fg, bg core.Color, style terminal.StyleMask) {
	// Skip empty strings
	if len(text) == 0 {
		return
	}

	// Ensure alpha channel is set
	fg.A = 255
	bg.A = 255

	for i, ch := range text {
		// Skip zero-width characters
		if ch == 0 || runewidth.RuneWidth(ch) == 0 {
			continue
		}
		term.DrawStyledCell(x+i, y, ch, fg, bg, style)
	}
}

func drawColorSpectrum(term *terminal.Terminal, x, y, width int) {
	for i := 0; i < width; i++ {
		hue := float64(i) / float64(width) * 360.0
		r, g, b := hslToRGB(hue, 1.0, 0.5)
		color := core.Color{R: r, G: g, B: b, A: 255}
		term.DrawCell(x+i, y, '▀', color, color)
	}
}

func handleMouseClick(term *terminal.Terminal, ev terminal.MouseEvent) {
	if ev.Buttons&tcell.ButtonPrimary != 0 {
		x, y := ev.Position.X, ev.Position.Y
		// Draw a small pattern at click location
		pattern := []struct{ dx, dy int }{
			{0, 0}, {1, 0}, {0, 1}, {1, 1}, // 2x2 square
			{-1, 0}, {0, -1}, {1, -1}, {-1, 1}, // surrounding points
		}
		color := core.Color{
			R: uint8(time.Now().UnixNano() % 256),
			G: uint8(time.Now().UnixNano() / 256 % 256),
			B: uint8(time.Now().UnixNano() / 65536 % 256),
			A: 255,
		}
		for _, p := range pattern {
			term.DrawCell(x+p.dx, y+p.dy, '•', color, core.Color{})
		}
	}
}

// HSL to RGB conversion helper
func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	var r, g, b float64

	if s == 0 {
		r, g, b = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q

		r = hueToRGB(p, q, h/360+1.0/3.0)
		g = hueToRGB(p, q, h/360)
		b = hueToRGB(p, q, h/360-1.0/3.0)
	}

	return uint8(r * 255), uint8(g * 255), uint8(b * 255)
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}
