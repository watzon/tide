// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package terminal

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/watzon/tide/pkg/core"
)

// Event types for the terminal
type KeyEvent struct {
	Key       tcell.Key
	Rune      rune
	Modifiers tcell.ModMask
	timestamp time.Time
}

type MouseEvent struct {
	Buttons   tcell.ButtonMask
	Position  core.Point
	timestamp time.Time
}

func (e KeyEvent) When() time.Time   { return e.timestamp }
func (e MouseEvent) When() time.Time { return e.timestamp }
