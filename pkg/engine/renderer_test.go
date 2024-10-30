// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package engine_test

import (
	"testing"

	"github.com/watzon/tide/pkg/core/color"
	"github.com/watzon/tide/pkg/core/geometry"
	"github.com/watzon/tide/pkg/engine"
)

type testBackend struct {
	cleared bool
	present bool
	size    geometry.Size
}

func newTestBackend() *testBackend {
	return &testBackend{size: geometry.Size{Width: 80, Height: 24}}
}

func (b *testBackend) Init() error                                    { return nil }
func (b *testBackend) Shutdown() error                                { return nil }
func (b *testBackend) Size() geometry.Size                            { return b.size }
func (b *testBackend) Clear()                                         { b.cleared = true }
func (b *testBackend) DrawCell(x, y int, ch rune, fg, bg color.Color) {}
func (b *testBackend) Present() error                                 { b.present = true; return nil }

func TestRenderer(t *testing.T) {
	t.Run("Render cycle", func(t *testing.T) {
		backend := newTestBackend()
		renderer := engine.NewRenderer(backend)

		err := renderer.Render()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !backend.cleared {
			t.Error("expected backend to be cleared")
		}

		if !backend.present {
			t.Error("expected backend to be presented")
		}
	})
}
