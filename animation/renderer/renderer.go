package renderer

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/talvor/tiled/animation/manager"
	"github.com/talvor/tiled/common"
	tsxr "github.com/talvor/tiled/tsx/renderer"
)

type Renderer struct {
	TSXRenderer      *tsxr.Renderer
	AnimationManager *manager.AnimationManager
}

func NewRenderer(am *manager.AnimationManager, r *tsxr.Renderer) *Renderer {
	return &Renderer{
		AnimationManager: am,
		TSXRenderer:      r,
	}
}

func (r *Renderer) Draw(class string, action string, opts *common.DrawOptions) error {
	ani, err := r.AnimationManager.GetAnimation(class, action)
	if err != nil {
		return err
	}
	frame := ani.GetCurrentFrame()
	for _, part := range frame.Parts {
		if part.TileID == -1 {
			continue
		}

		opts.OffsetX = float64(part.XOffset)
		opts.OffsetY = float64(part.YOffset)
		opts.FlipHorizontal = part.FlipHorizontal
		opts.FlipVertical = part.FlipVertical

		if part.Sprite != "" {
			if err := r.drawTile(part.Sprite, part.TileID, opts); err != nil {
				return err
			}
		} else {
			for _, sprite := range ani.Sprites {
				if err := r.drawTile(sprite, part.TileID, opts); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *Renderer) drawTile(name string, tileID int, opts *common.DrawOptions) error {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(opts.Op.GeoM)
	op.GeoM.Translate(opts.OffsetX, opts.OffsetY)
	drawOptions := &common.DrawOptions{
		Screen:         opts.Screen,
		Op:             op,
		FlipHorizontal: opts.FlipHorizontal,
		FlipVertical:   opts.FlipVertical,
	}
	if err := r.TSXRenderer.DrawTileWithName(name, uint32(tileID), drawOptions); err != nil {
		return fmt.Errorf("failed to draw tile %s: %w", name, err)
	}
	return nil
}
