package elements

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/demonodojo/rabbits/assets"
	"github.com/demonodojo/rabbits/game"
	"github.com/demonodojo/rabbits/game/network"
	"github.com/google/uuid"
)

type Star struct {
	game.Serial
	scale    float64
	Position game.Vector
	sprite   *ebiten.Image
	Name     string
	size     int
	selected bool
	color    color.Color
	channel  *network.MessageQueue
}

func NewStar(channel *network.MessageQueue, position game.Vector) *Star {
	names := []string{
		"Sol", "Sirius", "Canopus", "Rigil Kentaurus", "Arcturus",
		"Vega", "Capella", "Rigel", "Procyon", "Achernar",
		"Betelgeuse", "Altair", "Aldebaran", "Antares", "Spica",
		"Pollux", "Fomalhaut", "Deneb", "Regulus", "Castor",
		"Gacrux", "Bellatrix", "Elnath", "Miaplacidus", "Alnilam",
		"Alnair", "Alnitak", "Dubhe", "Mirfak", "Wezen",
		"Sargas", "Kaus Australis", "Avior", "Alkaid", "Menkent",
		"Acrux", "Alhena", "Peacock", "Mirzam", "Alphard",
		"Hadar", "Hamal", "Diphda", "Mimosa", "Regor",
		"Acamar", "Achernar", "Achird", "Acrab", "Acrux",
		"Acubens", "Adhafera", "Adhara", "Adhil", "Ain",
		"Ainalrami", "Aiolos", "Aladfar", "Alasia", "Albaldah",
		"Albali", "Albireo", "Alchiba", "Alcor", "Alcyone",
		"Barnard’s Star", "Baten Kaitos", "Beemim", "Beid", "Bellatrix",
		"Betelgeuse", "Bharani", "Biham", "Botein", "Brachium", "Bunda",
		"Canopus", "Capella", "Caph", "Castor", "Castula",
		"Cebalrai", "Celaeno", "Cervantes", "Chalawan", "Chamukuy",
		"Chara", "Chertan", "Copernicus", "Cor Caroli", "Cujam", "Cursa",
		"Dabih", "Dalim", "Deneb", "Deneb Algedi", "Denebola",
		"Diadem", "Diphda", "Dschubba", "Dubhe", "Dziban",
		"Edasich", "Electra", "Elgafar", "Elkurud", "Elnath",
		"Eltanin", "Enif", "Errai",
		"Fafnir", "Fang", "Fawaris", "Felis", "Fomalhaut",
		"Fulu", "Fumalsamakah", "Furud", "Fuyue",
		"Gacrux", "Giausar", "Gienah", "Ginan", "Gomeisa",
		"Grumium", "Gudja", "Guniibuu",
	}

	colors := []color.Color{
		color.RGBA{0, 0, 255, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 0, 255, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{0, 255, 255, 255},
		color.RGBA{0, 0, 0, 255},
		color.RGBA{125, 255, 255, 255},
	}
	sprite := assets.StarSprite
	scale := 1.0

	l := &Star{
		Serial: game.Serial{
			ID:        uuid.New(),
			ClassName: "Star",
			Action:    "Spawn",
		},
		Position: position,
		scale:    scale,
		sprite:   sprite,
		Name:     names[rand.Intn(100)],
		color:    colors[rand.Intn(8)],
		size:     rand.Intn(7) + 1,
		channel:  channel,
	}
	return l
}

func (s *Star) Update() {

}

func (l *Star) Draw(screen *ebiten.Image, geom ebiten.GeoM) {

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(l.scale, l.scale)
	//op.GeoM.Translate(-float64(bounds.Dx())*l.scale, -float64(bounds.Dy())*l.scale)
	op.GeoM.Translate(l.Position.X*l.scale, l.Position.Y*l.scale)
	op.GeoM.Concat(geom)

	//screen.DrawImage(l.sprite, op)
	x, y := op.GeoM.Apply(0, 0)
	x2, y2 := op.GeoM.Apply(16, 16)

	radius := game.EuclidianDistance(game.Vector{X: x, Y: y}, game.Vector{X: x2, Y: y2})
	c := color.RGBA{210, 210, 50, 1}
	selectedColor := l.color
	if l.selected {
		c = color.RGBA{70, 255, 255, 255}
		selectedColor = color.RGBA{200, 200, 200, 255}
	}
	vector.DrawFilledCircle(screen, float32(x), float32(y), float32(radius/4), c, false)
	vector.StrokeCircle(screen, float32(x), float32(y), 10, 5, selectedColor, true)
	vector.StrokeCircle(screen, float32(x), float32(y), (float32(l.size) * 8), 1, color.RGBA{0, 0, 255, 255}, false)
	//vector.DrawFilledCircle(screen, float32(x), float32(y), (float32(radius) * float32(l.size) / 2.0), color.RGBA{60, 60, 60, 1}, false)

	if l.selected {
		vector.StrokeCircle(screen, float32(x), float32(y), (float32(radius) * 30.0), 1, color.RGBA{128, 128, 128, 255}, false)
		text.Draw(screen, fmt.Sprintf("Name: %s", l.Name), assets.InfoFont, 0, 30, color.RGBA{70, 255, 255, 1})
		text.Draw(screen, fmt.Sprintf("Size: %d", l.size), assets.InfoFont, 0, 50, color.RGBA{70, 255, 255, 1})
	}
}

func (l *Star) Select() {
	l.selected = true
	l.Action = "EDIT"
	l.channel.Enqueue(l.ToJson())
}

func (l *Star) Toggle() {
	l.selected = !l.selected
	if l.selected {
		l.Action = "EDIT"
		l.channel.Enqueue(l.ToJson())
	}
}

func (l *Star) UnSelect() {
	l.selected = false
}

func (l *Star) Center() (float64, float64) {
	return l.Position.X, l.Position.Y
}

func (l *Star) Collider() game.Rect {

	return game.NewRect(
		l.Position.X-8,
		l.Position.Y-8,
		16,
		16,
	)
}

func (r *Star) ToJson() string {
	json, _ := json.Marshal(r)
	return string(json)
}

func (r *Star) CopyFrom(other *Star) {
	r.ID = other.ID
	r.Action = other.Action
	r.Position = other.Position
}
