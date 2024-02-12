package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/demonodojo/rabbits/assets"
)

type Lettuce struct {
	scale    float64
	position Vector
	sprite   *ebiten.Image
}

func NewLettuce(baseVelocity float64) *Lettuce {
	sprite := assets.LettuceSprite
	bounds := sprite.Bounds()
	scale := 4.0
	pos := Vector{
		X: float64(screenWidth-bounds.Dx()*4) * rand.Float64(),
		Y: float64(screenHeight-bounds.Dy()*4) * rand.Float64(),
	}

	l := &Lettuce{
		position: pos,
		scale:    scale,
		sprite:   sprite,
	}
	return l
}

func (m *Lettuce) Update() {

}

func (l *Lettuce) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(l.scale, l.scale)

	op.GeoM.Translate(l.position.X, l.position.Y)

	screen.DrawImage(l.sprite, op)
}

func (l *Lettuce) Collider() Rect {
	bounds := l.sprite.Bounds()

	return NewRect(
		l.position.X,
		l.position.Y,
		float64(bounds.Dx())*l.scale,
		float64(bounds.Dy())*l.scale,
	)
}
