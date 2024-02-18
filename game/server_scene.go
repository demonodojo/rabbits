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

type ServerScene struct {
	game     *Game
	server   *network.Server
	rabbits  map[uuid.UUID]*Rabbit
	lettuces []*Lettuce

	score         int
	scale         float64
	baseVelocity  float64
	velocityTimer *Timer
}

func NewServerScene(g *Game, server *network.Server) *ServerScene {
	s := &ServerScene{
		game:          g,
		rabbits:       make(map[uuid.UUID]*Rabbit),
		server:        server,
		baseVelocity:  baseMeteorVelocity,
		velocityTimer: NewTimer(meteorSpeedUpTime),
	}

	m := NewLettuce(s.baseVelocity)
	s.lettuces = append(s.lettuces, m)
	return s
}

func (g *ServerScene) Update() error {

	g.UpdateRabbits()

	for _, rabbit := range g.rabbits {
		rabbit.Update()
	}

	for _, l := range g.lettuces {
		l.Update()
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

	return nil
}

func (g *ServerScene) Draw(screen *ebiten.Image) {

	opts := &ebiten.DrawImageOptions{}
	// Ajusta la escala de la imagen. 1 es el tamaño original, valores mayores para hacer zoom in, menores para zoom out.
	opts.GeoM.Scale(g.scale, g.scale)
	// Dibuja la imagen en la pantalla con las opciones de escala.

	for _, rabbit := range g.rabbits {
		rabbit.Draw(screen)
	}

	for _, m := range g.lettuces {
		m.Draw(screen)
	}

	// for _, b := range g.bullets {
	// 	b.Draw(screen)
	// }

	text.Draw(screen, fmt.Sprintf("%06d", g.score), assets.ScoreFont, screenWidth/2-100, 50, color.White)
}

func (g *ServerScene) Reset() {
	g.rabbits = make(map[uuid.UUID]*Rabbit)
	g.lettuces = nil
	g.score = 0
	g.baseVelocity = baseMeteorVelocity
	g.velocityTimer.Reset()
}

func (g *ServerScene) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (s *ServerScene) SpawnElement(name string, element interface{}) {

}

func (s *ServerScene) UpdateRabbits() {

	messages := s.server.ReadAll()

	for _, m := range messages {
		jsonData := []byte(m.Message)
		var serial Serial
		if err := json.Unmarshal(jsonData, &serial); err != nil {
			log.Fatal(fmt.Errorf("Cannot unmarshal %s", m.Message))
			continue
		}
		switch serial.ClassName {
		case "Rabbit":
			var rabbit Rabbit
			if err := json.Unmarshal(jsonData, &rabbit); err != nil {
				log.Fatal(fmt.Errorf("Cannot unmarshal the Rabbit %s", m.Message))
				continue
			}
			existing := s.rabbits[rabbit.ID]
			if existing != nil {
				existing.CopyFrom(&rabbit)
			} else {
				newRabbit := NewRabbit(s.game)
				newRabbit.CopyFrom(&rabbit)
				s.rabbits[rabbit.ID] = newRabbit
			}
			rabbit.Update()
		case "Lettuce":
			log.Fatal("lettuce not implemented")
		default:
			log.Printf("Cannot unmarshal the Rabbit %s", m.Message)
		}
	}
}
