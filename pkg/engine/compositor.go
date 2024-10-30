// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package engine

import (
	"sort"

	"github.com/watzon/tide/pkg/core/geometry"
)

type Layer struct {
	Bounds geometry.Rect
	Z      int
	Draw   func(b Backend)
}

type Compositor struct {
	layers []Layer
}

func NewCompositor() *Compositor {
	return &Compositor{
		layers: make([]Layer, 0),
	}
}

func (c *Compositor) AddLayer(layer Layer) {
	c.layers = append(c.layers, layer)
}

func (c *Compositor) Compose(backend Backend) {
	// Sort layers by Z-index (lower Z-index drawn first)
	sortedLayers := make([]Layer, len(c.layers))
	copy(sortedLayers, c.layers)

	sort.Slice(sortedLayers, func(i, j int) bool {
		return sortedLayers[i].Z < sortedLayers[j].Z
	})

	// Draw layers in order
	for _, layer := range sortedLayers {
		layer.Draw(backend)
	}
}
