// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// cmd/example/main.go
// cmd/example/main.go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
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

	// Enable Unicode and combining characters for better rendering
	term.EnableUnicode()
	term.EnableCombiningChars()

	// Create channels for control
	quit := make(chan struct{})

	// Start event handling in a separate goroutine
	go func() {
		term.HandleEvents(func(ev terminal.Event) bool {
			switch ev := ev.(type) {
			case terminal.KeyEvent:
				if ev.Key == tcell.KeyCtrlC || ev.Key == tcell.KeyEsc {
					close(quit)
					return true
				}
			case terminal.MouseEvent:
				// Optional: Handle mouse events here
			}
			return false
		})
	}()

	// Get terminal size
	size := term.Size()

	// Calculate center position for our box
	boxWidth := 40
	boxHeight := 10
	startX := (size.Width - boxWidth) / 2
	startY := (size.Height - boxHeight) / 2

	// Colors
	border := core.Color{R: 75, G: 0, B: 130, A: 255}  // Indigo
	title := core.Color{R: 255, G: 215, B: 0, A: 255}  // Gold
	text := core.Color{R: 255, G: 255, B: 255, A: 255} // White
	bg := core.Color{R: 25, G: 25, B: 25, A: 255}      // Dark gray

	// Animation ticker
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	// Draw loop
	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			// Clear screen with background color
			term.Clear()

			// Draw box border
			for y := startY; y < startY+boxHeight; y++ {
				for x := startX; x < startX+boxWidth; x++ {
					// Corners
					if (x == startX || x == startX+boxWidth-1) &&
						(y == startY || y == startY+boxHeight-1) {
						term.DrawCell(x, y, '+', border, bg)
					} else if y == startY || y == startY+boxHeight-1 {
						// Top and bottom borders
						term.DrawCell(x, y, '-', border, bg)
					} else if x == startX || x == startX+boxWidth-1 {
						// Side borders
						term.DrawCell(x, y, '|', border, bg)
					} else {
						// Box background
						term.DrawCell(x, y, ' ', text, bg)
					}
				}
			}

			// Draw title with some fancy combining characters
			titleText := "Tide Terminal Demo ♥\u0308" // Heart with diaeresis
			titleX := startX + (boxWidth-len(titleText))/2
			for i, ch := range titleText {
				term.DrawStyledCell(titleX+i, startY, ch, title, bg, terminal.StyleBold)
			}

			// Draw some sample text with combining characters
			messages := []string{
				"Welcome to Tide!",
				"A modern TUI framework for Go",
				"",
				"Some Unicode examples:",
				"  • Combining: e\u0301 a\u0308 n\u0303", // é ä ñ
				"  • Symbols: ★ ◆ ○ ◇ ♠ ♣ ♥ ♦",
				"",
				"Press Ctrl+C or ESC to exit",
			}

			for i, msg := range messages {
				msgX := startX + (boxWidth-len(msg))/2
				for j, ch := range msg {
					term.DrawCell(msgX+j, startY+2+i, ch, text, bg)
				}
			}

			// Draw a simple animation
			spinChars := []rune{'◜', '◝', '◞', '◟'}
			spinChar := spinChars[time.Now().UnixNano()/200000000%4]
			term.DrawCell(startX+boxWidth/2, startY+boxHeight-2, spinChar, text, bg)

			// Present the frame
			term.Present()
		}
	}
}
