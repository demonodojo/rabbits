package game

import (
	"encoding/json"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/demonodojo/rabbits/assets"
	"github.com/google/uuid"
)

type Lettuce struct {
	Serial
	scale    float64
	Position Vector
	sprite   *ebiten.Image
}

func NewLettuce() *Lettuce {
	sprite := assets.LettuceSprite
	bounds := sprite.Bounds()
	scale := 4.0
	pos := Vector{
		X: float64(screenWidth-bounds.Dx()*4) * rand.Float64(),
		Y: float64(screenHeight-bounds.Dy()*4) * rand.Float64(),
	}

	l := &Lettuce{
		Serial: Serial{
			ID:        uuid.New(),
			ClassName: "Lettuce",
			Action:    "Spawn",
		},
		Position: pos,
		scale:    scale,
		sprite:   sprite,
	}
	return l
}

func (m *Lettuce) Update() {

}

func (l *Lettuce) Draw(screen *ebiten.Image, geom ebiten.GeoM) {

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(l.scale, l.scale)

	op.GeoM.Translate(l.Position.X, l.Position.Y)
	op.GeoM.Concat(geom)

	screen.DrawImage(l.sprite, op)
}

func (l *Lettuce) Collider() Rect {
	bounds := l.sprite.Bounds()

	return NewRect(
		l.Position.X,
		l.Position.Y,
		float64(bounds.Dx())*l.scale,
		float64(bounds.Dy())*l.scale,
	)
}

func (r *Lettuce) ToJson() string {
	json, _ := json.Marshal(r)
	return string(json)
}

func (r *Lettuce) CopyFrom(other *Lettuce) {
	r.ID = other.ID
	r.Action = other.Action
	r.Position = other.Position
}
