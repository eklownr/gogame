package tilemaps

import (
	"encoding/json"
	"os"
)

type TylemapLayers struct {
	Data   []int `json:"data"`
	Whith  int   `json:"width"`
	Height int   `json:"height"`
}
type TilemapJSON struct {
	Tiles []TylemapLayers `json:"layers"`
}

func newTilemapJSON(filepath string) (*TilemapJSON, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var tilemapJSON TilemapJSON
	err = json.Unmarshal(content, &tilemapJSON)
	if err != nil {
		return nil, err
	}
	return &tilemapJSON, nil
}
