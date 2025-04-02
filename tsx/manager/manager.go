package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/talvor/tiled/tsx"
)

var (
	ErrTilesetManagerNotLoaded = errors.New("tsx: tileset manager not loaded")
	ErrTilesetNotFound         = errors.New("tsx: tileset not found")
)

type TilesetManager struct {
	baseDir  string
	Tilesets map[string]*tsx.Tileset
	IsLoaded bool
}

func (tm *TilesetManager) GetTilesetBySource(source string) (*tsx.Tileset, error) {
	if !tm.IsLoaded {
		return nil, ErrTilesetManagerNotLoaded
	}
	for _, ts := range tm.Tilesets {
		if ts.Source == source {
			return ts, nil
		}
	}
	return nil, fmt.Errorf("source: %s %w", source, ErrTilesetNotFound)
}

func (tm *TilesetManager) GetTilesetByName(name string) (*tsx.Tileset, error) {
	if !tm.IsLoaded {
		return nil, ErrTilesetManagerNotLoaded
	}

	if _, ok := tm.Tilesets[name]; !ok {
		return nil, fmt.Errorf("name: %s %w", name, ErrTilesetNotFound)
	}

	return tm.Tilesets[name], nil
}

func (tm *TilesetManager) AddTileset(source string) error {
	ts, err := tsx.LoadFile(source)
	if err != nil {
		return err
	}

	tm.Tilesets[ts.Name] = ts

	return nil
}

func (tm *TilesetManager) LoadTilesetsFromDir(dir string) {
	loadTilesets(tm, dir)
}

func (tm *TilesetManager) DebugPrintTilesets() {
	for name := range tm.Tilesets {
		fmt.Println(name)
	}
}

func NewManager(baseDir string) (*TilesetManager, error) {
	tm := &TilesetManager{
		baseDir:  baseDir,
		Tilesets: make(map[string]*tsx.Tileset),
		IsLoaded: false,
	}

	if err := loadTilesets(tm, baseDir); err != nil {
		return nil, err
	}

	return tm, nil
}

func loadTilesets(tm *TilesetManager, baseDir string) error {
	tsxFiles, err := findTSXFiles(baseDir)
	if err != nil {
		return fmt.Errorf("error loading tilesets: %s %w", baseDir, err)
	}

	for _, tsxFile := range tsxFiles {
		ts, err := tsx.LoadFile(tsxFile)
		if err != nil {
			return fmt.Errorf("error loading tilesets: %s %w", baseDir, err)
		}

		tm.Tilesets[ts.Name] = ts
	}

	tm.IsLoaded = true
	return nil
}

func findTSXFiles(dir string) ([]string, error) {
	var tsxFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".tsx" {
			tsxFiles = append(tsxFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return tsxFiles, nil
}
