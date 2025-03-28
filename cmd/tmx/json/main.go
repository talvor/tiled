package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/talvor/tiled/tmx"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	tmxPath := path.Join(homeDir, "Documents/examples/maps/StartScene.tmx")

	m, err := tmx.LoadFile(tmxPath)
	if err != nil {
		panic(err)
	}

	mJSON, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(mJSON))
}
