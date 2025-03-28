package tiled

import (
	"github.com/talvor/tiled/tmx"
	tmxrenderer "github.com/talvor/tiled/tmx/renderer"
	"github.com/talvor/tiled/tsx"
	tsxrenderer "github.com/talvor/tiled/tsx/renderer"
)

func NewTilesetRenderer(tilesetBaseDir string) *tsxrenderer.Renderer {
	ts := tsx.NewTilesetManager(tilesetBaseDir)
	return tsxrenderer.NewRenderer(ts)
}

func NewMapRenderer(mapBaseDir string, tilesetBaseDir string) *tmxrenderer.Renderer {
	ts := tsx.NewTilesetManager(tilesetBaseDir)
	mm := tmx.NewMapManager(mapBaseDir)
	tsr := tsxrenderer.NewRenderer(ts)
	return tmxrenderer.NewRenderer(mm, tsr)
}
