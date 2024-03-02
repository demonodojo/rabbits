//go:build js && wasm

package main

import (
	"fmt"
	"github.com/demonodojo/rabbits/game"
	"github.com/demonodojo/rabbits/game/network"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := &game.Game{}
	var scene game.Scene
	url := "ws://localhost:8080/ws"
	fmt.Println("Iniciando...")
	if true {
		fmt.Println("Iniciando en javascript...")
		client, err := network.NewJSClient(url)
		if err != nil {
			fmt.Println("Error")
			//log.Fatal("dial:", err)
			return
		} else {
			fmt.Println("Creando escena...")
			scene = game.NewClientScene(g, client)
		}
		defer client.Close()
	} else {
		fmt.Println("Iniciando cliente...")
		_, err := network.NewJSClient(url)
		if err != nil {
			fmt.Println("Error")
		}
		scene = game.NewRabbitDirectScene(g)
	}

	fmt.Println("Iniciando juego...")
	g.SetScene(scene)
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
