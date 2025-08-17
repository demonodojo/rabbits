//go:build js && wasm

package main

import (
	"fmt"
	"github.com/demonodojo/rabbits/game"
	"github.com/demonodojo/rabbits/game/scenes"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := &game.Game{}
	var scene game.Scene
	fmt.Println("Iniciando...")
	
	// Por defecto, mostrar el menú principal
	scene = scenes.NewMenuScene(g)

	fmt.Println("Iniciando juego...")
	g.SetScene(scene)
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
