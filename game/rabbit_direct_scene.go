package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/demonodojo/rabbits/assets"
)

const (
	screenWidth  = 800
	screenHeight = 600

	lettuceSpawnTime = 1 * time.Second

	baseMeteorVelocity  = 0.25
	meteorSpeedUpAmount = 0.1
	meteorSpeedUpTime   = 5 * time.Second
)

type RabbitDirectScene struct {
	game              *Game
	player            *Player
	rabbit            *Rabbit
	lettuceSpawnTimer *Timer
	lettuces          []*Lettuce
	bullets           []*Bullet

	score         int
	scale         float64
	baseVelocity  float64
	velocityTimer *Timer
}

func NewRabbitDirectScene(g *Game) *RabbitDirectScene {
	s := &RabbitDirectScene{
		game:              g,
		lettuceSpawnTimer: NewTimer(lettuceSpawnTime),
		baseVelocity:      baseMeteorVelocity,
		velocityTimer:     NewTimer(meteorSpeedUpTime),
	}

	s.rabbit = NewRabbit(g)

	m := NewLettuce(s.baseVelocity)
	s.lettuces = append(s.lettuces, m)

	return s
}

func (g *RabbitDirectScene) Update() error {

	g.scale += 0.01
	if g.scale > 2 {
		g.scale = 0.5
	}

	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += meteorSpeedUpAmount
	}

	g.rabbit.Update()

	g.lettuceSpawnTimer.Update()
	if g.lettuceSpawnTimer.IsReady() {
		g.lettuceSpawnTimer.Reset()

		m := NewLettuce(g.baseVelocity)
		g.lettuces = append(g.lettuces, m)
	}

	for _, l := range g.lettuces {
		l.Update()
	}

	if ebiten.IsKeyPressed(ebiten.Key1) {
		ebiten.SetFullscreen(true)
	}

	if ebiten.IsKeyPressed(ebiten.Key2) {
		ebiten.SetFullscreen(false)
	}

	// for _, b := range g.bullets {
	// 	b.Update()
	// }

	// // Check for meteor/bullet collisions
	// for i, m := range g.meteors {
	// 	for j, b := range g.bullets {
	// 		if m.Collider().Intersects(b.Collider()) {
	// 			g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
	// 			g.bullets = append(g.bullets[:j], g.bullets[j+1:]...)
	// 			g.score++
	// 		}
	// 	}
	// }

	// Check for rabbit/lettuces collisions
	for i, l := range g.lettuces {
		if l.Collider().Intersects(g.rabbit.Collider()) {
			g.lettuces = append(g.lettuces[:i], g.lettuces[i+1:]...)
			g.score++
			if len(g.lettuces) == 0 {
				g.Reset()
			}
			//g.Reset()
			break
		}
	}

	return nil
}

func (g *RabbitDirectScene) Draw(screen *ebiten.Image) {

	opts := &ebiten.DrawImageOptions{}
	// Ajusta la escala de la imagen. 1 es el tamaño original, valores mayores para hacer zoom in, menores para zoom out.
	opts.GeoM.Scale(g.scale, g.scale)
	// Dibuja la imagen en la pantalla con las opciones de escala.

	g.rabbit.Draw(screen)

	for _, m := range g.lettuces {
		m.Draw(screen)
	}

	// for _, b := range g.bullets {
	// 	b.Draw(screen)
	// }

	text.Draw(screen, fmt.Sprintf("%06d", g.score), assets.ScoreFont, screenWidth/2-100, 50, color.White)
}

func (g *RabbitDirectScene) Reset() {
	g.rabbit = NewRabbit(g.game)
	g.lettuces = nil
	g.bullets = nil
	g.score = 0
	g.lettuceSpawnTimer.Reset()
	g.baseVelocity = baseMeteorVelocity
	g.velocityTimer.Reset()
}

func (g *RabbitDirectScene) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (s *RabbitDirectScene) SpawnElement(name string, element interface{}) {

}
