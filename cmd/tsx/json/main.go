package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/talvor/tiled/tsx"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	tilesetPath := path.Join(homeDir, "Documents/examples/tilesets/NinjaDark.tsx")

	tileset, err := tsx.LoadFile(tilesetPath)
	if err != nil {
		panic(err)
	}

	tsJSON, err := json.Marshal(tileset)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(tsJSON))
}
