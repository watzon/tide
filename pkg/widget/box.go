package widget

import "github.com/watzon/tide/pkg/engine"

// Box is a basic container widget
type Box struct {
	BaseWidget
	children []Widget
}

func NewBox() *Box {
	return &Box{
		children: make([]Widget, 0),
	}
}

func (b *Box) AppendChild(child Widget) {
	b.children = append(b.children, child)
}

func (b *Box) Build(context BuildContext) Widget {
	return b
}

func (b *Box) CreateRenderObject() RenderObject {
	box := NewBaseRenderBox()
	box.WithStyle(b.GetStyle())

	// Create render objects for children
	for _, child := range b.children {
		childRenderObj := child.CreateRenderObject()
		box.AppendChild(childRenderObj)
	}

	return box
}

func (b *Box) UpdateRenderObject(renderObject RenderObject) {
	if box, ok := renderObject.(*BaseRenderBox); ok {
		box.WithStyle(b.GetStyle())

		// Update children's render objects
		for i, child := range b.children {
			if i < len(box.Children()) {
				child.UpdateRenderObject(box.Children()[i])
			}
		}
	}
}

// Make sure BaseRenderBox properly paints children
func (r *BaseRenderBox) Paint(context engine.RenderContext) {
	r.PaintBackground(context)
	r.PaintBorder(context)
	r.PaintContent(context) // This should call Paint on all children
}

func (r *BaseRenderBox) PaintContent(context engine.RenderContext) {
	// Get content area
	contentRect := r.ContentRect()

	// Push offset for children
	context.PushOffset(contentRect.Min)

	// Paint each child
	for _, child := range r.children {
		child.Paint(context)
	}

	// Pop offset after painting children
	context.PopOffset()
}
