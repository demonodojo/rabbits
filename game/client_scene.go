package game

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/math/f64"

	"github.com/demonodojo/rabbits/assets"
	"github.com/demonodojo/rabbits/game/network"
)

type ClientScene struct {
	game          *Game
	camera        *Camera
	client        network.GenericClient
	rabbit        *Rabbit
	rabbits       map[uuid.UUID]*Rabbit
	lettuces      map[uuid.UUID]*Lettuce
	bullets       map[uuid.UUID]*Bullet
	rabbitsOrder  []uuid.UUID
	lettucesOrder []uuid.UUID
	bulletsOrder  []uuid.UUID

	score         int
	scale         float64
	baseVelocity  float64
	velocityTimer *Timer
}

func NewClientScene(g *Game, client network.GenericClient) *ClientScene {
	s := &ClientScene{
		game:          g,
		camera:        &Camera{ViewPort: f64.Vec2{screenWidth, screenHeight}},
		client:        client,
		baseVelocity:  baseMeteorVelocity,
		velocityTimer: NewTimer(meteorSpeedUpTime),
	}

	s.camera.Reset()
	s.rabbits = make(map[uuid.UUID]*Rabbit)
	s.lettuces = make(map[uuid.UUID]*Lettuce)
	s.bullets = make(map[uuid.UUID]*Bullet)
	s.rabbit = NewRabbit(g)
	s.rabbits[s.rabbit.ID] = s.rabbit

	s.UpdateRabbitsOrder()
	s.UpdateLettucesOrder()
	s.UpdateBulletsOrder()

	client.Write(s.rabbit.ToJson())
	return s
}

func (s *ClientScene) Update() error {

	if s.rabbit.Interact() {
		s.client.Write(s.rabbit.ToJson())
		s.rabbit.Action = "NONE"
	}

	s.UpdateRabbits()

	for _, r := range s.rabbits {
		r.Update()
	}

	s.camera.Update(s.rabbit)

	for _, l := range s.lettuces {
		l.Update()
		if l.Action == "Delete" {
			delete(s.lettuces, l.ID)
			s.UpdateLettucesOrder()
		}
	}

	for _, b := range s.bullets {
		b.Update()
		if b.Action == "DELETE" {
			delete(s.bullets, b.ID)
			s.UpdateBulletsOrder()
		}
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

	for _, id := range g.rabbitsOrder {
		r := g.rabbits[id]
		r.Draw(screen, g.camera.Matrix)
	}

	for _, id := range g.lettucesOrder {
		l := g.lettuces[id]
		l.Draw(screen, g.camera.Matrix)
	}

	for _, id := range g.bulletsOrder {
		l := g.bullets[id]
		l.Draw(screen, g.camera.Matrix)
	}

	text.Draw(screen, fmt.Sprintf("%06d", g.rabbit.Score), assets.ScoreFont, screenWidth/2-100, 50, color.White)
	text.Draw(screen, fmt.Sprintf("%06d", len(g.bullets)), assets.InfoFont, 10, 50, color.White)
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
				s.rabbit.Score = rabbit.Score
				s.rabbit.Speed = rabbit.Speed
				if EuclidianDistance(rabbit.Position, s.rabbit.Position) > 100.0 {
					s.rabbit.Position = rabbit.Position
				}
				continue
			} else {
				existing = s.rabbits[rabbit.ID]
			}
			if existing != nil {
				existing.CopyFrom(&rabbit)
			} else {
				newRabbit := NewRabbit(s.game)
				newRabbit.CopyFrom(&rabbit)
				s.rabbits[rabbit.ID] = newRabbit
				s.UpdateRabbitsOrder()
			}

		case "Lettuce":
			var lettuce Lettuce
			if err := json.Unmarshal(jsonData, &lettuce); err != nil {
				log.Fatal(fmt.Errorf("Cannot unmarshal the Lettuce %s", m))
				continue
			}
			existing := s.lettuces[lettuce.ID]
			if existing != nil {
				existing.CopyFrom(&lettuce)
			} else {
				l := NewLettuce()
				l.CopyFrom(&lettuce)
				s.lettuces[lettuce.ID] = l
				s.UpdateLettucesOrder()
			}

		case "Bullet":
			var bullet Bullet
			if err := json.Unmarshal(jsonData, &bullet); err != nil {
				log.Fatal(fmt.Errorf("Cannot unmarshal the Bullet %s", m))
				continue
			}
			existing := s.bullets[bullet.ID]
			if existing != nil {
				existing.CopyFrom(&bullet)
			} else {
				b := NewBullet(bullet.Position, bullet.Rotation)
				b.CopyFrom(&bullet)
				s.bullets[bullet.ID] = b
				s.UpdateBulletsOrder()
			}

		default:
			log.Printf("Cannot unmarshal the Message %s", m)
		}
	}
}

func (s *ClientScene) UpdateRabbitsOrder() {
	s.rabbitsOrder = GetOrderedIds(s.rabbits)
}

func (s *ClientScene) UpdateLettucesOrder() {
	s.lettucesOrder = GetOrderedIds(s.lettuces)
}

func (s *ClientScene) UpdateBulletsOrder() {
	s.bulletsOrder = GetOrderedIds(s.bullets)
}
