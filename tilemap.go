package main

import (
	"encoding/json"
	"os"
)

type TilemapLayerJSON struct {
	Data   []int  `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Name   string `json:"name"`
}

type tilesetMetaData struct {
	Gid  int    `json:"firstgid"`
	Path string `json:"source"`
}

type TilemapJSON struct {
	Layers   []TilemapLayerJSON `json:"layers"`
	TileSets []tilesetMetaData  `json:"tilesets"`
}

func (t *TilemapJSON) GenerateTilesets() ([]Tileset, error) {
	tilesets := make([]Tileset, 0)

	for _, tilesetData := range t.TileSets {
		tileset, err := NewTileSet("./assets/maps/"+tilesetData.Path, tilesetData.Gid)
		if err != nil {
			return nil, err
		}
		tilesets = append(tilesets, tileset)
	}

	return tilesets, nil
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	tilemapJSON := TilemapJSON{}

	err = json.Unmarshal(contents, &tilemapJSON)

	return &tilemapJSON, nil
}
