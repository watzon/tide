// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package terminal_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/watzon/tide/pkg/backend/terminal"
	"github.com/watzon/tide/pkg/core/color"
	"github.com/watzon/tide/pkg/core/geometry"
)

type testContext struct {
	screen tcell.Screen
	term   *terminal.Terminal
}

func setupTest(t *testing.T) *testContext {
	// Create a new simulation screen
	screen := tcell.NewSimulationScreen("")
	if err := screen.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}

	config := terminal.DefaultConfig()
	config.PollInterval = time.Millisecond * 10 // Faster polling for tests
	term, err := terminal.NewWithScreen(screen, config)
	if err != nil {
		t.Fatalf("failed to create terminal: %v", err)
	}

	return &testContext{
		screen: screen,
		term:   term,
	}
}

func TestTerminal(t *testing.T) {
	t.Run("initialization", func(t *testing.T) {
		ctx := setupTest(t)
		defer ctx.term.Shutdown()

		size := ctx.term.Size()
		if size.Width == 0 || size.Height == 0 {
			t.Error("expected non-zero terminal size")
		}
	})

	t.Run("resize handling", func(t *testing.T) {
		ctx := setupTest(t)
		defer ctx.term.Shutdown()

		var wg sync.WaitGroup
		wg.Add(1)

		resizeCalled := false
		newSize := geometry.Size{}

		ctx.term.OnResize(func(size geometry.Size) {
			resizeCalled = true
			newSize = size
			wg.Done()
		})

		// Simulate resize
		ctx.screen.SetSize(100, 50)

		// Use EventResize instead of InjectResize
		ctx.screen.PostEvent(tcell.NewEventResize(100, 50))

		// Wait for the resize event to be processed
		wg.Wait()

		if !resizeCalled {
			t.Error("resize callback was not called")
		}

		if newSize.Width != 100 || newSize.Height != 50 {
			t.Errorf("expected size (100,50), got (%d,%d)", newSize.Width, newSize.Height)
		}
	})

	t.Run("focus handling", func(t *testing.T) {
		ctx := setupTest(t)
		defer ctx.term.Shutdown()

		var wg sync.WaitGroup
		wg.Add(1)

		focusChanged := false
		focusState := false

		ctx.term.OnFocusChange(func(focused bool) {
			focusChanged = true
			focusState = focused
			wg.Done()
		})

		// Post focus event instead of injecting
		ctx.screen.PostEvent(tcell.NewEventFocus(true))

		// Wait for the focus event to be processed
		wg.Wait()

		if !focusChanged {
			t.Error("focus callback was not called")
		}

		if !focusState {
			t.Error("focus state not updated correctly")
		}
	})

	t.Run("basic drawing", func(t *testing.T) {
		ctx := setupTest(t)
		defer ctx.term.Shutdown()

		// Draw a character
		ctx.term.DrawCell(0, 0, 'A',
			color.Color{R: 255, G: 255, B: 255, A: 255},
			color.Color{R: 0, G: 0, B: 0, A: 255},
		)
		ctx.term.Present()

		// Get content using type assertion to access SimulationScreen methods
		simScreen := ctx.screen.(tcell.SimulationScreen)
		mainc, _, _, _ := simScreen.GetContent(0, 0)
		if mainc != 'A' {
			t.Errorf("expected 'A', got %c", mainc)
		}
	})

	t.Run("styled drawing", func(t *testing.T) {
		ctx := setupTest(t)
		defer ctx.term.Shutdown()

		tests := []struct {
			name  string
			style terminal.StyleMask
			check func(tcell.Style) bool
		}{
			{
				name:  "bold",
				style: terminal.StyleBold,
				check: func(s tcell.Style) bool {
					_, _, attrs := s.Decompose()
					return attrs&tcell.AttrBold != 0
				},
			},
			{
				name:  "underline",
				style: terminal.StyleUnderline,
				check: func(s tcell.Style) bool {
					_, _, attrs := s.Decompose()
					return attrs&tcell.AttrUnderline != 0
				},
			},
			{
				name:  "italic",
				style: terminal.StyleItalic,
				check: func(s tcell.Style) bool {
					_, _, attrs := s.Decompose()
					return attrs&tcell.AttrItalic != 0
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctx.term.DrawStyledCell(0, 0, 'X',
					color.Color{R: 255, G: 255, B: 255, A: 255},
					color.Color{R: 0, G: 0, B: 0, A: 255},
					tt.style,
				)
				ctx.term.Present()

				simScreen := ctx.screen.(tcell.SimulationScreen)
				_, _, style, _ := simScreen.GetContent(0, 0)
				if !tt.check(style) {
					t.Errorf("expected style %s to be set", tt.name)
				}
			})
		}
	})

	t.Run("region drawing", func(t *testing.T) {
		ctx := setupTest(t)
		defer ctx.term.Shutdown()

		region := geometry.NewRect(1, 1, 3, 3)
		style := tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorBlue)

		ctx.term.DrawRegion(region, style, '█')
		ctx.term.Present()

		// Check all cells in region
		for y := region.Min.Y; y < region.Max.Y; y++ {
			for x := region.Min.X; x < region.Max.X; x++ {
				ch, _, _, _ := ctx.screen.GetContent(x, y)
				if ch != '█' {
					t.Errorf("expected '█' at (%d,%d), got %c", x, y, ch)
				}
			}
		}

		// Check boundary
		ch, _, _, _ := ctx.screen.GetContent(0, 0)
		if ch == '█' {
			t.Error("character drawn outside region")
		}
	})

	t.Run("cursor management", func(t *testing.T) {
		ctx := setupTest(t)
		defer ctx.term.Shutdown()

		tests := []struct {
			name string
			x, y int
		}{
			{"origin", 0, 0},
			{"middle", 10, 10},
			{"negative coords", -1, -1},
			{"out of bounds", 1000, 1000},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctx.term.SetCursor(tt.x, tt.y)
				ctx.term.Present()
				// Note: Can't actually verify cursor position in simulation
				// Just checking that it doesn't panic
			})
		}

		// Test cursor hide
		ctx.term.HideCursor()
		ctx.term.Present()
	})

	t.Run("resize handling", func(t *testing.T) {
		ctx := setupTest(t)
		defer ctx.term.Shutdown()

		var wg sync.WaitGroup
		wg.Add(1)

		resizeCalled := false
		newSize := geometry.Size{}

		ctx.term.OnResize(func(size geometry.Size) {
			resizeCalled = true
			newSize = size
			wg.Done()
		})

		// Simulate resize using type assertion
		simScreen := ctx.screen.(tcell.SimulationScreen)
		simScreen.SetSize(100, 50)

		// Post resize event
		ev := tcell.NewEventResize(100, 50)
		ctx.screen.PostEvent(ev)

		// Wait for the resize event to be processed
		wg.Wait()

		if !resizeCalled {
			t.Error("resize callback was not called")
		}

		if newSize.Width != 100 || newSize.Height != 50 {
			t.Errorf("expected size (100,50), got (%d,%d)", newSize.Width, newSize.Height)
		}
	})

	t.Run("mouse handling", func(t *testing.T) {
		ctx := setupTest(t)
		defer ctx.term.Shutdown()

		// Test mouse click
		ev := tcell.NewEventMouse(1, 1, tcell.Button1, tcell.ModNone)
		ctx.screen.PostEvent(ev)

		// Test mouse release
		ev = tcell.NewEventMouse(1, 1, tcell.ButtonNone, tcell.ModNone)
		ctx.screen.PostEvent(ev)

		// No assertion here as we're just ensuring it doesn't panic
		// Real mouse handling would be tested through the event system
	})

	t.Run("focus handling", func(t *testing.T) {
		ctx := setupTest(t)
		defer ctx.term.Shutdown()

		var wg sync.WaitGroup
		wg.Add(1)

		focusChanged := false
		focusState := false

		ctx.term.OnFocusChange(func(focused bool) {
			focusChanged = true
			focusState = focused
			wg.Done()
		})

		ev := tcell.NewEventFocus(true)
		ctx.screen.PostEvent(ev)

		wg.Wait()

		if !focusChanged {
			t.Error("focus callback was not called")
		}

		if !focusState {
			t.Error("focus state not updated correctly")
		}
	})

	t.Run("suspend/resume", func(t *testing.T) {
		ctx := setupTest(t)
		defer ctx.term.Shutdown()

		suspendCalled := false
		resumeCalled := false

		ctx.term.OnSuspend(func() {
			suspendCalled = true
		})

		ctx.term.OnResume(func() {
			resumeCalled = true
		})

		if err := ctx.term.Suspend(); err != nil {
			t.Errorf("unexpected error on suspend: %v", err)
		}

		if !suspendCalled {
			t.Error("suspend callback was not called")
		}

		if err := ctx.term.Resume(); err != nil {
			t.Errorf("unexpected error on resume: %v", err)
		}

		if !resumeCalled {
			t.Error("resume callback was not called")
		}
	})
}

func TestUnicodeSupport(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		unicode  bool
	}{
		{"ASCII only", "Hello", 5, false},
		{"CJK chars", "你好", 4, true},
		{"Combining chars", "é", 1, true},
		{"Mixed content", "Hello 世界", 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := setupTest(t)
			defer ctx.term.Shutdown()

			if tt.unicode {
				ctx.term.EnableUnicode()
			} else {
				ctx.term.DisableUnicode()
			}

			width := ctx.term.StringWidth(tt.input)
			if width != tt.expected {
				t.Errorf("got width %d, want %d", width, tt.expected)
			}
		})
	}
}

func TestTerminalConcurrency(t *testing.T) {
	ctx := setupTest(t)
	defer ctx.term.Shutdown()

	const goroutines = 10
	const operations = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	errc := make(chan error, goroutines*operations)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				if err := func() error {
					defer func() {
						if r := recover(); r != nil {
							errc <- fmt.Errorf("panic: %v", r)
						}
					}()

					ctx.term.DrawCell(0, 0, 'X',
						color.Color{R: 255, G: 255, B: 255, A: 255},
						color.Color{R: 0, G: 0, B: 0, A: 255},
					)
					ctx.term.Present()
					_ = ctx.term.Size()
					ctx.term.SetCursor(0, 0)
					ctx.term.HideCursor()
					return nil
				}(); err != nil {
					errc <- err
				}
			}
		}()
	}

	wg.Wait()
	close(errc)

	for err := range errc {
		t.Errorf("concurrent operation error: %v", err)
	}
}

func TestCombiningCharacters(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		combining     bool
		expectedWidth int
		expectedCells int
	}{
		{
			name:          "Combining enabled - acute accent",
			input:         "e\u0301", // é composed of 'e' and combining acute accent
			combining:     true,
			expectedWidth: 1,
			expectedCells: 1,
		},
		{
			name:          "Combining disabled - acute accent",
			input:         "e\u0301",
			combining:     false,
			expectedWidth: 2,
			expectedCells: 2,
		},
		{
			name:          "Combining enabled - multiple marks",
			input:         "a\u0308\u0301", // ä́ with diaeresis and acute
			combining:     true,
			expectedWidth: 1,
			expectedCells: 1,
		},
		{
			name:          "Combining disabled - multiple marks",
			input:         "a\u0308\u0301",
			combining:     false,
			expectedWidth: 3,
			expectedCells: 3,
		},
		{
			name:          "Combining enabled - heart with diaeresis",
			input:         "♥\u0308", // Heart with diaeresis
			combining:     true,
			expectedWidth: 1,
			expectedCells: 1,
		},
		{
			name:          "Combining disabled - heart with diaeresis",
			input:         "♥\u0308",
			combining:     false,
			expectedWidth: 2,
			expectedCells: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := setupTest(t)
			defer ctx.term.Shutdown()

			ctx.term.EnableUnicode()
			if tt.combining {
				ctx.term.EnableCombiningChars()
			} else {
				ctx.term.DisableCombiningChars()
			}

			// Draw each character
			x := 0
			for _, ch := range tt.input {
				ctx.term.DrawStyledCell(x, 0, ch,
					color.Color{R: 255, G: 255, B: 255, A: 255},
					color.Color{R: 0, G: 0, B: 0, A: 255},
					0,
				)
				if !tt.combining || !unicode.IsMark(ch) {
					x++
				}
			}

			ctx.term.Present()

			// Verify the number of cells used
			usedCells := 0
			simScreen := ctx.screen.(tcell.SimulationScreen)
			for i := 0; i < tt.expectedWidth+1; i++ {
				mainc, combc, _, _ := simScreen.GetContent(i, 0)
				if mainc != ' ' || len(combc) > 0 {
					usedCells++
				}
			}

			if usedCells != tt.expectedCells {
				t.Errorf("got %d cells used, want %d", usedCells, tt.expectedCells)
			}
		})
	}
}
