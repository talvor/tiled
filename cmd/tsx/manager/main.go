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
	tm := manager.NewManager([]string{tilesetsDir, otherTilesetDir})

	tm.AddTilesetGroup("player", []string{
		"Player_Base_Running",
		"Shirt_Green_Running",
		"Pants_Blue_Running",
		"Shoes_Brown_Running",
		"Medium_Hair_Brown_Running",
		"Hands_Bare_Running",
	})

	log.Println("Total Tilesets:", len(tm.TilesetsByName))
	log.Println("Total TilesetGroups:", len(tm.TilesetGroups))

	tilesetGroup, err := tm.GetTilesetGroup("player")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Total Tilesets in group:", len(tilesetGroup))
}
