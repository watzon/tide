package main

import (
	"log"

	"github.com/watzon/tide/pkg/backend/terminal"
	"github.com/watzon/tide/pkg/core/color"
	"github.com/watzon/tide/pkg/core/geometry"
	"github.com/watzon/tide/pkg/engine"
	"github.com/watzon/tide/pkg/widget"
)

func main() {
	// Initialize terminal
	term, err := terminal.New()
	if err != nil {
		log.Fatalf("failed to initialize terminal: %v", err)
	}
	defer term.Shutdown()

	term.HideCursor()
	term.Clear()

	// Create render context
	ctx := engine.NewTerminalContext(term)

	// Create a simple text widget
	text := widget.NewText("Hello, World!")
	text.WithStyle(widget.NewWidgetStyle().
		WithForeground(color.White).
		WithBackground(color.Blue))

	// Create and mount element
	element := widget.NewElement(text)
	element.Mount(nil)

	// Set constraints based on terminal size
	termSize := term.Size()
	constraints := widget.NewConstraints(
		geometry.Size{Width: 0, Height: 0},
		geometry.Size{Width: termSize.Width, Height: termSize.Height},
	)

	// Apply constraints to text widget
	if textWidget, ok := element.Widget().(*widget.Text); ok {
		textWidget.WithConstraints(constraints)
	}

	// Layout and paint
	element.LayoutPhase()
	if renderObj := element.RenderObject(); renderObj != nil {
		// Debug print
		log.Printf("Render object size: %v", renderObj.Size())

		// Paint
		renderObj.Paint(ctx)
	}

	// Present to screen
	if err := ctx.Present(); err != nil {
		log.Fatalf("failed to present: %v", err)
	}

	// Wait for any event before closing
	term.HandleEvents(func(event terminal.Event) bool {
		return true // Return true to stop handling events
	})
}
