package main

import (
	"log"
	"os"
	"path"

	"github.com/talvor/tiled/tsx/manager"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	tilesetsDir := path.Join(homeDir, "Documents/examples/tilesets")
	otherTilesetDir := path.Join(homeDir, "Documents/tilesets/Cute_Fantasy/")
	tm, err := manager.NewManager(tilesetsDir)
	if err != nil {
		log.Fatal(err)
	}

	tm.LoadTilesetsFromDir(otherTilesetDir)

	tm.DebugPrintTilesets()
}
