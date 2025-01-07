package game

import (
	"encoding/json"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/demonodojo/rabbits/assets"
)

const (
	bulletSpeedPerSecond = 400.0
	bulletLife           = 60
)

type Bullet struct {
	Serial
	Position       Vector
	Rotation       float64
	Life           int
	sprite         *ebiten.Image
	lastUpdateTime time.Time
	scale          float64
}

func NewBullet(pos Vector, rotation float64) *Bullet {
	sprite := assets.LaserSprite

	scale := .2
	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) * scale / 2
	halfH := float64(bounds.Dy()) * scale / 2

	pos.X -= halfW
	pos.Y -= halfH

	b := &Bullet{
		Serial: Serial{
			ID:        uuid.New(),
			ClassName: "Bullet",
			Action:    "Spawn",
		},
		Position:       pos,
		Rotation:       rotation,
		sprite:         sprite,
		lastUpdateTime: time.Now(),
		scale:          scale,
		Life:           bulletLife,
	}

	return b
}

func (b *Bullet) Update() {

	now := time.Now()
	delta := now.Sub(b.lastUpdateTime)
	b.lastUpdateTime = now

	deltaMs := delta.Seconds() * 1000
	updateFactor := deltaMs / 16.666

	speed := bulletSpeedPerSecond / float64(ebiten.TPS())

	b.Position.X += math.Sin(b.Rotation) * speed * updateFactor
	b.Position.Y += math.Cos(b.Rotation) * -speed * updateFactor
	b.Life--
	if b.Life <= 0 {
		b.Action = "DELETE"
	}
}

func (b *Bullet) Draw(screen *ebiten.Image, geom ebiten.GeoM) {
	bounds := b.sprite.Bounds()
	halfW := float64(bounds.Dx()) * b.scale / 2
	halfH := float64(bounds.Dy()) * b.scale / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(b.scale, b.scale)
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(b.Rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(b.Position.X, b.Position.Y)

	op.GeoM.Concat(geom)

	screen.DrawImage(b.sprite, op)
}

func (b *Bullet) Collider() Rect {
	bounds := b.sprite.Bounds()

	return NewRect(
		b.Position.X,
		b.Position.Y,
		float64(bounds.Dx())*b.scale,
		float64(bounds.Dy())*b.scale,
	)
}

func (b *Bullet) ToJson() string {
	json, _ := json.Marshal(b)
	return string(json)
}

func (b *Bullet) CopyFrom(other *Bullet) {
	b.ID = other.ID
	b.Action = other.Action
	b.Position = other.Position
	b.Rotation = other.Rotation
	b.Life = other.Life

}
