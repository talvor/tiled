package renderer

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/talvor/tiled/animation"
	"github.com/talvor/tiled/animation/manager"
	"github.com/talvor/tiled/common"
	"github.com/talvor/tiled/tsx"
	tsxr "github.com/talvor/tiled/tsx/renderer"
)

type TilesetResolver interface {
	GetTilesets(names []string) []*tsx.Tileset
	GetTilesetByName(name string) *tsx.Tileset
}

type Renderer struct {
	TSXRenderer      *tsxr.Renderer
	AnimationManager *manager.AnimationManager
	TilesetResolver  TilesetResolver
}

func NewRenderer(am *manager.AnimationManager, r *tsxr.Renderer) *Renderer {
	return &Renderer{
		AnimationManager: am,
		TSXRenderer:      r,
	}
}

func (r *Renderer) SetTilesetResolver(resolver TilesetResolver) {
	r.TilesetResolver = resolver
}

func (r *Renderer) Draw(class string, action string, opts *common.DrawOptions) error {
	ani, err := r.AnimationManager.GetAnimation(class, action)
	if err != nil {
		return err
	}
	frame := ani.GetNextFrame()

	drawFunc := func(part animation.Part) error {
		opts.OffsetX = float64(part.XOffset)
		opts.OffsetY = float64(part.YOffset)
		opts.FlipHorizontal = part.FlipHorizontal
		opts.FlipVertical = part.FlipVertical

		if part.Tileset != "" {
			if err := r.drawTile(part.Tileset, part.TileID, opts); err != nil {
				return err
			}
		} else {
			for _, tileset := range ani.Tilesets {
				if err := r.drawTile(tileset, part.TileID, opts); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return r.processFrame(&frame, drawFunc)
}

func (r *Renderer) processFrame(frame *animation.Frame, processFunc func(part animation.Part) error) error {
	for _, part := range frame.Parts {
		if part.TileID == -1 {
			continue
		}
		if err := processFunc(part); err != nil {
			return fmt.Errorf("failed to process part: %w", err)
		}
	}
	return nil
}

func (r *Renderer) GetCollider(class string, action string, colliderName string) *image.Rectangle {
	// If the tileset resolver is not set, we cannot resolve tilesets to get the colliders
	if r.TilesetResolver == nil {
		return nil
	}
	// Get the animation
	ani, err := r.AnimationManager.GetAnimation(class, action)
	if err != nil {
		return nil
	}

	if ani.ColliderRect == nil {
		collisionRect := image.Rectangle{}
		// Iterate through the frames and get the collision rectangles
		for _, frame := range ani.Frames {
			// Iterate through the parts of the frame
			r.processFrame(&frame, func(part animation.Part) error {
				if part.TileID == -1 {
					return nil
				}
				if part.Tileset != "" {
					ts := r.TilesetResolver.GetTilesetByName(part.Tileset)
					rect, _ := ts.GetTileCollisionRect(uint32(part.TileID), colliderName)
					collisionRect = collisionRect.Union(rect)
				} else {
					for _, tileset := range ani.Tilesets {
						ts := r.TilesetResolver.GetTilesetByName(tileset)
						rect, _ := ts.GetTileCollisionRect(uint32(part.TileID), colliderName)
						collisionRect = collisionRect.Union(rect)
					}
				}
				return nil
			})
		}

		// Save the collision rectangle to the animation for future use
		ani.ColliderRect = &collisionRect
	}

	return ani.ColliderRect
}

func (r *Renderer) drawTile(tileset string, tileID int, opts *common.DrawOptions) error {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(opts.Op.GeoM)
	op.GeoM.Translate(opts.OffsetX, opts.OffsetY)
	drawOptions := &common.DrawOptions{
		Screen:         opts.Screen,
		Op:             op,
		FlipHorizontal: opts.FlipHorizontal,
		FlipVertical:   opts.FlipVertical,
	}
	if err := r.TSXRenderer.DrawTileWithName(tileset, uint32(tileID), drawOptions); err != nil {
		return fmt.Errorf("failed to draw tile %s: %w", tileset, err)
	}
	return nil
}
