package painter

import (
	"github.com/roman-mazur/architecture-lab-3/ui"
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	Update(state *TextureState)
}

type Fill struct {
	Color color.Color
}

type Update struct{}

type Reset struct{}

type TextureOperation interface {
	Do(t screen.Texture)
	Update(state *TextureState)
}

type BgRect struct {
	X1 float32
	Y1 float32
	X2 float32
	Y2 float32
}

type Figure struct {
	X float32
	Y float32
}

type Move struct {
	X float32
	Y float32
}

type OperationList []Operation

var UpdateOp = Update{}

func (op Update) Update(_ *TextureState) {}

func (op Fill) Do(t screen.Texture) {
	t.Fill(t.Bounds(), op.Color, screen.Src)
}

func (op Fill) Update(state *TextureState) {
	state.backgroundColor = &op
}

var ResetOp = Reset{}

func (op Reset) Update(state *TextureState) {
	state.backgroundColor = &Fill{Color: color.Black}
	state.backgroundRect = nil
	state.figureCenters = nil
}

func (op BgRect) Do(t screen.Texture) {
	t.Fill(
		image.Rect(
			int(op.X1*float32(t.Size().X)),
			int(op.Y1*float32(t.Size().Y)),
			int(op.X2*float32(t.Size().X)),
			int(op.Y2*float32(t.Size().Y)),
		),
		color.Black,
		screen.Src,
	)
}

func (op BgRect) Update(state *TextureState) {
	state.backgroundRect = &op
}

func (op Figure) Do(t screen.Texture) {
	ui.DrawFigure(
		t,
		image.Pt(
			int(op.X*float32(t.Size().X)),
			int(op.Y*float32(t.Size().Y)),
		),
	)
}

func (op Figure) Update(state *TextureState) {
	state.figureCenters = append(state.figureCenters, &op)
}

func (op Move) Update(state *TextureState) {
	for _, fig := range state.figureCenters {
		fig.X += op.X
		fig.Y += op.Y
	}
}
