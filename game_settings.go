package main

import (
	"concept/entities"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Player       *entities.Player
	Enemies      []*entities.Enemy
	Potions      []*entities.Potion
	TilemapJSON  *TilemapJSON
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

	for _, layer := range g.TilemapJSON.Layers {
		for index, id := range layer.Data {
			x := index % layer.Width
			y := index / layer.Width

			x *= 16
			y *= 16

			srcX := (id - 1) % 28
			srcY := (id - 1) / 28

			srcX *= 16
			srcY *= 16

			opts.GeoM.Translate(float64(x), float64(y))
			opts.GeoM.Translate(g.Camera.X, g.Camera.Y)
			screen.DrawImage(
				g.TilemapImage.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				&opts,
			)
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func NewGameSetting(
	playerImg *ebiten.Image,
	enemieImg *ebiten.Image,
	potionImg *ebiten.Image,
	tilemapJSON *TilemapJSON,
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
		TilemapImage: tileImg,
		Camera:       camera,
	}
}
