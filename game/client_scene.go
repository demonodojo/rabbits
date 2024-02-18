package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/demonodojo/rabbits/assets"
	"github.com/demonodojo/rabbits/game/network"
)

type ClientScene struct {
	game     *Game
	client   *network.Client
	rabbit   *Rabbit
	lettuces []*Lettuce

	score         int
	scale         float64
	baseVelocity  float64
	velocityTimer *Timer
}

func NewClientScene(g *Game, client *network.Client) *ClientScene {
	s := &ClientScene{
		game:          g,
		client:        client,
		baseVelocity:  baseMeteorVelocity,
		velocityTimer: NewTimer(meteorSpeedUpTime),
	}

	s.rabbit = NewRabbit(g)
	client.Write(s.rabbit.ToJson())
	return s
}

func (s *ClientScene) Update() error {

	if s.rabbit.Update() {
		s.client.Write(s.rabbit.ToJson())
	}

	for _, l := range s.lettuces {
		l.Update()
	}

	if ebiten.IsKeyPressed(ebiten.Key1) {
		ebiten.SetFullscreen(true)
	}

	if ebiten.IsKeyPressed(ebiten.Key2) {
		ebiten.SetFullscreen(false)
	}

	return nil
}

func (g *ClientScene) Draw(screen *ebiten.Image) {

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

func (g *ClientScene) Reset() {
	g.rabbit = NewRabbit(g.game)
	g.lettuces = nil
	g.score = 0
	g.baseVelocity = baseMeteorVelocity
	g.velocityTimer.Reset()
}

func (g *ClientScene) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (s *ClientScene) SpawnElement(name string, element interface{}) {

}
