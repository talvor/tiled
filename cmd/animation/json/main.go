package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/talvor/tiled/animation"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	animationPath := path.Join(homeDir, "Documents/tilesets/Cute_Fantasy/animations/animations.ani")

	a, err := animation.LoadFile(animationPath)
	if err != nil {
		panic(err)
	}

	animJSON, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(animJSON))
}
