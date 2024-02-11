package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/ThreeDotsLabs/meteors/assets"
)

type Rabbit struct {
	game *Game

	position Vector
	rotation float64
	speed    float64
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
		game:          game,
		position:      pos,
		scale:         scale,
		bounds:        bounds,
		halfW:         halfW,
		halfH:         halfH,
		rotation:      0,
		sprite:        sprite,
		spriteR:       spriteR,
		shootCooldown: NewTimer(shootCooldown),
	}
}

func (r *Rabbit) Update() {
	rotationSpeed := rotationPerSecond / float64(ebiten.TPS())
	speedPerSecond := 0.1

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		r.rotation -= rotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		r.rotation += rotationSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		r.speed += speedPerSecond
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		r.speed -= speedPerSecond
	}

	r.position.X += math.Sin(r.rotation) * r.speed
	r.position.Y += math.Cos(r.rotation) * -r.speed

	r.shootCooldown.Update()

}

func (r *Rabbit) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(r.scale, r.scale)
	op.GeoM.Translate(-r.halfW, -r.halfH)
	op.GeoM.Rotate(r.rotation)
	op.GeoM.Translate(r.halfW, r.halfH)

	x := r.position.X
	y := r.position.Y
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
	if x != r.position.X || y != r.position.Y {
		sprite = r.spriteR
	}

	op.GeoM.Translate(x, y)

	screen.DrawImage(sprite, op)
}

func (r *Rabbit) Collider() Rect {

	return NewRect(
		r.position.X,
		r.position.Y,
		float64(r.halfW*2),
		float64(r.halfH*2),
	)
}
