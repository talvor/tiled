package main

import (
	"fmt"
	"os"
	"path"

	"github.com/talvor/tiled/tmx"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	tmxPath := path.Join(homeDir, "Documents/tilesets/Cute_Fantasy/maps/Home.tmx")

	m, err := tmx.LoadFile(tmxPath)
	if err != nil {
		panic(err)
	}

	l, err := m.GetLayer("collider")
	if err != nil {
		panic(err)
	}

	fmt.Println(l.GetTileRectFromIndex(1, m))
	// mJSON, err := json.Marshal(m)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(mJSON))
}
