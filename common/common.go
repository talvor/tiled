package common

import "github.com/hajimehoshi/ebiten/v2"

type DrawOptions struct {
	Screen         *ebiten.Image
	Op             *ebiten.DrawImageOptions
	FlipHorizontal bool
	FlipVertical   bool
	OffsetX        float64
	OffsetY        float64
}
