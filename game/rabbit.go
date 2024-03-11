package game

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/demonodojo/rabbits/assets"
	"github.com/google/uuid"
)

type Rabbit struct {
	Serial
	game           *Game
	Position       Vector  `json:"position"`
	Rotation       float64 `json:"rotation"`
	Speed          float64 `json:"speed"`
	Score          int32   `json:"score"`
	Heat           int64   `json:"heat"`
	Load           int64   `json:"load"`
	lastUpdateTime time.Time
	scale          float64
	bounds         image.Rectangle
	halfW          float64
	halfH          float64
	sprite         *ebiten.Image
	spriteR        *ebiten.Image

	shootCooldown *Timer
}

func NewRabbit(game *Game) *Rabbit {
	sprite := assets.PlayerSprite
	spriteR := assets.RabbitSpriteR

	scale := .4
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
			Action:    "Spawn",
		},
		game:           game,
		Position:       pos,
		Score:          0,
		Heat:           0,
		Load:           0,
		scale:          scale,
		bounds:         bounds,
		halfW:          halfW,
		halfH:          halfH,
		Rotation:       0,
		sprite:         sprite,
		spriteR:        spriteR,
		lastUpdateTime: time.Now(),
	}
}

func (r *Rabbit) Update() {
	now := time.Now()
	delta := now.Sub(r.lastUpdateTime)
	r.lastUpdateTime = now

	deltaMs := delta.Seconds() * 1000
	updateFactor := deltaMs / 16.666
	r.Position.X += math.Sin(r.Rotation) * updateFactor * r.Speed
	r.Position.Y += math.Cos(r.Rotation) * updateFactor * (-r.Speed)
	if r.Heat > 0 {
		r.Heat--
	}
	if r.Load > 0 {
		r.Load--
	}
}

func (r *Rabbit) Interact() bool {
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

	if ebiten.IsKeyPressed(ebiten.KeyF) {
		if r.Heat < 100 && r.Load == 0 {
			r.Action = "FIRE"
			r.Load = 30
			r.Heat += 30
			interaction = true
		}
	}

	r.Update()

	return interaction
}

func (r *Rabbit) advancedPosition() (Vector, float64) {
	position := Vector{
		X: r.Position.X + 19.0 + math.Sin(r.Rotation)*40,
		Y: r.Position.Y + 15.0 - math.Cos(r.Rotation)*40,
	}
	return position, r.Rotation
}

func (r *Rabbit) Draw(screen *ebiten.Image, geom ebiten.GeoM) {

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(r.scale, r.scale)
	op.GeoM.Translate(-r.halfW, -r.halfH)
	op.GeoM.Rotate(r.Rotation)
	op.GeoM.Translate(r.halfW, r.halfH)

	x := r.Position.X
	y := r.Position.Y
	op.GeoM.Translate(x, y)
	op.GeoM.Concat(geom)

	screen.DrawImage(r.sprite, op)
	text.Draw(screen, fmt.Sprintf("%f %f", r.halfH, r.halfW), assets.InfoFont, 10, 70, color.White)
	text.Draw(screen, fmt.Sprintf("Heat %d", r.Heat), assets.InfoFont, 10, 90, color.White)
}

func (r *Rabbit) Fired() {
	r.Speed = 0
	r.Score -= 10
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
	r.Action = other.Action
	r.Position = other.Position
	r.Rotation = other.Rotation
	r.Speed = other.Speed
	r.Score = other.Score
}
