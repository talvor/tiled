package tiled

import (
	anim "github.com/talvor/tiled/animation/manager"
	anir "github.com/talvor/tiled/animation/renderer"
	tmxm "github.com/talvor/tiled/tmx/manager"
	tmxr "github.com/talvor/tiled/tmx/renderer"
	tsxm "github.com/talvor/tiled/tsx/manager"
	tsxr "github.com/talvor/tiled/tsx/renderer"
)

func NewAnimationRenderer(animationBaseDirs []string, tilesetBaseDirs []string) *anir.Renderer {
	am := anim.NewManager(animationBaseDirs)
	tsm := tsxm.NewManager(tilesetBaseDirs)
	tsr := tsxr.NewRenderer(tsm)

	return anir.NewRenderer(am, tsr)
}

func NewTilesetRenderer(tilesetBaseDirs []string) *tsxr.Renderer {
	ts := tsxm.NewManager(tilesetBaseDirs)
	return tsxr.NewRenderer(ts)
}

func NewMapRenderer(mapBaseDirs []string, tilesetBaseDirs []string) *tmxr.Renderer {
	ts := tsxm.NewManager(tilesetBaseDirs)

	mm := tmxm.NewManager(mapBaseDirs)

	tsr := tsxr.NewRenderer(ts)
	return tmxr.NewRenderer(mm, tsr)
}
