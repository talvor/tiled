package main

import (
	"fmt"
	"os"
	"path"

	"github.com/talvor/tiled/tsx"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	tilesetsDir := path.Join(homeDir, "Documents/examples/tilesets")

	// Create a new tileset manager and load all tilesets from the directory
	tsm := tsx.NewTilesetManager(tilesetsDir)
	tsm.AddTileset(path.Join(tilesetsDir, "NinjaDark.tsx"))

	tileset, _ := tsm.GetTilesetBySource(path.Join(tilesetsDir, "NinjaDark.tsx"))
	fmt.Println(tileset)
	tileset, _ = tsm.GetTilesetByName("NinjaDark")
	fmt.Println(tileset)
}
