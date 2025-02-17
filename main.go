package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Img *ebiten.Image
	X   float64
	Y   float64
}

type Enemy struct {
	*Sprite
	FollowsPlayer bool
}

type Potion struct {
	*Sprite
	HealingPoints uint
}

type Player struct {
	*Sprite
	Health uint
}

type Game struct {
	Player  *Player
	Enemies []*Enemy
	Potions []*Potion
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
		if g.Player.X > potion.X {
			g.Player.Health += potion.HealingPoints
			log.Printf("Picked up a potion %d", g.Player.Health)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}

	// Draw player
	opts.GeoM.Translate(g.Player.X, g.Player.Y)
	screen.DrawImage(
		g.Player.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
		&opts,
	)
	opts.GeoM.Reset()

	for _, potion := range g.Potions {
		opts.GeoM.Translate(potion.X, potion.Y)
		screen.DrawImage(
			potion.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}

	for _, enemie := range g.Enemies {
		opts.GeoM.Translate(enemie.X, enemie.Y)
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

	gameInit := Game{
		Player: &Player{
			Sprite: &Sprite{
				Img: playerImg,
				X:   100.0,
				Y:   100.0,
			},
			Health: 50,
		},
		Enemies: []*Enemy{
			{
				&Sprite{
					Img: enemieImg,
					X:   50.0,
					Y:   50.0,
				},
				false,
			},
			{
				&Sprite{
					Img: enemieImg,
					X:   70.0,
					Y:   70.0,
				},
				true,
			},
			{
				&Sprite{
					Img: enemieImg,
					X:   150.0,
					Y:   150.0,
				},
				true,
			},
		},
		Potions: []*Potion{
			{
				&Sprite{
					Img: potionImg,
					X:   125.0,
					Y:   125.0,
				},
				20,
			},
			{
				&Sprite{
					Img: potionImg,
					X:   225.0,
					Y:   200.0,
				},
				15,
			},
		},
	}

	if err := ebiten.RunGame(&gameInit); err != nil {
		log.Fatal(err)
	}
}
