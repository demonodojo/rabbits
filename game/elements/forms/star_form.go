package forms

import (
	"encoding/json"

	"image/color"

	"github.com/demonodojo/rabbits/game"
	"github.com/demonodojo/rabbits/game/elements"
	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/google/uuid"
)

type StarForm struct {
	game.Serial
	scale float64
	ui    *ebitenui.UI
	root  *widget.Container
	star  *elements.Star
}

func NewStarForm(star *elements.Star) *StarForm {

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use a row layout to layout the textinput widgets
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)))),
	)

	// This adds the root container to the UI, so that it will be rendered.
	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	l := &StarForm{
		Serial: game.Serial{
			ID:        uuid.New(),
			ClassName: "StarForm",
			Action:    "Spawn",
		},
		ui:   eui,
		root: rootContainer,
		star: star,
	}

	name := NewTextInput("Name", star.Name)
	// add the button as a child of the container
	rootContainer.AddChild(NewButton("save", func(args *widget.ButtonClickedEventArgs) {
		l.Action = "SUBMIT"
		star.Name = name.GetText()
	}))
	rootContainer.AddChild(name)

	return l
}

func (f *StarForm) Update() {
	f.ui.Update()
}

func (f *StarForm) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	f.ui.Draw(screen)
}

func (r *StarForm) ToJson() string {
	json, _ := json.Marshal(r)
	return string(json)
}

func (r *StarForm) CopyFrom(other *StarForm) {
	r.ID = other.ID
	r.Action = other.Action
}
