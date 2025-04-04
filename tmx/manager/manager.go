package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/talvor/tiled/common"
	"github.com/talvor/tiled/tmx"
)

var ErrMapNotFound = errors.New("tsx: map not found")

type MapManager struct {
	Maps map[string]*tmx.Map
}

func (mm *MapManager) GetMapByName(name string) (*tmx.Map, error) {
	if m, ok := mm.Maps[name]; ok {
		return m, nil
	}

	return nil, ErrMapNotFound
}

func (mm *MapManager) LoadMapsFromDir(dir string) {
	loadMaps(mm, dir)
}

func (mm *MapManager) DebugPrintMaps() {
	for name := range mm.Maps {
		fmt.Println(name)
	}
}

func NewManager(baseDirs []string) *MapManager {
	mm := &MapManager{
		Maps: make(map[string]*tmx.Map),
	}

	// Load maps from the base directories
	for _, baseDir := range baseDirs {
		if err := common.PathShouldBeDirectory(baseDir); err != nil {
			fmt.Printf("Error: %s %v\n", baseDir, err)
			continue
		}

		if err := loadMaps(mm, baseDir); err != nil {
			fmt.Printf("Error loading maps from %s: %v\n", baseDir, err)
		}
	}

	return mm
}

func loadMaps(mm *MapManager, baseDir string) error {
	tmxFiles, err := findTMXFiles(baseDir)
	if err != nil {
		return fmt.Errorf("error loading maps: %s %w", baseDir, err)
	}

	for _, tmxFile := range tmxFiles {
		t, err := tmx.LoadFile(tmxFile)
		if err != nil {
			return fmt.Errorf("error loading maps: %s %w", baseDir, err)
		}

		mm.Maps[t.Class] = t
	}

	return nil
}

func findTMXFiles(dir string) ([]string, error) {
	var tmxFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".tmx" {
			tmxFiles = append(tmxFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return tmxFiles, nil
}
