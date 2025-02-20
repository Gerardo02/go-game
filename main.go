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

	playerImg, _, err := ebitenutil.NewImageFromFile("./assets/sprites/player.png")
	if err != nil {
		log.Println("player error")
		log.Fatal(err)
	}

	enemieImg, _, err := ebitenutil.NewImageFromFile("./assets/sprites/enemie.png")
	if err != nil {
		log.Println("enemie error")
		log.Fatal(err)
	}

	potionImg, _, err := ebitenutil.NewImageFromFile("./assets/sprites/potion.png")
	if err != nil {
		log.Println("potion error")
		log.Fatal(err)
	}

	tileImg, _, err := ebitenutil.NewImageFromFile("./assets/sprites/tileset.png")
	if err != nil {
		log.Println("Tile img error")
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("./assets/maps/spawn.json")
	if err != nil {
		log.Println("Tile maps error")
		log.Fatal(err)
	}

	tilesets, err := tilemapJSON.GenerateTilesets()
	if err != nil {
		log.Println("Tile sets error")
		log.Fatal(err)
	}

	camera := NewCamera(0, 0)

	gameInit := NewGameSetting(
		playerImg,
		enemieImg,
		potionImg,
		tilemapJSON,
		tilesets,
		tileImg,
		camera,
	)

	if err := ebiten.RunGame(gameInit); err != nil {
		log.Fatal(err)
	}
}
