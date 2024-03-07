//go:build !js || !wasm

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/demonodojo/rabbits/game"
	"github.com/demonodojo/rabbits/game/network"
	"github.com/demonodojo/rabbits/game/scenes"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := &game.Game{}
	var scene game.Scene
	serverMode := flag.Bool("server", false, "Inits the application in server mode")
	clientMode := flag.Bool("client", false, "Inits the application in client mode")
	directMode := flag.Bool("direct", false, "Inits the application in direct mode")
	starsMode := flag.Bool("stars", false, "Inits the application in stars mode")
	url := "ws://localhost:8080/ws"
	// Parsea los flags desde los argumentos de línea de comandos
	flag.Parse()

	if *serverMode {
		fmt.Println("Iniciando en modo servidor...")
		server := network.Server{Port: ":8080"}
		server.Start()
		scene = game.NewServerScene(g, &server)
	} else if *directMode {
		scene = game.NewRabbitDirectScene(g)
	} else if *starsMode {
		scene = scenes.NewStarsDirectScene(g)
	} else if *clientMode {
		fmt.Println("Iniciando en modo cliente...")
		client, err := network.NewClient(url)
		if err != nil {
			log.Fatal("dial:", err)
		} else {
			scene = game.NewClientScene(g, client)
		}
		defer client.Close()
	} else {
		scene = game.NewRabbitDirectScene(g)
	}

	fmt.Println("Iniciando juego...")
	g.SetScene(scene)

	ebiten.SetWindowTitle("Rabbits")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
