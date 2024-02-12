package main

import (
	"github.com/demonodojo/rabbits/game"
	"github.com/demonodojo/rabbits/game/network"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	if true {
		server := network.Server{Port: ":8080"}
		go server.Start() // Inicia el servidor en un goroutine
	}

	g := game.NewGame()

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
