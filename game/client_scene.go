package game

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/demonodojo/rabbits/assets"
	"github.com/demonodojo/rabbits/game/network"
)

type ClientScene struct {
	game     *Game
	client   *network.Client
	rabbit   *Rabbit
	rabbits  map[uuid.UUID]*Rabbit
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
	s.rabbits = make(map[uuid.UUID]*Rabbit)
	s.rabbit = NewRabbit(g)
	client.Write(s.rabbit.ToJson())
	return s
}

func (s *ClientScene) Update() error {

	if s.rabbit.Interact() {
		s.client.Write(s.rabbit.ToJson())
	}

	s.UpdateRabbits()

	for _, r := range s.rabbits {
		r.Update()
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

	for _, r := range g.rabbits {
		r.Draw(screen)
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

func (s *ClientScene) UpdateRabbits() {

	messages := s.client.ReadAll()

	for _, m := range messages {
		jsonData := []byte(m)
		var serial Serial
		if err := json.Unmarshal(jsonData, &serial); err != nil {
			log.Fatal(fmt.Errorf("Cannot unmarshal %s", m))
			continue
		}
		switch serial.ClassName {
		case "Rabbit":
			var rabbit Rabbit
			if err := json.Unmarshal(jsonData, &rabbit); err != nil {
				log.Fatal(fmt.Errorf("Cannot unmarshal the Rabbit %s", m))
				continue
			}
			var existing *Rabbit
			if s.rabbit.ID == rabbit.ID {
				existing = s.rabbit
			} else {
				existing = s.rabbits[rabbit.ID]
			}
			if existing != nil {
				existing.CopyFrom(&rabbit)
			} else {
				newRabbit := NewRabbit(s.game)
				newRabbit.CopyFrom(&rabbit)
				s.rabbits[rabbit.ID] = newRabbit
			}
		case "Lettuce":
			log.Fatal("lettuce not implemented")
		default:
			log.Printf("Cannot unmarshal the Rabbit %s", m)
		}
	}
}
