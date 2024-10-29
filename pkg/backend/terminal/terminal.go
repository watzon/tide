// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package terminal

import (
	"fmt"
	"sync"
	"time"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/watzon/tide/pkg/core"
)

// StyleMask represents different text style attributes
type StyleMask uint16

const (
	StyleBold StyleMask = 1 << iota
	StyleBlink
	StyleReverse
	StyleUnderline
	StyleDim
	StyleItalic
	StyleStrikethrough
)

// MouseMode represents different mouse handling modes
type MouseMode int

const (
	MouseDisabled MouseMode = iota
	MouseClick
	MouseDrag
	MouseMotion
)

// Event represents a terminal event
type Event interface {
	When() time.Time
}

// Terminal represents a terminal backend
type Terminal struct {
	screen            tcell.Screen
	style             tcell.Style
	colorOptimizer    *ColorOptimizer
	clipboardProvider ClipboardProvider

	// State
	size      core.Size
	mouseMode MouseMode
	focused   bool
	suspended bool
	lock      sync.RWMutex
	eventChan chan Event
	stopChan  chan struct{}

	// Callbacks
	onResize      func(core.Size)
	onFocusChange func(bool)
	onSuspend     func()
	onResume      func()

	// Unicode
	unicodeMode    bool
	combiningChars bool
	title          string // Track the current window title

	// Buffer management
	frontBuffer *Buffer
	backBuffer  *Buffer
}

// Config holds terminal configuration
type Config struct {
	EnableMouse   bool
	MouseMode     MouseMode
	ColorMode     tcell.Color
	PollInterval  time.Duration
	HandleSuspend bool
	HandleResize  bool
	CaptureEvents bool
}

// DefaultConfig returns the default terminal configuration
func DefaultConfig() *Config {
	return &Config{
		EnableMouse:   true,
		MouseMode:     MouseClick,
		ColorMode:     tcell.ColorDefault,
		PollInterval:  time.Millisecond * 50,
		HandleSuspend: true,
		HandleResize:  true,
		CaptureEvents: true,
	}
}

// New creates a new terminal with default configuration
func New() (*Terminal, error) {
	return NewWithConfig(DefaultConfig())
}

// NewWithConfig creates a new terminal with the provided configuration
func NewWithConfig(config *Config) (*Terminal, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("failed to create screen: %w", err)
	}

	return NewWithScreen(screen, config)
}

// NewWithScreen creates a new terminal with a provided screen
func NewWithScreen(screen tcell.Screen, config *Config) (*Terminal, error) {
	if err := screen.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize screen: %w", err)
	}

	width, height := screen.Size()
	size := core.Size{Width: width, Height: height}

	t := &Terminal{
		screen:         screen,
		style:          tcell.StyleDefault,
		size:           size,
		mouseMode:      config.MouseMode,
		eventChan:      make(chan Event, 100),
		stopChan:       make(chan struct{}),
		combiningChars: true,
		frontBuffer:    NewBuffer(size),
		backBuffer:     NewBuffer(size),
	}

	if config.EnableMouse {
		t.EnableMouse()
	}

	if config.CaptureEvents {
		go t.eventLoop(config.PollInterval)
	}

	if t.SupportsUnicode() {
		t.EnableUnicode()
	}

	return t, nil
}

// Screen management

func (t *Terminal) Init() error {
	return nil
}

func (t *Terminal) Shutdown() error {
	close(t.stopChan)
	t.screen.Fini()
	return nil
}

func (t *Terminal) Suspend() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.suspended = true
	t.screen.Fini()
	if t.onSuspend != nil {
		t.onSuspend()
	}
	return nil
}

func (t *Terminal) Resume() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if err := t.screen.Init(); err != nil {
		return err
	}
	t.suspended = false
	if t.onResume != nil {
		t.onResume()
	}
	return nil
}

func (t *Terminal) Sync() {
	t.screen.Sync()
}

// Drawing operations

func (t *Terminal) Clear() {
	t.screen.Clear()
}

func (t *Terminal) DrawCell(x, y int, ch rune, fg, bg core.Color) {
	t.DrawStyledCell(x, y, ch, fg, bg, 0)
}

func (t *Terminal) DrawStyledCell(x, y int, ch rune, fg, bg core.Color, style StyleMask) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	// Start with optimized colors
	tcellStyle := tcell.StyleDefault.
		Foreground(t.optimizeColor(fg)).
		Background(t.optimizeColor(bg))

	// Apply style attributes
	if style&StyleBold != 0 {
		tcellStyle = tcellStyle.Bold(true)
	}
	if style&StyleBlink != 0 {
		tcellStyle = tcellStyle.Blink(true)
	}
	if style&StyleReverse != 0 {
		tcellStyle = tcellStyle.Reverse(true)
	}
	if style&StyleUnderline != 0 {
		tcellStyle = tcellStyle.Underline(true)
	}
	if style&StyleDim != 0 {
		tcellStyle = tcellStyle.Dim(true)
	}
	if style&StyleItalic != 0 {
		tcellStyle = tcellStyle.Italic(true)
	}
	if style&StyleStrikethrough != 0 {
		tcellStyle = tcellStyle.StrikeThrough(true)
	}

	// When combining is disabled or it's a regular character
	if !t.combiningChars {
		// Draw combining marks as standalone characters when combining is disabled
		if unicode.IsMark(ch) {
			// Use a visible character for combining marks when they're standalone
			t.backBuffer.SetCell(x, y, '\u25CC', []rune{ch}, tcellStyle)
			t.screen.SetContent(x, y, '\u25CC', []rune{ch}, tcellStyle)
			return
		}
	}

	// Handle combining characters when enabled
	if t.unicodeMode && t.combiningChars && unicode.IsMark(ch) {
		prevCell, exists := t.backBuffer.GetCell(x-1, y)
		if exists && prevCell.Rune != ' ' {
			combining := append(prevCell.Combining, ch)
			t.backBuffer.SetCell(x-1, y, prevCell.Rune, combining, tcellStyle)
			t.screen.SetContent(x-1, y, prevCell.Rune, combining, tcellStyle)
			return
		}
	}

	// Normal character handling
	t.backBuffer.SetCell(x, y, ch, nil, tcellStyle)
	t.screen.SetContent(x, y, ch, nil, tcellStyle)
}

func (t *Terminal) DrawRegion(region core.Rect, style tcell.Style, ch rune) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	for y := region.Min.Y; y < region.Max.Y; y++ {
		for x := region.Min.X; x < region.Max.X; x++ {
			t.screen.SetContent(x, y, ch, nil, style)
		}
	}
}

func (t *Terminal) StringWidth(s string) int {
	if !t.unicodeMode {
		return len(s)
	}

	if !t.combiningChars {
		// When combining chars are disabled, count each rune separately
		return len([]rune(s))
	}

	// Use runewidth for normal Unicode width calculation
	return runewidth.StringWidth(s)
}

func (t *Terminal) Present() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.backBuffer.dirty {
		return nil
	}

	t.backBuffer.lock.RLock()
	t.frontBuffer.lock.RLock()
	defer t.backBuffer.lock.RUnlock()
	defer t.frontBuffer.lock.RUnlock()

	for y := 0; y < t.size.Height; y++ {
		for x := 0; x < t.size.Width; x++ {
			pos := core.Point{X: x, Y: y}

			backCell, backExists := t.backBuffer.cells[pos]
			frontCell, frontExists := t.frontBuffer.cells[pos]

			if backExists && frontExists &&
				backCell.Rune == frontCell.Rune &&
				backCell.Style == frontCell.Style &&
				equalRunes(backCell.Combining, frontCell.Combining) {
				continue
			}

			if backExists {
				if !t.combiningChars && unicode.IsMark(backCell.Rune) {
					// For disabled combining chars, show with dotted circle
					t.screen.SetContent(x, y, '\u25CC', []rune{backCell.Rune}, backCell.Style)
				} else {
					t.screen.SetContent(x, y, backCell.Rune, backCell.Combining, backCell.Style)
				}
			} else {
				t.screen.SetContent(x, y, ' ', nil, tcell.StyleDefault)
			}
		}
	}

	cursor := t.backBuffer.GetCursor()
	t.screen.ShowCursor(cursor.X, cursor.Y)

	t.screen.Show()
	t.backBuffer.dirty = false
	return nil
}

// Size and cursor management

func (t *Terminal) Size() core.Size {
	t.lock.RLock()
	defer t.lock.RUnlock()

	width, height := t.screen.Size()
	return core.Size{Width: width, Height: height}
}

func (t *Terminal) SetCursor(x, y int) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.backBuffer.SetCursor(x, y)
}

func (t *Terminal) GetCursor() core.Point {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.backBuffer.GetCursor()
}

func (t *Terminal) HideCursor() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.screen.HideCursor()
}

// Mouse handling

func (t *Terminal) EnableMouse() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.screen.EnableMouse()
}

func (t *Terminal) DisableMouse() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.screen.DisableMouse()
}

func (t *Terminal) SetMouseMode(mode MouseMode) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.mouseMode = mode
	switch mode {
	case MouseDisabled:
		t.screen.DisableMouse()
	case MouseClick:
		t.screen.EnableMouse(tcell.MouseButtonEvents)
	case MouseDrag:
		t.screen.EnableMouse(tcell.MouseButtonEvents, tcell.MouseDragEvents)
	case MouseMotion:
		t.screen.EnableMouse(tcell.MouseButtonEvents, tcell.MouseDragEvents, tcell.MouseMotionEvents)
	}
}

// Event handling

func (t *Terminal) eventLoop(pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopChan:
			return
		case <-ticker.C:
			ev := t.screen.PollEvent()
			if ev == nil {
				continue
			}

			go func(ev tcell.Event) {
				t.lock.Lock()
				defer t.lock.Unlock()

				switch ev := ev.(type) {
				case *tcell.EventResize:
					width, height := ev.Size()
					t.size = core.Size{Width: width, Height: height}
					t.screen.Sync()
					if t.onResize != nil {
						t.onResize(t.size)
					}
				case *tcell.EventMouse:
					t.handleMouse(ev)
				case *tcell.EventKey:
					t.handleKey(ev)
				case *tcell.EventFocus:
					t.focused = ev.Focused
					if t.onFocusChange != nil {
						t.onFocusChange(t.focused)
					}
				}
			}(ev)
		}
	}
}

func (t *Terminal) handleMouse(ev *tcell.EventMouse) {
	// Skip if mouse events are disabled
	if t.mouseMode == MouseDisabled {
		return
	}

	x, y := ev.Position()
	buttons := ev.Buttons()

	// Create mouse event
	event := MouseEvent{
		Buttons:   buttons,
		Position:  core.Point{X: x, Y: y},
		timestamp: ev.When(),
	}

	// Handle based on mouse mode
	switch t.mouseMode {
	case MouseClick:
		// Only send button press events (Primary, Secondary, Middle)
		if buttons&(tcell.ButtonPrimary|tcell.ButtonSecondary|tcell.ButtonMiddle) != 0 {
			t.eventChan <- event
		}
	case MouseDrag:
		// Send button events and drag events
		if buttons != tcell.ButtonNone {
			t.eventChan <- event
		}
	case MouseMotion:
		// Send all mouse events
		t.eventChan <- event
	}
}

func (t *Terminal) handleKey(ev *tcell.EventKey) {
	// Create key event
	event := KeyEvent{
		Key:       ev.Key(),
		Rune:      ev.Rune(),
		Modifiers: ev.Modifiers(),
		timestamp: ev.When(),
	}

	// Send the event through the channel
	t.eventChan <- event
}

// Clipboard operations

// SetClipboard sets the clipboard content
func (t *Terminal) SetClipboard(content string) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	// Try system clipboard first
	if t.clipboardProvider == nil {
		t.clipboardProvider = &SystemClipboard{}
	}

	if err := t.clipboardProvider.Set(content); err != nil {
		// Fall back to in-memory clipboard
		fallback := &FallbackClipboard{}
		t.clipboardProvider = fallback
		return fallback.Set(content)
	}
	return nil
}

// GetClipboard retrieves the clipboard content
func (t *Terminal) GetClipboard() (string, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if t.clipboardProvider == nil {
		t.clipboardProvider = &SystemClipboard{}
	}

	content, err := t.clipboardProvider.Get()
	if err != nil {
		// Fall back to in-memory clipboard
		fallback := &FallbackClipboard{}
		t.clipboardProvider = fallback
		return fallback.Get()
	}
	return content, nil
}

// Callbacks

func (t *Terminal) OnResize(callback func(core.Size)) {
	t.onResize = callback
}

func (t *Terminal) OnFocusChange(callback func(bool)) {
	t.onFocusChange = callback
}

func (t *Terminal) OnSuspend(callback func()) {
	t.onSuspend = callback
}

func (t *Terminal) OnResume(callback func()) {
	t.onResume = callback
}

// Add this new method
func (t *Terminal) HandleEvents(handler func(Event) bool) {
	for {
		select {
		case <-t.stopChan:
			return
		case event := <-t.eventChan:
			if handler(event) {
				return
			}
		}
	}
}

func (t *Terminal) EnableUnicode() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.unicodeMode = true
}

func (t *Terminal) DisableUnicode() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.unicodeMode = false
}

func (t *Terminal) EnableCombiningChars() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.combiningChars = true
}

func (t *Terminal) DisableCombiningChars() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.combiningChars = false
}

// SetTitle sets the terminal window title
func (t *Terminal) SetTitle(title string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.title = title
	// t.screen.SetTitle(title) // FIXME: Seems broken on this version of tcell
}

// GetTitle returns the current terminal window title
func (t *Terminal) GetTitle() string {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.title
}

// SwapBuffers swaps the front and back buffers
func (t *Terminal) SwapBuffers() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.frontBuffer, t.backBuffer = t.backBuffer, t.frontBuffer
}

// DrawText draws a string of text, handling combining characters appropriately
func (t *Terminal) DrawText(x, y int, text string, fg, bg core.Color, style StyleMask) {
	currentX := x
	for _, ch := range text {
		t.DrawStyledCell(currentX, y, ch, fg, bg, style)
		if !t.combiningChars || !unicode.IsMark(ch) {
			currentX++
		}
	}
}
