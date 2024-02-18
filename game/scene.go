package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (int, int)
	SpawnElement(name string, g interface{})
}

type SceneManager struct {
	currentScene Scene
}

func (s *SceneManager) Update() error {
	return s.currentScene.Update()
}

func (s *SceneManager) Draw(screen *ebiten.Image) {
	s.currentScene.Draw(screen)
}
