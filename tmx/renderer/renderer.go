package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/talvor/tiled/tmx"
	tsxrenderer "github.com/talvor/tiled/tsx/renderer"
)

type Renderer struct {
	TsxRenderer *tsxrenderer.Renderer
	MapManager  *tmx.MapManager
}

func NewRenderer(mm *tmx.MapManager, tsxRenderer *tsxrenderer.Renderer) *Renderer {
	return &Renderer{
		TsxRenderer: tsxRenderer,
		MapManager:  mm,
	}
}

func (r *Renderer) DrawMapLayer(mapName string, layerName string, screen *ebiten.Image) error {
	m, err := r.MapManager.GetMapByName(mapName)
	if err != nil {
		return err
	}

	layer, err := m.GetLayer(layerName)
	if err != nil {
		return err
	}

	for idx, tileId := range layer.Tiles {
		// for i := 0; i < 500; i++ {
		// tileId := layer.Tiles[i]
		ts, id := m.DecodeTileGID(tileId)
		if ts == nil {
			continue
		}

		posX, posY := layer.GetTilePositionFromIndex(idx, m)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(posX), float64(posY))
		if err := r.TsxRenderer.DrawTileWithSource(ts.Source, uint32(id), &tsxrenderer.DrawOptions{
			Screen: screen,
			Op:     op,
		}); err != nil {
			return err
		}
	}
	return nil
}

type DrawOptions struct {
	Screen         *ebiten.Image
	Op             *ebiten.DrawImageOptions
	FlipHorizontal bool
	FlipVertical   bool
}
