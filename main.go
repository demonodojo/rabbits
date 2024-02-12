package main

import (
	"flag"
	"fmt"
	"github.com/demonodojo/rabbits/game"
	"github.com/demonodojo/rabbits/game/network"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	serverMode := flag.Bool("server", false, "Inits the application in server mode")

	// Parsea los flags desde los argumentos de línea de comandos
	flag.Parse()

	if *serverMode {
		fmt.Println("Iniciando en modo servidor...")
		server := network.Server{Port: ":8080"}
		server.Start()
	} else {
		fmt.Println("Iniciando en modo cliente...")
	}

	g := game.NewGame()

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
