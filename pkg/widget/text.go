// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package widget

import (
	"strings"

	"github.com/watzon/tide/pkg/core/geometry"
	"github.com/watzon/tide/pkg/engine"
)

// Text is a widget that displays text
type Text struct {
	BaseWidget // This gives us GetConstraints, GetStyle, etc.
	content    string
}

func NewText(content string) *Text {
	return &Text{
		content: content,
		BaseWidget: BaseWidget{
			style: NewWidgetStyle(), // Initialize with default style
		},
	}
}

func (t *Text) Build(context BuildContext) Widget {
	return t
}

// TextRenderObject handles rendering of text content
type TextRenderObject struct {
	BaseRenderObject
	content string
}

func NewTextRenderObject(style WidgetStyle, content string) *TextRenderObject {
	return &TextRenderObject{
		BaseRenderObject: BaseRenderObject{
			style: style,
		},
		content: content,
	}
}

func (r *TextRenderObject) Paint(context engine.RenderContext) {
	// Paint background using BaseRenderObject's functionality
	r.BaseRenderObject.Paint(context)

	// Split content into lines
	lines := strings.Split(r.content, "\n")

	// Paint each line
	for y, line := range lines {
		if y >= r.size.Height {
			break // Don't exceed height
		}
		for x, ch := range line {
			if x >= r.size.Width {
				break // Don't exceed width
			}
			context.DrawCell(x, y, ch, r.style.ForegroundColor, r.style.BackgroundColor)
		}
	}
}

func (r *TextRenderObject) Layout(constraints Constraints) geometry.Size {
	lines := strings.Split(r.content, "\n")

	// Calculate required size
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	height := len(lines)

	// Apply constraints
	r.size = constraints.Constrain(geometry.Size{
		Width:  width,
		Height: height,
	})

	return r.size
}

func (t *Text) CreateRenderObject() RenderObject {
	return NewTextRenderObject(t.GetStyle(), t.content)
}

func (t *Text) UpdateRenderObject(renderObject RenderObject) {
	if textRenderObj, ok := renderObject.(*TextRenderObject); ok {
		textRenderObj.style = t.GetStyle()
		textRenderObj.content = t.content
	}
}

func (t *Text) GetContent() string {
	return t.content
}

func (t *Text) WithContent(content string) *Text {
	t.content = content
	return t
}
