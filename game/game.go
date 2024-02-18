package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	currentScene Scene
}

// Implementa ebiten.Game para Game
func (g *Game) Update() error {
	return g.currentScene.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.currentScene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.currentScene.Layout(outsideWidth, outsideHeight)
}

// SetScene cambia la escena actual del juego
func (g *Game) SetScene(scene Scene) {
	g.currentScene = scene
}

// GetCurrentScene devuelve la escena actual del juego
func (g *Game) GetCurrentScene() Scene {
	return g.currentScene
}
