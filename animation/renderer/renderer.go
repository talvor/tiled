package renderer

import "github.com/talvor/tiled/animation/manager"

type Renderer struct {
	AnimationManager *manager.AnimationManager
}

func NewRenderer(am *manager.AnimationManager) *Renderer {
	return &Renderer{
		AnimationManager: am,
	}
}
