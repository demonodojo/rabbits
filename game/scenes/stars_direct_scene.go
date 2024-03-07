package scenes

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/math/f64"

	"github.com/demonodojo/rabbits/assets"
	"github.com/demonodojo/rabbits/game"
	"github.com/demonodojo/rabbits/game/elements"
	"github.com/demonodojo/rabbits/game/elements/forms"
	"github.com/demonodojo/rabbits/game/network"
)

const (
	screenWidth  = 800
	screenHeight = 600

	lettuceSpawnTime = 1 * time.Second

	baseMeteorVelocity  = 0.25
	meteorSpeedUpAmount = 0.1
	meteorSpeedUpTime   = 5 * time.Second
)

type StarsDirectScene struct {
	game      *game.Game
	camera    *game.Camera
	offscreen *ebiten.Image
	player    *game.Player
	rabbit    *game.Rabbit
	stars     []*elements.Star
	starForm  *forms.StarForm
	bullets   []*game.Bullet
	spawns    *network.MessageQueue

	score         int
	scale         float64
	baseVelocity  float64
	velocityTimer *game.Timer
}

func NewStarsDirectScene(g *game.Game) *StarsDirectScene {
	s := &StarsDirectScene{
		game:          g,
		camera:        &game.Camera{ViewPort: f64.Vec2{screenWidth, screenHeight}},
		baseVelocity:  baseMeteorVelocity,
		velocityTimer: game.NewTimer(meteorSpeedUpTime),
		spawns:        network.NewMessageQueue(),
	}

	s.camera.Reset()
	s.rabbit = game.NewRabbit(g)
	rand.Seed(time.Now().UnixNano())
	stars := map[string]*elements.Star{}
	for i := 0; i < 50; i++ {
		pos := game.Vector{
			X: (float64(screenWidth*4-16) * rand.Float64()) - screenWidth*1.5,
			Y: (float64(screenHeight*4-16) * rand.Float64()) - screenHeight*1.5,
		}
		m := elements.NewStar(s.spawns, pos)
		stars[m.Name] = m
		s.stars = append(s.stars, m)
	}

	for _, v := range stars {
		s.stars = append(s.stars, v)
	}

	return s
}

func (g *StarsDirectScene) Update() error {

	g.Spawn()

	g.rabbit.Interact()

	g.Interact()

	g.scale += 0.01
	if g.scale > 2 {
		g.scale = 0.5
	}

	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += meteorSpeedUpAmount
	}

	for _, l := range g.stars {
		l.Update()
	}

	if ebiten.IsKeyPressed(ebiten.Key1) {
		ebiten.SetFullscreen(true)
	}

	if ebiten.IsKeyPressed(ebiten.Key2) {
		ebiten.SetFullscreen(false)
	}

	g.camera.Update(g.rabbit)

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

	// Check for rabbit/stars collision

	for i, l := range g.stars {
		if l.Collider().Intersects(g.rabbit.Collider()) {
			g.stars = append(g.stars[:i], g.stars[i+1:]...)
			g.score++
			if len(g.stars) == 0 {
				g.Reset()
			}
			//g.Reset()
			break
		}
	}

	return nil
}

func (s *StarsDirectScene) Interact() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		var selected *elements.Star
		lastRadius := 5000.0
		x, y := ebiten.CursorPosition()

		wx, wy := s.camera.ScreenToWorld(x, y)
		for _, l := range s.stars {
			radius := game.EuclidianDistance(game.Vector{X: l.Position.X, Y: l.Position.Y}, game.Vector{X: float64(wx), Y: float64(wy)})
			if radius < 20 && radius < lastRadius {
				if selected != nil {
					selected.UnSelect()
				}
				selected = l
				lastRadius = radius
			} else {
				l.UnSelect()
			}
		}
		if selected != nil {
			selected.Toggle()
		}
	}
}

func (g *StarsDirectScene) Draw(screen *ebiten.Image) {

	opts := &ebiten.DrawImageOptions{}
	// Ajusta la escala de la imagen. 1 es el tamaño original, valores mayores para hacer zoom in, menores para zoom out.
	opts.GeoM.Scale(g.scale, g.scale)
	// Dibuja la imagen en la pantalla con las opciones de escala.

	g.rabbit.Draw(screen, g.camera.Matrix)

	for _, l := range g.stars {
		l.Draw(screen, g.camera.Matrix)
	}

	// for _, b := range g.bullets {
	// 	b.Draw(screen)
	// }

	text.Draw(screen, fmt.Sprintf("%06d", g.score), assets.ScoreFont, screenWidth/2-100, 50, color.White)
}

func (g *StarsDirectScene) Reset() {
	g.rabbit = game.NewRabbit(g.game)
	g.stars = nil
	g.bullets = nil
	g.score = 0
	g.baseVelocity = baseMeteorVelocity
	g.velocityTimer.Reset()
}

func (g *StarsDirectScene) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (s *StarsDirectScene) SpawnElement(name string, element interface{}) {

}

func (s *StarsDirectScene) Spawn() {

	messages := s.spawns.ReadAll()

	for _, m := range messages {
		jsonData := []byte(m)
		var serial game.Serial
		if err := json.Unmarshal(jsonData, &serial); err != nil {
			log.Fatal(fmt.Errorf("Cannot unmarshal %s", m))
			continue
		}
		switch serial.ClassName {
		case "Form":
			var form game.Form
			if err := json.Unmarshal(jsonData, &form); err != nil {
				log.Fatal(fmt.Errorf("Cannot unmarshal the Form %s", m))
				continue
			}

			newForm := game.NewForm()
			newForm.CopyFrom(&form)
		default:
			log.Printf("Cannot unmarshal the Rabbit %s", m)
		}
	}
}
