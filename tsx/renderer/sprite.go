package renderer

import (
	"errors"
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/talvor/tiled/tsx"
)

var (
	ErrInvalidIdType     = errors.New("invalid id type")
	ErrTileset           = errors.New("error loading tileset")
	ErrNoAmimationFrames = errors.New("no animation frames")
	ErrNoComplexPart     = errors.New("no part found for complex sprite")
)

type DrawOptions struct {
	Screen         *ebiten.Image
	Op             *ebiten.DrawImageOptions
	FlipHorizontal bool
	FlipVertical   bool
}

type SpriteDrawer interface {
	Draw(id interface{}, opts *DrawOptions) error
}

// A simple sprite is a sprite that is made up of a single tileset
// and can be drawn with a single tile id
type SimpleSprite struct {
	Tileset  string
	Renderer *Renderer
}

func NewSimpleSprite(tileSet string, renderer *Renderer) *SimpleSprite {
	return &SimpleSprite{
		Tileset:  tileSet,
		Renderer: renderer,
	}
}

func (ss *SimpleSprite) Draw(id interface{}, opts *DrawOptions) error {
	switch id.(type) {
	case int:
		return drawSpriteByID(ss.Tileset, uint32(id.(int)), ss.Renderer, opts)
	case uint32:
		return drawSpriteByID(ss.Tileset, id.(uint32), ss.Renderer, opts)
	case string:
		return drawSpriteByName(ss.Tileset, id.(string), ss.Renderer, opts)
	}
	return fmt.Errorf("invalid id type: %w", ErrInvalidIdType)
}

func (ss *SimpleSprite) DrawWithAnimation(name string, duration int, opts *DrawOptions) error {
	return drawSpriteWithAnimation(ss.Tileset, name, duration, ss.Renderer, opts)
}

// A compound sprite is a sprite that is made up of multiple tilesets
// and can be drawn with a single tile id
type CompoundSprite struct {
	Tilesets []string
	Renderer *Renderer
}

func NewCompoundSprite(tileSets []string, renderer *Renderer) *CompoundSprite {
	return &CompoundSprite{
		Tilesets: tileSets,
		Renderer: renderer,
	}
}

func (cs *CompoundSprite) Draw(id interface{}, opts *DrawOptions) error {
	switch id.(type) {
	case int:
		for _, tileset := range cs.Tilesets {
			if err := drawSpriteByID(tileset, uint32(id.(int)), cs.Renderer, opts); err != nil {
				return err
			}
		}
	case uint32:
		for _, tileset := range cs.Tilesets {
			if err := drawSpriteByID(tileset, id.(uint32), cs.Renderer, opts); err != nil {
				return err
			}
		}
	case string:
		for _, tileset := range cs.Tilesets {
			if err := drawSpriteByName(tileset, id.(string), cs.Renderer, opts); err != nil {
				return err
			}
		}
	default:
		return ErrInvalidIdType
	}
	return nil
}

func (cs *CompoundSprite) DrawWithAnimation(name string, duration int, opts *DrawOptions) error {
	for _, tileset := range cs.Tilesets {
		if err := drawSpriteWithAnimation(tileset, name, duration, cs.Renderer, opts); err != nil {
			return err
		}
	}
	return nil
}

// A complex sprite has multiple parts that can be drawn with a single tile id
type ComplexSprite struct {
	Tileset  string
	Parts    map[uint32][]uint32
	Columns  int
	Renderer *Renderer
}

func NewComplexSprite(tileSet string, columns int, renderer *Renderer) *ComplexSprite {
	parts := make(map[uint32][]uint32)
	return &ComplexSprite{
		Tileset:  tileSet,
		Parts:    parts,
		Columns:  columns,
		Renderer: renderer,
	}
}

func (cs *ComplexSprite) AddPart(id uint32, parts []uint32) {
	cs.Parts[id] = parts
}

func (cs *ComplexSprite) Draw(id interface{}, opts *DrawOptions) error {
	ts, err := cs.Renderer.TilesetManager.GetTilesetByName(cs.Tileset)
	if err != nil {
		return fmt.Errorf("failed to find tileset with name %s: %w", cs.Tileset, ErrTileset)
	}

	tileId := uint32(0)
	switch id.(type) {
	case int:
		tileId = uint32(id.(int))
	case uint32:
		tileId = id.(uint32)
	default:
		return ErrInvalidIdType
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(opts.Op.GeoM)
	drawOptions := &DrawOptions{
		Screen:         opts.Screen,
		Op:             op,
		FlipHorizontal: opts.FlipHorizontal,
		FlipVertical:   opts.FlipVertical,
	}

	parts, ok := cs.Parts[tileId]
	if !ok {
		return ErrNoComplexPart
	}

	dx := float64(0)
	for idx, part := range parts {
		if err := drawSpriteByID(cs.Tileset, part, cs.Renderer, drawOptions); err != nil {
			return err
		}
		if (idx+1)%cs.Columns == 0 {
			drawOptions.Op.GeoM.Translate(-dx, float64(ts.TileHeight))
			dx = 0
		} else {
			drawOptions.Op.GeoM.Translate(float64(ts.TileWidth), 0)
			dx += float64(ts.TileWidth)
		}
	}
	return nil
}

func getTileset(tileset string, tilesetManager *tsx.TilesetManager) (*tsx.Tileset, error) {
	return tilesetManager.GetTilesetByName(tileset)
}

func drawSpriteByID(
	tileset string,
	ID uint32,
	renderer *Renderer,
	opts *DrawOptions,
) error {
	ts, err := renderer.TilesetManager.GetTilesetByName(tileset)
	if err != nil {
		return fmt.Errorf("failed to find tileset with name %s: %w", tileset, ErrTileset)
	}

	img, err := renderer.loadTilesetImage(ts)
	if err != nil {
		return fmt.Errorf("failed to load tileset image for tileset %s: %w", tileset, ErrTileset)
	}

	rect, err := ts.GetTileRect(ID)
	if err != nil {
		return fmt.Errorf("failed to get tile rect for tile %d in tileset %s: %w", ID, tileset, ErrTileset)
	}

	img = transformImage(img.SubImage(rect).(*ebiten.Image), opts)

	opts.Screen.DrawImage(img, opts.Op)

	return nil
}

func drawSpriteByName(tileset string, name string, er *Renderer, opts *DrawOptions) error {
	tile, err := getTileByName(tileset, name, er)
	if err != nil {
		return err
	}

	return drawSpriteByID(tileset, tile.ID, er, opts)
}

func drawSpriteWithAnimation(tileset string, name string, duration int, er *Renderer, opts *DrawOptions) error {
	tile, err := getTileByName(tileset, name, er)
	if err != nil {
		return err
	}

	if tile.Animation.Frames == nil || len(tile.Animation.Frames) == 0 {
		return fmt.Errorf("no animation frames found for tile %s in tileset %s: %w", name, tileset, ErrNoAmimationFrames)
	}

	tileID := tile.ID
	if tile.Animation.Frames != nil && len(tile.Animation.Frames) > 0 {
		animationIdx := int(time.Now().UnixMilli()) / duration % len(tile.Animation.Frames)
		frame := tile.Animation.Frames[animationIdx]
		tileID = frame.ID
	}
	return drawSpriteByID(tileset, tileID, er, opts)
}

func getTileByName(tileset string, name string, er *Renderer) (*tsx.Tile, error) {
	ts, err := er.TilesetManager.GetTilesetByName(tileset)
	if err != nil {
		return nil, fmt.Errorf("failed to find tileset with name %s: %w", tileset, ErrTileset)
	}

	tile, err := ts.GetTileByType(name)
	if err != nil {
		return nil, fmt.Errorf("failed to find tile with name %s in tileset %s: %w", name, tileset, ErrTileset)
	}

	return tile, nil
}

func transformImage(img *ebiten.Image, opts *DrawOptions) *ebiten.Image {
	if opts.FlipHorizontal {
		result := ebiten.NewImage(img.Bounds().Dx(), img.Bounds().Dy())
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(img.Bounds().Dx()), 0)
		result.DrawImage(img, op)
		img = result
	}
	if opts.FlipVertical {
		result := ebiten.NewImage(img.Bounds().Dx(), img.Bounds().Dy())
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1, -1)
		op.GeoM.Translate(0, float64(img.Bounds().Dy()))
		result.DrawImage(img, op)
		img = result
	}

	return img
}
