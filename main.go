package main

import (
	"flag"
	"fmt"
	"github.com/demonodojo/rabbits/game"
	"github.com/demonodojo/rabbits/game/network"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := &game.Game{}
	var scene game.Scene
	serverMode := flag.Bool("server", false, "Inits the application in server mode")
	clientMode := flag.Bool("client", false, "Inits the application in client mode")
	url := "ws://localhost:8080/ws"
	// Parsea los flags desde los argumentos de línea de comandos
	flag.Parse()

	if *serverMode {
		fmt.Println("Iniciando en modo servidor...")
		server := network.Server{Port: ":8080"}
		server.Start()
		scene = game.NewServerScene(g, &server)
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
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
