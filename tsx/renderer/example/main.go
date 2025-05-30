package main

import (
	"image/color"
	"log"
	"path"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/talvor/tiled/common"
	"github.com/talvor/tiled/tsx/manager"
	"github.com/talvor/tiled/tsx/renderer"
)

var (
	r                *renderer.Renderer
	simpleSprite     *renderer.SimpleSprite
	runningAnimation *renderer.SimpleAnimation
	attackAnimation  *renderer.SimpleAnimation
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	tilesetDir := path.Join(path.Dir(filename), "tilesets")
	tm := manager.NewManager([]string{tilesetDir})

	r = renderer.NewRenderer(tm)
	simpleSprite = renderer.NewSimpleSprite("player", r)
	runningAnimation = renderer.NewSimpleAnimation(simpleSprite, 100, []int{}, nil)
	attackAnimation = renderer.NewSimpleAnimation(simpleSprite, 150, []int{}, nil)
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(16, 16)

	// Draw all tiles from tileset
	ts := r.TilesetManager.GetTilesetByName(simpleSprite.Tileset)
	if ts == nil {
		log.Fatal("Tileset not found")
	}
	columns := ts.Columns
	rows := ts.TileCount / columns

	for row := range rows {
		for column := range columns {
			if err := simpleSprite.Draw((row*columns)+column, &common.DrawOptions{Screen: screen, Op: op}); err != nil {
				panic(err)
			}
			op.GeoM.Translate(32, 0)
		}

		op.GeoM.Translate(float64(columns*-32), 48)
	}

	op.GeoM.Reset()
	op.GeoM.Translate(384, 16)

	// Draw running animations
	// Draw running down
	runningAnimation.SetFrames([]int{18, 19, 20, 21, 22, 23})
	runningAnimation.DrawAnimation(&common.DrawOptions{Screen: screen, Op: op})
	op.GeoM.Translate(32, 0)
	// Draw running right
	runningAnimation.SetFrames([]int{24, 25, 26, 27, 28, 29})
	runningAnimation.DrawAnimation(&common.DrawOptions{Screen: screen, Op: op})
	op.GeoM.Translate(32, 0)
	// Draw running right
	runningAnimation.DrawAnimation(&common.DrawOptions{Screen: screen, Op: op, FlipHorizontal: true})
	op.GeoM.Translate(32, 0)
	// Draw running up
	runningAnimation.SetFrames([]int{30, 31, 32, 33, 34, 35})
	runningAnimation.DrawAnimation(&common.DrawOptions{Screen: screen, Op: op})

	op.GeoM.Translate(-96, 48)

	// Draw attack animations
	attackAnimation.SetFrames([]int{36, 37, 38, 39})
	attackAnimation.DrawAnimation(&common.DrawOptions{Screen: screen, Op: op})
	op.GeoM.Translate(32, 0)
	attackAnimation.SetFrames([]int{42, 43, 44, 45})
	attackAnimation.DrawAnimation(&common.DrawOptions{Screen: screen, Op: op})
	op.GeoM.Translate(32, 0)
	attackAnimation.DrawAnimation(&common.DrawOptions{Screen: screen, Op: op, FlipHorizontal: true})
	op.GeoM.Translate(32, 0)
	attackAnimation.SetFrames([]int{48, 49, 50, 51})
	attackAnimation.DrawAnimation(&common.DrawOptions{Screen: screen, Op: op})
	op.GeoM.Translate(32, 0)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 512
}

func main() {
	ebiten.SetWindowSize(640, 512)
	ebiten.SetWindowTitle("Render an image")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
