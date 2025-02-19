package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Jueguito loco")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImg, _, err := ebitenutil.NewImageFromFile("./assets/images/player.png")
	if err != nil {
		log.Fatal(err)
	}

	enemieImg, _, err := ebitenutil.NewImageFromFile("./assets/images/enemie.png")
	if err != nil {
		log.Fatal(err)
	}

	potionImg, _, err := ebitenutil.NewImageFromFile("./assets/images/potion.png")
	if err != nil {
		log.Fatal(err)
	}

	tileImg, _, err := ebitenutil.NewImageFromFile("./assets/images/tileset.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("./assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}
	camera := NewCamera(0, 0)

	gameInit := NewGameSetting(
		playerImg,
		enemieImg,
		potionImg,
		tilemapJSON,
		tileImg,
		camera,
	)

	if err := ebiten.RunGame(gameInit); err != nil {
		log.Fatal(err)
	}
}
