package main

import (
	"image"
	"image/color"
	"log"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/talvor/tiled/common"
	"github.com/talvor/tiled/tsx"
	"github.com/talvor/tiled/tsx/manager"
	"github.com/talvor/tiled/tsx/renderer"
)

var (
	r              *renderer.Renderer
	walkAnimation  *renderer.SimpleAnimation
	runAnimation   *renderer.TimedAnimation
	compoundSprite *renderer.CompoundSprite
	bodyTileset    *tsx.Tileset
	collisionRect1 image.Rectangle
	collisionRect2 image.Rectangle
)

func init() {
	homeDir, _ := os.UserHomeDir()
	tilesetsDir := path.Join(homeDir, "Downloads/mana_seed_character_base/character")
	tm := manager.NewManager([]string{tilesetsDir})
	bodyTileset, _ = tm.GetTilesetByName("char_a_p1_0bas_humn_v01")
	collisionRect1, _ = bodyTileset.GetTileCollisionRect(48, "collider")
	collisionRect2, _ = bodyTileset.GetTileCollisionRect(49, "full_tile_collider")
	r = renderer.NewRenderer(tm)
	compoundSprite = renderer.NewCompoundSprite([]string{
		"char_a_p1_0bas_humn_v01",
		// "char_a_p1_1out_boxr_v01",
		// "char_a_p1_1out_undi_v01",
		"char_a_p1_1out_pfpn_v04",
		"char_a_p1_4har_dap1_v01",
		"char_a_p1_5hat_pnty_v04",
	}, r)
	walkAnimation = renderer.NewSimpleAnimation(compoundSprite, 135, []int{48, 49, 50, 51, 52, 53}, nil)
	runAnimation, _ = renderer.NewTimedAnimation(compoundSprite, []int{48, 49, 54, 51, 52, 55}, []uint32{80, 55, 125, 80, 55, 125}, nil)
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	op := &ebiten.DrawImageOptions{}

	// Draw simple sprite
	dx := float64(0)
	dy := float64(0)

	moveRight := func(x, y float64) {
		dx += x
		dy += y
		op.GeoM.Translate(x, y)
	}

	nextLine := func(y float64) {
		dx = 16
		dy += y
		op.GeoM.Reset()
		op.GeoM.Translate(dx, dy)
	}

	moveRight(16, 16)

	walkAnimation.DrawAnimation(&common.DrawOptions{Screen: screen, Op: op})

	cx := float64(collisionRect1.Min.X) + dx
	cy := float64(collisionRect1.Min.Y) + dy
	ebitenutil.DrawRect(screen, cx, cy, float64(collisionRect1.Dx()), float64(collisionRect1.Dy()), color.RGBA{255, 0, 0, 50})

	nextLine(64)
	runAnimation.DrawAnimation(&common.DrawOptions{Screen: screen, Op: op})
	cx = float64(collisionRect2.Min.X) + dx
	cy = float64(collisionRect2.Min.Y) + dy
	ebitenutil.DrawRect(screen, cx, cy, float64(collisionRect2.Dx()), float64(collisionRect2.Dy()), color.RGBA{255, 0, 0, 50})

	op.GeoM.Reset()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 640
}

func main() {
	ebiten.SetWindowSize(640, 640)
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
