package game

import (
	"encoding/json"

	"image/color"
	"log"
	"math/rand"

	"golang.org/x/image/font/gofont/goregular"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/google/uuid"
)

type Form struct {
	Serial
	scale    float64
	Position Vector
	ui       *ebitenui.UI
	root     *widget.Container
}

func NewForm() *Form {

	pos := Vector{
		X: (float64(screenWidth) * rand.Float64()) - screenWidth*1.5,
		Y: (float64(screenHeight) * rand.Float64()) - screenHeight*1.5,
	}

	// This creates the root container for this UI.
	// All other UI elements must be added to this container.
	rootContainer := widget.NewContainer()

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

	l := &Form{
		Serial: Serial{
			ID:        uuid.New(),
			ClassName: "Form",
			Action:    "Spawn",
		},
		Position: pos,
		ui:       eui,
		root:     rootContainer,
	}
	return l
}

func (f *Form) Update() {
	f.ui.Update()
}

func (f *Form) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	f.ui.Draw(screen)
}

func (r *Form) ToJson() string {
	json, _ := json.Marshal(r)
	return string(json)
}

func (r *Form) CopyFrom(other *Form) {
	r.ID = other.ID
	r.Action = other.Action
	r.Position = other.Position
}
