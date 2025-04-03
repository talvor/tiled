package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/talvor/tiled/common"
	"github.com/talvor/tiled/tmx/manager"
	tsxrenderer "github.com/talvor/tiled/tsx/renderer"
)

type Renderer struct {
	TsxRenderer *tsxrenderer.Renderer
	MapManager  *manager.MapManager
}

func NewRenderer(mm *manager.MapManager, tsxRenderer *tsxrenderer.Renderer) *Renderer {
	return &Renderer{
		TsxRenderer: tsxRenderer,
		MapManager:  mm,
	}
}

func (r *Renderer) DrawMapLayer(mapName string, layerName string, opts *common.DrawOptions) error {
	m, err := r.MapManager.GetMapByName(mapName)
	if err != nil {
		return err
	}

	layer, err := m.GetLayer(layerName)
	if err != nil {
		return err
	}

	for idx, tileId := range layer.Tiles {
		ts, id := m.DecodeTileGID(tileId)
		if ts == nil {
			continue
		}

		posX, posY := layer.GetTilePositionFromIndex(idx, m)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Concat(opts.Op.GeoM)
		op.GeoM.Translate(float64(posX), float64(posY))

		if err := r.TsxRenderer.DrawTileWithSource(ts.Source, uint32(id), &common.DrawOptions{
			Screen: opts.Screen,
			Op:     op,
		}); err != nil {
			return err
		}
	}
	return nil
}
