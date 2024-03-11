package assets

import (
	"embed"
	"image"
	_ "image/png"
	"io/fs"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

//go:embed *
var assets embed.FS

var RabbitSprite = mustLoadImage("tile_0106.png")
var RabbitSpriteR = mustLoadImage("tile_0160.png")
var PlayerSprite = mustLoadImage("player.png")
var CarrierSprite = mustLoadImage("player.png")
var LettuceSprite = mustLoadImage("tile_0094.png")
var StarSprite = mustLoadImage("star0.png")
var MeteorSprites = mustLoadImages("meteors/*.png")
var LaserSprite = mustLoadImage("laser.png")
var ScoreFont = mustLoadFont("font.ttf", 48)
var InfoFont = mustLoadFont("font.ttf", 14)
var StarFont = regularFont(12)

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func mustLoadImages(path string) []*ebiten.Image {
	matches, err := fs.Glob(assets, path)
	if err != nil {
		panic(err)
	}

	images := make([]*ebiten.Image, len(matches))
	for i, match := range matches {
		images[i] = mustLoadImage(match)
	}

	return images
}

func mustLoadFont(name string, size float64) font.Face {
	f, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}

	tt, err := opentype.Parse(f)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err)
	}

	return face
}

func regularFont(size float64) font.Face {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}
