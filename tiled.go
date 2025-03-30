package tiled

import (
	"github.com/talvor/tiled/tmx"
	tmxrenderer "github.com/talvor/tiled/tmx/renderer"
	"github.com/talvor/tiled/tsx"
	tsxrenderer "github.com/talvor/tiled/tsx/renderer"
)

func NewTilesetRenderer(tilesetBaseDir string) (*tsxrenderer.Renderer, error) {
	ts, err := tsx.NewTilesetManager(tilesetBaseDir)
	if err != nil {
		return nil, err
	}
	return tsxrenderer.NewRenderer(ts), nil
}

func NewMapRenderer(mapBaseDir string, tilesetBaseDir string) (*tmxrenderer.Renderer, error) {
	ts, err := tsx.NewTilesetManager(tilesetBaseDir)
	if err != nil {
		return nil, err
	}

	mm, err := tmx.NewMapManager(mapBaseDir)
	if err != nil {
		return nil, err
	}

	tsr := tsxrenderer.NewRenderer(ts)
	return tmxrenderer.NewRenderer(mm, tsr), nil
}
