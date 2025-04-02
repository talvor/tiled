package tiled

import (
	tmxm "github.com/talvor/tiled/tmx/manager"
	tmxr "github.com/talvor/tiled/tmx/renderer"
	tsxm "github.com/talvor/tiled/tsx/manager"
	tsxr "github.com/talvor/tiled/tsx/renderer"
)

func NewTilesetRenderer(tilesetBaseDir string) (*tsxr.Renderer, error) {
	ts, err := tsxm.NewManager(tilesetBaseDir)
	if err != nil {
		return nil, err
	}
	return tsxr.NewRenderer(ts), nil
}

func NewMapRenderer(mapBaseDir string, tilesetBaseDir string) (*tmxr.Renderer, error) {
	ts, err := tsxm.NewManager(tilesetBaseDir)
	if err != nil {
		return nil, err
	}

	mm, err := tmxm.NewManager(mapBaseDir)
	if err != nil {
		return nil, err
	}

	tsr := tsxr.NewRenderer(ts)
	return tmxr.NewRenderer(mm, tsr), nil
}
