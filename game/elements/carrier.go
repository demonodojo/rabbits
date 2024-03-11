package elements

import (
	"encoding/json"
	"image"
	"image/color"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/demonodojo/rabbits/assets"
	"github.com/demonodojo/rabbits/game"
	"github.com/demonodojo/rabbits/game/network"
	"github.com/google/uuid"
)

type Carrier struct {
	game.Serial
	scale    float64
	Position game.Vector
	Rotation float64 `json:"rotation"`
	Speed    float64 `json:"speed"`
	Score    int32   `json:"score"`
	//lastUpdateTime time.Time,
	sprite   *ebiten.Image
	Name     string
	size     int
	selected bool
	bounds   image.Rectangle
	halfW    float64
	halfH    float64
	color    color.Color
	channel  *network.MessageQueue
}

func NewCarrier(channel *network.MessageQueue, position game.Vector) *Carrier {

	colors := []color.Color{
		color.RGBA{0, 0, 255, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 0, 255, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{0, 255, 255, 255},
		color.RGBA{0, 0, 0, 255},
		color.RGBA{125, 255, 255, 255},
	}
	sprite := assets.CarrierSprite
	scale := 1.0

	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) * scale / 2
	halfH := float64(bounds.Dy()) * scale / 2

	l := &Carrier{
		Serial: game.Serial{
			ID:        uuid.New(),
			ClassName: "Carrier",
			Action:    "Spawn",
		},
		Position: position,
		Rotation: 0,
		Speed:    0,
		Score:    0,
		scale:    scale,
		sprite:   sprite,
		bounds:   bounds,
		halfW:    halfW,
		halfH:    halfH,
		color:    colors[rand.Intn(8)],
		size:     rand.Intn(7) + 1,
		channel:  channel,
		//		lastUpdateTime: time.Now(),
	}
	return l
}

func (c *Carrier) Update() {
	//now := time.Now()
	//delta := now.Sub(c.lastUpdateTime)
	//c.lastUpdateTime = now

	//deltaMs := delta.Seconds() * 1000
	//updateFactor := deltaMs / 16.666
	//c.Position.X += math.Sin(c.Rotation) * updateFactor * c.Speed
	//c.Position.Y += math.Cos(c.Rotation) * updateFactor * (-c.Speed)
}

func (c *Carrier) Draw(screen *ebiten.Image, geom ebiten.GeoM) {

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(c.scale, c.scale)
	op.GeoM.Translate(-c.halfW, -c.halfH)
	op.GeoM.Rotate(c.Rotation)
	op.GeoM.Translate(c.halfW, c.halfH)

	x := c.Position.X
	y := c.Position.Y
	op.GeoM.Translate(x, y)
	op.GeoM.Concat(geom)
}

func (l *Carrier) Select() {
	l.selected = true

}

func (l *Carrier) Toggle() {
	l.selected = !l.selected
}

func (l *Carrier) UnSelect() {
	l.selected = false
}

func (l *Carrier) Edit() {
	l.Action = "EDIT"
	l.channel.Enqueue(l.ToJson())
}

func (l *Carrier) Center() (float64, float64) {
	return l.Position.X + l.halfW, l.Position.Y + l.halfH
}

func (l *Carrier) Collider() game.Rect {

	return game.NewRect(
		l.Position.X-8,
		l.Position.Y-8,
		16,
		16,
	)
}

func (r *Carrier) ToJson() string {
	json, _ := json.Marshal(r)
	return string(json)
}

func (r *Carrier) CopyFrom(other *Carrier) {
	r.ID = other.ID
	r.Action = other.Action
	r.Position = other.Position
}
