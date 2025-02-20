package main

import (
	"concept/entities"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	Player       *entities.Player
	Enemies      []*entities.Enemy
	Potions      []*entities.Potion
	TilemapJSON  *TilemapJSON
	tilesets     []Tileset
	TilemapImage *ebiten.Image
	Camera       *Camera
}

func (g *Game) Update() error {
	// react to inputs (key presses)

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.Player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.Player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.Player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.Player.Y += 2
	}

	for _, enemy := range g.Enemies {
		if enemy.FollowsPlayer {
			if enemy.X < g.Player.X {
				enemy.X += 1
			} else if enemy.X > g.Player.X {
				enemy.X -= 1
			}

			if enemy.Y < g.Player.Y {
				enemy.Y += 1
			} else if enemy.Y > g.Player.Y {
				enemy.Y -= 1
			}
		}
	}

	for _, potion := range g.Potions {
		if g.Player.X >= potion.X {
			g.Player.Health += potion.HealingPoints
			log.Printf("Picked up a potion %d", g.Player.Health)
		}
	}

	g.Camera.FollowTarget(g.Player.X+8, g.Player.Y+8, 320, 240)
	g.Camera.Constraint(
		float64(g.TilemapJSON.Layers[0].Width)*16.0,
		float64(g.TilemapJSON.Layers[0].Height)*16.0,
		320,
		240,
	)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// void background
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}
	for i, layer := range g.TilemapJSON.Layers {
		for j, id := range layer.Data {
			if id == 0 {
				continue
			}

			x := j % layer.Width
			y := j / layer.Width

			x *= 16
			y *= 16

			img := g.tilesets[i].Img(id)

			opts.GeoM.Translate(float64(x), float64(y))
			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16.0))
			opts.GeoM.Translate(g.Camera.X, g.Camera.Y)

			screen.DrawImage(img, &opts)

			opts.GeoM.Reset()
		}
	}

	// Draw player
	opts.GeoM.Translate(g.Player.X, g.Player.Y)
	opts.GeoM.Translate(g.Camera.X, g.Camera.Y)
	screen.DrawImage(
		g.Player.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
		&opts,
	)
	opts.GeoM.Reset()

	for _, potion := range g.Potions {
		opts.GeoM.Translate(potion.X, potion.Y)
		opts.GeoM.Translate(g.Camera.X, g.Camera.Y)
		screen.DrawImage(
			potion.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}

	for _, enemie := range g.Enemies {
		opts.GeoM.Translate(enemie.X, enemie.Y)
		opts.GeoM.Translate(g.Camera.X, g.Camera.Y)
		screen.DrawImage(
			enemie.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}

	drawFPSTPS(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func NewGameSetting(
	playerImg *ebiten.Image,
	enemieImg *ebiten.Image,
	potionImg *ebiten.Image,
	tilemapJSON *TilemapJSON,
	tilesets []Tileset,
	tileImg *ebiten.Image,
	camera *Camera,
) *Game {
	return &Game{
		Player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   150.0,
				Y:   80.0,
			},
			Health: 50,
		},
		Enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: enemieImg,
					X:   50.0,
					Y:   50.0,
				},
				FollowsPlayer: false,
			},
			{
				Sprite: &entities.Sprite{
					Img: enemieImg,
					X:   75.0,
					Y:   75.0,
				},
				FollowsPlayer: false,
			},
			{
				Sprite: &entities.Sprite{
					Img: enemieImg,
					X:   150.0,
					Y:   150.0,
				},
				FollowsPlayer: false,
			},
		},
		Potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   300.0,
					Y:   80.0,
				},
				HealingPoints: 20,
			},
		},
		TilemapJSON:  tilemapJSON,
		tilesets:     tilesets,
		TilemapImage: tileImg,
		Camera:       camera,
	}
}

func drawFPSTPS(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.3f", ebiten.ActualFPS()), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.3f", ebiten.ActualTPS()), 0, 20)
}
