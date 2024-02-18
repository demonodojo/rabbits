package game

import (
	"encoding/json"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/demonodojo/rabbits/assets"
	"github.com/google/uuid"
)

type Rabbit struct {
	Serial
	game     *Game
	Position Vector  `json:"position"`
	Rotation float64 `json:"rotation"`
	Speed    float64 `json:"speed"`
	scale    float64
	bounds   image.Rectangle
	halfW    float64
	halfH    float64
	sprite   *ebiten.Image
	spriteR  *ebiten.Image

	shootCooldown *Timer
}

func NewRabbit(game *Game) *Rabbit {
	sprite := assets.RabbitSprite
	spriteR := assets.RabbitSpriteR

	scale := 4.0
	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) * scale / 2
	halfH := float64(bounds.Dy()) * scale / 2

	pos := Vector{
		X: screenWidth/2 - halfW,
		Y: screenHeight/2 - halfH,
	}

	return &Rabbit{
		Serial: Serial{
			ID:        uuid.New(),
			ClassName: "Rabbit",
		},
		game:     game,
		Position: pos,
		scale:    scale,
		bounds:   bounds,
		halfW:    halfW,
		halfH:    halfH,
		Rotation: 0,
		sprite:   sprite,
		spriteR:  spriteR,
	}
}

func (r *Rabbit) Update() bool {
	rotationSpeed := rotationPerSecond / float64(ebiten.TPS())
	SpeedPerSecond := 0.1
	interaction := false
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		r.Rotation -= rotationSpeed
		interaction = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		r.Rotation += rotationSpeed
		interaction = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		r.Speed += SpeedPerSecond
		interaction = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		r.Speed -= SpeedPerSecond
		interaction = true
	}

	r.Position.X += math.Sin(r.Rotation) * r.Speed
	r.Position.Y += math.Cos(r.Rotation) * -r.Speed

	return interaction
}

func (r *Rabbit) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(r.scale, r.scale)
	op.GeoM.Translate(-r.halfW, -r.halfH)
	op.GeoM.Rotate(r.Rotation)
	op.GeoM.Translate(r.halfW, r.halfH)

	x := r.Position.X
	y := r.Position.Y
	if x > screenWidth-r.halfW*2 {
		x = screenWidth - r.halfW*2
	} else if x < 0 {
		x = 0
	}
	if y > screenHeight-r.halfH*2 {
		y = screenHeight - r.halfH*2
	} else if y < 0 {
		y = 0
	}

	sprite := r.sprite
	if x != r.Position.X || y != r.Position.Y {
		sprite = r.spriteR
	}

	op.GeoM.Translate(x, y)

	screen.DrawImage(sprite, op)
}

func (r *Rabbit) Collider() Rect {

	return NewRect(
		r.Position.X,
		r.Position.Y,
		float64(r.halfW*2),
		float64(r.halfH*2),
	)
}

func (r *Rabbit) ToJson() string {
	json, _ := json.Marshal(r)
	return string(json)
}

func (r *Rabbit) CopyFrom(other *Rabbit) {
	r.ID = other.ID
	r.Position = other.Position
	r.Rotation = other.Rotation
	r.Speed = other.Speed
}
