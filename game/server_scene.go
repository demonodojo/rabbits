package game

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/math/f64"

	"github.com/demonodojo/rabbits/assets"
	"github.com/demonodojo/rabbits/game/network"
)

type ServerScene struct {
	game              *Game
	camera            *Camera
	server            *network.Server
	lettuceSpawnTimer *Timer
	rabbits           map[uuid.UUID]*Rabbit
	lettuces          map[uuid.UUID]*Lettuce
	bullets           map[uuid.UUID]*Bullet
	lastUpdateTime    time.Time

	score         int
	scale         float64
	baseVelocity  float64
	velocityTimer *Timer
	mutex         sync.Mutex
}

func NewServerScene(g *Game, server *network.Server) *ServerScene {
	s := &ServerScene{
		game:              g,
		camera:            &Camera{ViewPort: f64.Vec2{screenWidth, screenHeight}},
		rabbits:           make(map[uuid.UUID]*Rabbit),
		lettuces:          make(map[uuid.UUID]*Lettuce),
		bullets:           make(map[uuid.UUID]*Bullet),
		lettuceSpawnTimer: NewTimer(lettuceSpawnTime),
		server:            server,
		baseVelocity:      baseMeteorVelocity,
		velocityTimer:     NewTimer(meteorSpeedUpTime),
		lastUpdateTime:    time.Now(),
	}

	return s
}

func (s *ServerScene) Update() error {

	s.CheckTime()
	s.UpdateRabbits()

	s.lettuceSpawnTimer.Update()
	if s.lettuceSpawnTimer.IsReady() {
		s.lettuceSpawnTimer.Reset()

		if len(s.lettuces) < 20 {
			l := NewLettuce()
			s.lettuces[l.ID] = l
			s.server.Broadcast(l.ToJson())
		}
	}

	for _, l := range s.lettuces {
		l.Update()
	}

	for _, rabbit := range s.rabbits {
		rabbit.Update()
	}

	for _, l := range s.lettuces {
		l.Update()
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

	for _, b := range s.bullets {
		b.Update()
		if b.Action == "DELETE" {
			delete(s.bullets, b.ID)
			s.server.Broadcast(b.ToJson())
		}
	}

	for _, r := range s.rabbits {
		for _, b := range s.bullets {
			if r.Collider().Intersects(b.Collider()) {
				r.Fired()
				b.Action = "DELETE"
				s.server.Broadcast(r.ToJson())
			}
		}
	}

	for _, r := range s.rabbits {
		for il, l := range s.lettuces {
			if l.Collider().Intersects(r.Collider()) {
				delete(s.lettuces, il)
				l.Action = "Delete"
				s.server.Broadcast(l.ToJson())
				r.Score++
				r.Action = "Score"
				s.server.Broadcast(r.ToJson())
				break
			}
		}
	}

	return nil
}

func (g *ServerScene) Draw(screen *ebiten.Image) {

	opts := &ebiten.DrawImageOptions{}
	// Ajusta la escala de la imagen. 1 es el tamaño original, valores mayores para hacer zoom in, menores para zoom out.
	opts.GeoM.Scale(g.scale, g.scale)
	// Dibuja la imagen en la pantalla con las opciones de escala.

	for _, rabbit := range g.rabbits {
		rabbit.Draw(screen, g.camera.Matrix)
	}

	for _, m := range g.lettuces {
		m.Draw(screen, g.camera.Matrix)
	}

	for _, b := range g.bullets {
		b.Draw(screen, g.camera.Matrix)
	}

	text.Draw(screen, fmt.Sprintf("%06d", g.score), assets.ScoreFont, screenWidth/2-100, 50, color.White)
	text.Draw(screen, fmt.Sprintf("%06d", len(g.bullets)), assets.InfoFont, 10, 50, color.White)
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
			log.Fatal(fmt.Errorf("cannot unmarshal %s", m.Message))
			continue
		}

		switch serial.ClassName {
		case "Rabbit":
			var rabbit Rabbit
			if err := json.Unmarshal(jsonData, &rabbit); err != nil {
				log.Fatal(fmt.Errorf("cannot unmarshal the Rabbit %s", m.Message))
				continue
			}

			if serial.Action == "FIRE" {
				position, rotation := rabbit.advancedPosition()
				b := NewBullet(position, rotation)
				s.bullets[b.ID] = b
				s.server.Broadcast(b.ToJson())
			} else {

				s.server.Broadcast(m.Message)
				existing := s.rabbits[rabbit.ID]
				if existing != nil {
					existing.CopyFrom(&rabbit)
				} else {
					newRabbit := NewRabbit(s.game)
					newRabbit.CopyFrom(&rabbit)
					s.rabbits[rabbit.ID] = newRabbit
				}
			}
		case "Lettuce":
			log.Fatal("lettuce not implemented")
		default:
			log.Printf("cannot unmarshal the Rabbit %s", m.Message)
		}
	}
}

func (s *ServerScene) CheckTime() {
	now := time.Now()
	delta := now.Sub(s.lastUpdateTime)
	s.lastUpdateTime = now

	// Convierte delta a milisegundos para una comparación fácil
	deltaMs := delta.Seconds() * 1000

	// Suponiendo que estás apuntando a 60 FPS, verifica si el delta de tiempo excede los 16.666 ms
	if deltaMs > (16.666 * 2) {
		fmt.Printf("Sobrepasado: %v ms\n", deltaMs-16.666)
	}

}
