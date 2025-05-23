package main

import (
	"os"
	"path"

	"github.com/talvor/tiled/animation/manager"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	animationsDir := path.Join(homeDir, "Documents/examples/animations")
	am := manager.NewManager([]string{animationsDir})
	am.DebugPrintAnimations()
}
