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
	ErrTilesetGroupNotFound    = errors.New("tsx: tileset group not found")
)

type TilesetManager struct {
	baseDir          string
	TilesetsByName   map[string]*tsx.Tileset
	TilesetsBySource map[string]*tsx.Tileset
	TilesetGroups    map[string]tsx.TilesetGroup
}

func (tm *TilesetManager) GetTilesetBySource(source string) (*tsx.Tileset, error) {
	if _, ok := tm.TilesetsBySource[source]; !ok {
		return nil, fmt.Errorf("source: %s %w", source, ErrTilesetNotFound)
	}

	return tm.TilesetsBySource[source], nil
}

func (tm *TilesetManager) GetTilesetByName(name string) (*tsx.Tileset, error) {
	if _, ok := tm.TilesetsByName[name]; !ok {
		return nil, fmt.Errorf("name: %s %w", name, ErrTilesetNotFound)
	}

	return tm.TilesetsByName[name], nil
}

func (tm *TilesetManager) HasTilesetByName(name string) bool {
	if _, ok := tm.TilesetsByName[name]; !ok {
		return false
	}

	return true
}

func (tm *TilesetManager) HasTilesetBySource(source string) bool {
	if _, ok := tm.TilesetsBySource[source]; !ok {
		return false
	}

	return true
}

func (tm *TilesetManager) AddTileset(sourcePath string) (*tsx.Tileset, error) {
	ts, err := tsx.LoadFile(sourcePath)
	if err != nil {
		return nil, err
	}

	tm.TilesetsByName[ts.Name] = ts
	tm.TilesetsBySource[ts.Source] = ts

	return ts, nil
}

func (tm *TilesetManager) AddTilesetGroupBySource(name string, sources []string) error {
	var tilesets tsx.TilesetGroup
	for _, source := range sources {
		var ts *tsx.Tileset
		ts, _ = tm.GetTilesetBySource(source)
		if ts == nil {
			var err error
			ts, err = tm.AddTileset(source)
			if err != nil {
				return err
			}
		}
		tilesets = append(tilesets, ts)
	}

	tm.TilesetGroups[name] = tilesets

	return nil
}

func (tm *TilesetManager) AddTilesetGroup(name string, names []string) error {
	var tilesets tsx.TilesetGroup
	for _, name := range names {
		ts, err := tm.GetTilesetByName(name)
		if err != nil {
			return err
		}
		tilesets = append(tilesets, ts)
	}

	tm.TilesetGroups[name] = tilesets

	return nil
}

func (tm *TilesetManager) GetTilesetGroup(name string) (tsx.TilesetGroup, error) {
	if _, ok := tm.TilesetGroups[name]; !ok {
		return nil, fmt.Errorf("group: %s %w", name, ErrTilesetGroupNotFound)
	}

	return tm.TilesetGroups[name], nil
}

func (tm *TilesetManager) LoadTilesetsFromDir(dir string) {
	loadTilesets(tm, dir)
}

func (tm *TilesetManager) DebugPrintTilesets() {
	for name := range tm.TilesetsByName {
		fmt.Println(name)
	}
}

func NewManager(baseDir string) (*TilesetManager, error) {
	tm := &TilesetManager{
		baseDir:          baseDir,
		TilesetsByName:   make(map[string]*tsx.Tileset),
		TilesetsBySource: make(map[string]*tsx.Tileset),
		TilesetGroups:    make(map[string]tsx.TilesetGroup),
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

		tm.TilesetsByName[ts.Name] = ts
		tm.TilesetsBySource[ts.Source] = ts
	}

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
