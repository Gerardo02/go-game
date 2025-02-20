package main

import (
	"encoding/json"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tileset interface {
	Img(id int) *ebiten.Image
}

type TileJSON struct {
	ID     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imageWidth"`
	Height int    `json:"imageHeight"`
}

type DynTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

type UniformTilesetJSON struct {
	Path string `json:"image"`
}

type DynTileset struct {
	imgs []*ebiten.Image
	gid  int
}

type UniformTileset struct {
	img *ebiten.Image
	gid int
}

func (u *UniformTileset) Img(id int) *ebiten.Image {
	// tileID - tilesetOffset = id relative to tileset
	id -= u.gid

	srcX := id % 28
	srcY := id / 28

	srcX *= 16
	srcY *= 16

	return u.img.SubImage(
		image.Rect(srcX, srcY, srcX+16, srcY+16),
	).(*ebiten.Image)
}

func (d *DynTileset) Img(id int) *ebiten.Image {
	// tileID - tilesetOffset = id relative to tileset
	id -= d.gid

	return d.imgs[id]
}

func NewTileSet(path string, gid int) (Tileset, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if strings.Contains(path, "buildings") {
		var dynTilesetJSON DynTilesetJSON
		err = json.Unmarshal(contents, &dynTilesetJSON)
		if err != nil {
			return nil, err
		}

		dynTileset := DynTileset{
			gid:  gid,
			imgs: make([]*ebiten.Image, 0),
		}

		for _, tileJSON := range dynTilesetJSON.Tiles {
			tileJSONPath := tileJSON.Path
			tileJSONPath = filepath.Clean(tileJSONPath)
			tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = filepath.Join("./assets/", tileJSONPath)

			img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
			if err != nil {
				return nil, err
			}

			dynTileset.imgs = append(dynTileset.imgs, img)
		}
		return &dynTileset, nil

	}

	var uniformTilesetJSON UniformTilesetJSON
	err = json.Unmarshal(contents, &uniformTilesetJSON)
	if err != nil {
		return nil, err
	}

	tileJSONPath := uniformTilesetJSON.Path
	tileJSONPath = filepath.Clean(tileJSONPath)
	tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
	tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
	tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
	tileJSONPath = filepath.Join("./assets/", tileJSONPath)

	img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
	if err != nil {
		return nil, err
	}

	uniformTileset := UniformTileset{
		img: img,
		gid: gid,
	}

	return &uniformTileset, nil
}
