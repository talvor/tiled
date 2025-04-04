package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/talvor/tiled/tmx/manager"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	mapsDir := path.Join(homeDir, "Documents/examples/maps")

	tm := manager.NewManager([]string{mapsDir})
	js, err := json.Marshal(tm)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(js))
}
