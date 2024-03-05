package game

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/math/f64"
)

type Camera struct {
	ViewPort   f64.Vec2
	Position   f64.Vec2
	ZoomFactor float64
	Rotation   int
	Matrix     ebiten.GeoM
	attached   bool
	dragStart  f64.Vec2
	posStart   f64.Vec2
	dragging   bool
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %d, S: %f",
		c.Position, c.Rotation, c.ZoomFactor,
	)
}

func (c *Camera) viewportCenter() f64.Vec2 {
	return f64.Vec2{
		c.ViewPort[0] * 0.5,
		c.ViewPort[1] * 0.5,
	}
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position[0], -c.Position[1])
	// We want to scale and rotate around center of image / screen
	m.Translate(-c.viewportCenter()[0], -c.viewportCenter()[1])
	m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
	return m
}

func (c *Camera) Render(world, screen *ebiten.Image) {
	screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: c.worldMatrix(),
	})
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		// When scaling it can happened that matrix is not invertable
		return math.NaN(), math.NaN()
	}
}

func (c *Camera) Reset() {
	c.Position[0] = 0
	c.Position[1] = 0
	c.Rotation = 0
	c.ZoomFactor = 0
	c.attached = false
	c.dragging = false

}

func (c *Camera) Update(r *Rabbit) error {
	if c.attached {
		c.Position[0] = r.Position.X + r.halfW - c.viewportCenter()[0]
		c.Position[1] = r.Position.Y + r.halfH - c.viewportCenter()[1]
		c.viewportCenter()
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		c.Position[0] -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		c.Position[0] += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		c.Position[1] -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		c.Position[1] += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		if c.ZoomFactor > -2400 {
			c.ZoomFactor -= 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		if c.ZoomFactor < 2400 {
			c.ZoomFactor += 1
		}
	}

	_, y := ebiten.Wheel()
	if y != 0 {
		c.ZoomFactor += y * 5
	}

	if ebiten.IsKeyPressed(ebiten.KeyE) {
		c.Rotation -= 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		c.Rotation += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		c.Reset()
	}

	if ebiten.IsKeyPressed(ebiten.KeyO) {
		c.attached = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		c.attached = false
	}

	c.updateDrag()

	c.Matrix = c.worldMatrix()
	return nil
}

func (c *Camera) updateDrag() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		c.dragStart = f64.Vec2{
			float64(x),
			float64(y),
		}
		c.posStart = c.Position
		c.dragging = true
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		c.dragging = false
	}
	if c.dragging {
		x, y := ebiten.CursorPosition()
		c.Position[0] = c.posStart[0] - float64(x) + c.dragStart[0]
		c.Position[1] = c.posStart[1] - float64(y) + c.dragStart[1]
	}
}
