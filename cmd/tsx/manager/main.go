package main

import (
	"os"
	"path"

	"github.com/talvor/tiled/tsx"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	tilesetsDir := path.Join(homeDir, "Documents/examples/tilesets")
	otherTilesetDir := path.Join(homeDir, "Documents/tilesets/Cute_Fantasy/")
	tm := tsx.NewTilesetManager(tilesetsDir)
	tm.LoadTilesetsFromDir(otherTilesetDir)

	tm.DebugPrintTilesets()
}
