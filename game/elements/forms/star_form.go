package forms

import (
	"encoding/json"

	"image/color"
	"log"

	"golang.org/x/image/font/gofont/goregular"

	"github.com/demonodojo/rabbits/game"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/google/uuid"
)

type StarForm struct {
	game.Serial
	scale float64
	ui    *ebitenui.UI
	root  *widget.Container
}

func NewStarForm() *StarForm {

	// This creates the root container for this UI.
	// All other UI elements must be added to this container.
	rootContainer := widget.NewContainer(
		// The container will use a plain color as its background. This is not required if you wish
		// the container to be transparent.
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0xaa, 0xff})),
		// Containers have the concept of a Layout. This is how children of this container should be
		// displayed within the bounds of this container.
		// The container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// This adds the root container to the UI, so that it will be rendered.
	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	// This loads a font and creates a font face.
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal("Error Parsing Font", err)
	}
	fontFace := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 32,
	})

	// This creates a text widget that says "Hello World!"
	helloWorldLabel := widget.NewText(
		widget.TextOpts.Text("Hello World!", fontFace, color.White),
	)

	// To display the text widget, we have to add it to the root container.
	rootContainer.AddChild(helloWorldLabel)

	l := &StarForm{
		Serial: game.Serial{
			ID:        uuid.New(),
			ClassName: "StarForm",
			Action:    "Spawn",
		},
		ui:   eui,
		root: rootContainer,
	}
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
