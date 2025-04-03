package tsx

import (
	"errors"
	"image"
	"image/color"
	"path"
)

var (
	ErrTileTypeNotFound  = errors.New("tsx: tile type not found")
	ErrTileIDOutOfBounds = errors.New("tsx: tile id out of bounds")
)

type Tileset struct {
	baseDir    string
	Source     string
	Name       string `xml:"name,attr"`
	TileWidth  int    `xml:"tilewidth,attr"`
	TileHeight int    `xml:"tileheight,attr"`
	TileCount  int    `xml:"tilecount,attr"`
	Spacing    int    `xml:"spacing,attr"`
	Margin     int    `xml:"margin,attr"`
	Columns    int    `xml:"columns,attr"`
	Image      Image  `xml:"image"`
	Tiles      []Tile `xml:"tile"`
}

// GetTileRect returns the rectangle of a tile in the tileset.
func (ts *Tileset) GetTileRect(tileID uint32) (image.Rectangle, error) {
	if tileID >= uint32(ts.TileCount) {
		return image.Rectangle{}, ErrTileIDOutOfBounds
	}

	tilesetColumns := ts.Columns

	if tilesetColumns == 0 {
		tilesetColumns = ts.Image.Width / (ts.TileWidth + ts.Spacing)
	}

	x := int(tileID) % tilesetColumns
	y := int(tileID) / tilesetColumns

	xOffset := int(x)*ts.Spacing + ts.Margin
	yOffset := int(y)*ts.Spacing + ts.Margin

	rect := image.Rect(x*ts.TileWidth+xOffset,
		y*ts.TileHeight+yOffset,
		(x+1)*ts.TileWidth+xOffset,
		(y+1)*ts.TileHeight+yOffset)
	return rect, nil
}

// GetTileCollisionRect returns the collision rectangle of a tile in the tileset.
// The collision is determined by an Object inside the tiles ObjectGroup or by the tile itself.
func (ts *Tileset) GetTileCollisionRect(tileID uint32, collisionObjectName string) (image.Rectangle, error) {
	if tileID >= uint32(ts.TileCount) {
		return image.Rectangle{}, ErrTileIDOutOfBounds
	}
	for _, t := range ts.Tiles {
		if t.ID == tileID {
			for _, og := range t.ObjectGroups {
				for _, o := range og.Objects {
					if o.Name == collisionObjectName {
						return image.Rect(int(o.X), int(o.Y), int(o.X+o.Width), int(o.Y+o.Height)), nil
					}
				}
			}
		}
	}
	return image.Rect(0, 0, ts.TileWidth, ts.TileHeight), nil
}

func (ts *Tileset) GetTileByType(tileType string) (*Tile, error) {
	for _, t := range ts.Tiles {
		if t.Type == tileType {
			return &t, nil
		}
	}

	return nil, ErrTileTypeNotFound
}

func (ts *Tileset) GetTileByID(tileID uint32) (*Tile, error) {
	for _, t := range ts.Tiles {
		if t.ID == tileID {
			return &t, nil
		}
	}

	return nil, ErrTileIDOutOfBounds
}

func (ts *Tileset) TileHasAnimation(tileID uint32) bool {
	for _, t := range ts.Tiles {
		if t.ID == tileID && len(t.Animation.Frames) > 0 {
			return true
		}
	}

	return false
}

func (ts *Tileset) GetTileAnimation(tileID uint32) (*Animation, error) {
	for _, t := range ts.Tiles {
		if t.ID == tileID {
			return &t.Animation, nil
		}
	}

	return nil, ErrTileIDOutOfBounds
}

func (ts *Tileset) decodeImage() {
	if ts.Image.Source == "" {
		return
	}

	ts.Image.Source = path.Join(ts.baseDir, ts.Image.Source)
}

type TilesetGroup = []*Tileset

type Image struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type Tile struct {
	ID           uint32         `xml:"id,attr"`
	Type         string         `xml:"type,attr"`
	Animation    Animation      `xml:"animation"`
	ObjectGroups []*ObjectGroup `xml:"objectgroup"`
}

type Animation struct {
	Frames []Frame `xml:"frame"`
}

type Frame struct {
	ID       uint32 `xml:"tileid,attr"`
	Duration int    `xml:"duration,attr"`
}

type ObjectGroup struct {
	ID         uint32     `xml:"id,attr"`
	Name       string     `xml:"name,attr"`
	Class      string     `xml:"class,attr"`
	Color      *HexColor  `xml:"color,attr"`
	Opacity    float32    `xml:"opacity,attr"`
	Visible    bool       `xml:"visible,attr"`
	OffsetX    int        `xml:"offsetx,attr"`
	OffsetY    int        `xml:"offsety,attr"`
	DrawOrder  string     `xml:"draworder,attr"`
	ParallaxX  float32    `xml:"parallaxx,attr"`
	ParallaxY  float32    `xml:"parallaxy,attr"`
	Properties Properties `xml:"properties>property"`
	Objects    []*Object  `xml:"object"`
}

type HexColor struct {
	c color.RGBA
}

type Properties []*Property

type Property struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

type Object struct {
	ID             uint32      `xml:"id,attr"`
	Name           string      `xml:"name,attr"`
	Class          string      `xml:"class,attr"`
	X              float64     `xml:"x,attr"`
	Y              float64     `xml:"y,attr"`
	Width          float64     `xml:"width,attr"`
	Height         float64     `xml:"height,attr"`
	Rotation       float64     `xml:"rotation,attr"`
	GID            uint32      `xml:"gid,attr"`
	Visible        bool        `xml:"visible,attr"`
	Properties     Properties  `xml:"properties>property"`
	Ellipses       []*Ellipse  `xml:"ellipse"`
	Polygons       []*Polygon  `xml:"polygon"`
	PolyLines      []*PolyLine `xml:"polyline"`
	Text           *Text       `xml:"text"`
	TemplateSource string      `xml:"template,attr"`
	TemplateLoaded bool        `xml:"-"`
	Template       *Template
}

type Ellipse struct{}

type Polygon struct {
	Points *Points `xml:"points,attr"`
}

type PolyLine struct {
	Points *Points `xml:"points,attr"`
}

type Point struct {
	X float64
	Y float64
}

type Points []*Point

type Template struct {
	Tileset *Tileset `xml:"tileset"`
	Object  *Object  `xml:"object"`
}

type Text struct {
	Text          string    `xml:",chardata"`
	FontFamily    string    `xml:"fontfamily,attr"`
	Size          int       `xml:"pixelsize,attr"`
	Wrap          bool      `xml:"wrap,attr"`
	Color         *HexColor `xml:"color,attr"`
	Bold          bool      `xml:"bold,attr"`
	Italic        bool      `xml:"italic,attr"`
	Underline     bool      `xml:"underline,attr"`
	Strikethrough bool      `xml:"strikeout,attr"`
	Kerning       bool      `xml:"kerning,attr"`
	HAlign        string    `xml:"halign,attr"`
	VAlign        string    `xml:"valign,attr"`
}
