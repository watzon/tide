// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package engine

// Renderer manages the rendering pipeline
type Renderer struct {
	backend    Backend
	compositor *Compositor
}

func NewRenderer(backend Backend) *Renderer {
	return &Renderer{
		backend:    backend,
		compositor: NewCompositor(),
	}
}

func (r *Renderer) Render() error {
	r.backend.Clear()
	r.compositor.Compose(r.backend)
	return r.backend.Present()
}
