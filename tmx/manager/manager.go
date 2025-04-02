package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/talvor/tiled/tmx"
)

var (
	ErrMapManagerNotLoaded = errors.New("tsx: map manager not loaded")
	ErrMapNotFound         = errors.New("tsx: map not found")
)

type MapManager struct {
	baseDir  string
	Maps     map[string]*tmx.Map
	IsLoaded bool
}

func (mm *MapManager) GetMapByName(name string) (*tmx.Map, error) {
	if !mm.IsLoaded {
		return nil, ErrMapManagerNotLoaded
	}

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

func NewManager(baseDir string) (*MapManager, error) {
	mm := &MapManager{
		baseDir:  baseDir,
		Maps:     make(map[string]*tmx.Map),
		IsLoaded: false,
	}

	if err := loadMaps(mm, baseDir); err != nil {
		return nil, err
	}

	return mm, nil
}

func loadMaps(mm *MapManager, baseDir string) error {
	tmxFiles, err := findTMXFiles(baseDir)
	if err != nil {
		return fmt.Errorf("error loading maps: %s %w", mm.baseDir, err)
	}

	for _, tmxFile := range tmxFiles {
		t, err := tmx.LoadFile(tmxFile)
		if err != nil {
			return fmt.Errorf("error loading maps: %s %w", mm.baseDir, err)
		}

		mm.Maps[t.Class] = t
	}

	mm.IsLoaded = true
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
