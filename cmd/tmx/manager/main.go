package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/talvor/tiled/tmx/manager"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	mapsDir := path.Join(homeDir, "Documents/examples/maps")

	tm, err := manager.NewManager(mapsDir)
	if err != nil {
		log.Fatal(err)
	}

	js, err := json.Marshal(tm)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(js))
}
