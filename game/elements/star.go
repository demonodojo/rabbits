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
		"Absolutno", "Acamar", "Achernar", "Achird", "Acrab", "Acrux", "Acubens",
		"Adhafera", "Adhara", "Adhil", "Ain", "Ainalrami", "Aiolos", "Aladfar",
		"Alasia", "Alathfar", "Albaldah", "Albali", "Albireo", "Alchiba", "Alcor",
		"Alcyone", "Aldebaran", "Alderamin", "Aldhanab", "Aldhibah", "Aldulfin", "Alfirk",
		"Algedi", "Algenib", "Algieba", "Algol", "Algorab", "Alhena", "Alioth", "Aljanah",
		"Alkaid", "Alkalurops", "Alkaphrah", "Alkarab", "Alkes", "Almaaz", "Almach",
		"Al Minliar al Asad", "Alnair", "Alnasl", "Alnilam", "Alnitak", "Alniyat", "Alphard",
		"Alphecca", "Alpheratz", "Alpherg", "Alrakis", "Alrescha", "Alruba", "Alsafi", "Alsciaukat",
		"Alsephina", "Alshain", "Alshat", "Altair", "Altais", "Alterf", "Aludra", "Alula Australis",
		"Alula Borealis", "Alya", "Alzirr", "Amadioha", "Amansinaya", "Anadolu", "Añañuca", "Ancha",
		"Angetenar", "Aniara", "Ankaa", "Anser", "Antares", "Arcalís", "Arcturus", "Arkab Posterior",
		"Arkab Prior", "Arneb", "Ascella", "Asellus Australis", "Asellus Borealis", "Ashlesha", "Asellus Primus",
		"Asellus Secundus", "Asellus Tertius", "Aspidiske", "Asterope", "Atakoraka", "Athebyne", "Atik", "Atlas",
		"Atria", "Avior", "Axólotl", "Ayeyarwady", "Azelfafage", "Azha", "Azmidi",
		"Baekdu", "Barnard's Star", "Baten Kaitos", "Batsũ̀", "Beemim", "Beid", "Belel", "Bélénos", "Bellatrix",
		"Berehynia", "Betelgeuse", "Bharani", "Bibhā", "Biham", "Bosona", "Botein", "Brachium", "Bubup", "Buna",
		"Bunda", "Canopus", "Capella", "Caph", "Castor", "Castula", "Cebalrai", "Ceibo", "Celaeno", "Cervantes",
		"Chalawan", "Chamukuy", "Chaophraya", "Chara", "Chasoň", "Chechia", "Chertan", "Citadelle", "Citalá",
		"Cocibolca", "Copernicus", "Cor Caroli", "Cujam", "Cursa", "Dabih", "Dalim", "Danfeng", "Deneb", "Deneb Algedi",
		"Denebola", "Diadem", "Dilmun", "Dingolay", "Diphda", "Dìwö", "Diya", "Dofida", "Dombay", "Dschubba", "Dubhe", "Dziban",
		"Ebla", "Edasich", "Electra", "Elgafar", "Elkurud", "Elnath", "Eltanin", "Emiw", "Enif", "Errai",
		"Fafnir", "Fang", "Fawaris", "Felis", "Felixvarela", "Filetdor", "Flegetonte", "Fomalhaut", "Formosa",
		"Franz", "Fulu", "Fumalsamakah", "Funi", "Furud", "Fuyue",
		"Gacrux", "Gakyid", "Gar", "Garnet Star", "Geminga", "Giausar", "Gienah", "Ginan", "Gloas", "Gnomon", "Gomeisa", "Graffias",
		"Guahayona", "Grumium", "Gudja", "Gumala", "Guniibuu",
		"Hadar", "Haedus", "Hamal", "Hassaleh", "Hatysa", "Helvetios", "Heze", "Hoggar", "Homam", "Horna", "Hunahpú", "Hunor",
		"Iklil", "Illyrian", "Imai", "Inquill", "Intan", "Intercrus", "Irena", "Itonda", "Izar",
		"Jabbah", "Jishui",
		"Kaffaljidhma", "Kaewkosin", "Kalausi", "Kamuy", "Kang", "Karaka", "Kaus Australis", "Kaus Borealis", "Kaus Media", "Kaveh",
		"Keid", "Khambalia", "Kitalpha", "Kochab", "Koeia", "Koit", "Komondor", "Kornephoros", "Kosjenka", "Kraz", "Kuma", "Kurhah",
		"La Superba", "Larawag", "Lerna", "Lesath", "Libertas", "Lich", "Liesma", "Lilii Borea", "Lionrock", "Lucilinburhuc", "Lusitânia",
		"Maasym", "Macondo", "Mago", "Mahasim", "Mahsati", "Maia", "Malmok", "Marfik", "Markab", "Markeb", "Márohu", "Marsic", "Matar",
		"Matza", "Maru", "Mazaalai", "Mebsuta", "Megrez", "Meissa", "Mekbuda", "Meleph", "Menkalinan", "Menkar", "Menkent", "Menkib",
		"Merak", "Merga", "Meridiana", "Merope", "Mesarthim", "Miaplacidus", "Mimosa", "Minchir", "Minelauva", "Mintaka", "Mira", "Mirach",
		"Miram", "Mirfak", "Mirzam", "Misam", "Mizar", "Moldoveanu", "Mönch", "Montuno", "Morava", "Moriah", "Mothallah", "Mouhoun",
		"Mpingo", "Muliphein", "Muphrid", "Muscida", "Musica", "Muspelheim",
		"Nahn", "Naledi", "Naos", "Nashira", "Násti", "Natasha", "Navi", "Nekkar", "Nembus", "Nenque", "Nervia", "Nihal", "Nikawiy",
		"Noquisi", "Nosaxa", "Nunki", "Nusakan", "Nushagak", "Nyamien",
		"Ogma", "Okab", "Orkaria",
		"Paikauhale", "Parumleo", "Peacock", "Petra", "Phact", "Phecda", "Pherkad", "Phoenicia", "Piautos", "Pincoya", "Pipirima",
		"Pipoltr", "Pleione", "Poerava", "Polaris", "Polaris Australis", "Polis", "Pollux", "Porrima", "Praecipua", "Prima Hyadum",
		"Procyon", "Propus", "Proxima Centauri",
		"Ran", "Rana", "Rapeto", "Rasalas", "Rasalgethi", "Rasalhague", "Rastaban", "Regor", "Regulus", "Revati", "Rigel",
		"Rigil Kentaurus", "Rosalíadecastro", "Rotanev", "Ruchbah", "Rukbat",
		"Sabik", "Saclateni", "Sadachbia", "Sadalbari", "Sadalmelik", "Sadalsuud", "Sadr", "Sagarmatha", "Saiph", "Salm", "Sāmaya",
		"Sansuna", "Sargas", "Sarin", "Sceptrum", "Scheat", "Schedar", "Secunda Hyadum", "Segin", "Seginus", "Sham", "Shama", "Sharjah",
		"Shaula", "Sheliak", "Sheratan", "Sika", "Sirius", "Situla", "Skat", "Sol", "Solaris", "Spica", "Sterrennacht", "Stribor", "Sualocin",
		"Subra", "Suhail", "Sulafat", "Syrma",
		"Tabit", "Taika", "Taiyangshou", "Taiyi", "Talitha", "Tangra", "Tania Australis", "Tania Borealis", "Tapecue", "Tarazed",
		"Tarf", "Taygeta", "Tegmine", "Tejat", "Terebellum", "Tevel", "Thabit", "Theemin", "Thuban", "Tiaki", "Tianguan", "Tianyi",
		"Timir", "Tislit", "Titawin", "Tojil", "Toliman", "Tonatiuh", "Torcular", "Tuiren", "Tupã", "Tupi", "Tureis",
		"Ukdah", "Uklun", "Unukalhai", "Unurgunite", "Uruk", "Uúba",
		"Vega", "Veritate", "Vindemiatrix",
		"Wasat", "Wattle", "Wazn", "Wezen", "Wouri", "Wurren",
		"Xamidimura", "Xihe", "Xuange",
		"Yed Posterior", "Yed Prior", "Yildun",
		"Zaniah", "Zaurak", "Zavijava", "Zembra", "Zhang", "Zibal", "Zosma", "Zubenelgenubi", "Zubenelhakrabi", "Zubeneschamali",
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
	vector.DrawFilledCircle(screen, float32(x), float32(y), float32(radius/4), c, true)
	vector.StrokeCircle(screen, float32(x), float32(y), 10, 5, selectedColor, true)
	vector.StrokeCircle(screen, float32(x), float32(y), (float32(l.size) * 8), 1, color.RGBA{0, 0, 255, 255}, true)
	//vector.DrawFilledCircle(screen, float32(x), float32(y), (float32(radius) * float32(l.size) / 2.0), color.RGBA{60, 60, 60, 1}, false)

	if l.selected {
		vector.StrokeCircle(screen, float32(x), float32(y), (float32(radius) * 30.0), 1, color.RGBA{128, 128, 128, 255}, true)
		text.Draw(screen, fmt.Sprintf("Name: %s", l.Name), assets.InfoFont, 0, 30, color.RGBA{70, 255, 255, 1})
		text.Draw(screen, fmt.Sprintf("Size: %d", l.size), assets.InfoFont, 0, 50, color.RGBA{70, 255, 255, 1})
	}

	if radius > 16 {
		text.Draw(screen, fmt.Sprintf(l.Name), assets.StarFont, int(x)-len(l.Name)*2, int(y+radius/2+30), color.RGBA{255, 255, 255, 255})
	}
}

func (l *Star) Select() {
	l.selected = true

}

func (l *Star) Toggle() {
	l.selected = !l.selected
}

func (l *Star) UnSelect() {
	l.selected = false
}

func (l *Star) Edit() {
	l.Action = "EDIT"
	l.channel.Enqueue(l.ToJson())
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
