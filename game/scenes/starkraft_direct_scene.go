package scenes

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
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

type StarKraftDirectScene struct {
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

func NewStarKraftDirectScene(g *game.Game) *StarKraftDirectScene {
	s := &StarKraftDirectScene{
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

func (g *StarKraftDirectScene) Update() error {

	if g.starForm != nil {
		g.starForm.Update()
		if g.starForm.Action == "SUBMIT" {
			g.starForm = nil
		}
	} else {
		g.rabbit.Interact()

		selected := g.Interact()
		if !selected {
			g.camera.Update(g.rabbit)
		}

	}

	g.Spawn()

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

func (s *StarKraftDirectScene) Interact() bool {
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
		return selected != nil
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		var selected *elements.Star
		lastRadius := 5000.0
		x, y := ebiten.CursorPosition()

		wx, wy := s.camera.ScreenToWorld(x, y)
		for _, l := range s.stars {
			radius := game.EuclidianDistance(game.Vector{X: l.Position.X, Y: l.Position.Y}, game.Vector{X: float64(wx), Y: float64(wy)})
			if radius < 20 && radius < lastRadius {
				selected = l
				lastRadius = radius
			}
		}
		if selected != nil {
			selected.Edit()
		}
		return selected != nil
	}
	return false
}

func (g *StarKraftDirectScene) Draw(screen *ebiten.Image) {

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

	if g.starForm != nil {
		g.starForm.Draw(screen, g.camera.Matrix)
		return
	}
}

func (g *StarKraftDirectScene) Reset() {
	g.rabbit = game.NewRabbit(g.game)
	g.stars = nil
	g.bullets = nil
	g.score = 0
	g.baseVelocity = baseMeteorVelocity
	g.velocityTimer.Reset()
}

func (g *StarKraftDirectScene) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (s *StarKraftDirectScene) SpawnElement(name string, element interface{}) {

}

func (s *StarKraftDirectScene) Spawn() {

	messages := s.spawns.ReadAll()

	for _, m := range messages {
		jsonData := []byte(m)
		var serial game.Serial
		if err := json.Unmarshal(jsonData, &serial); err != nil {
			log.Fatal(fmt.Errorf("Cannot unmarshal %s", m))
			continue
		}
		switch serial.ClassName {
		case "Star":
			if serial.Action == "EDIT" {
				existing := s.starById(serial.ID)
				newForm := forms.NewStarForm(existing)
				s.starForm = newForm
			}

		default:
			log.Printf("Cannot unmarshal the Elemetnt %s", m)
		}
	}
}

func (s *StarKraftDirectScene) starById(ID uuid.UUID) *elements.Star {
	for _, star := range s.stars {
		if star.ID == ID {
			return star
		}
	}
	return nil
}
