package main

import (
	"image/color"
	"log"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	anim "github.com/talvor/tiled/animation/manager"
	"github.com/talvor/tiled/animation/renderer"
	"github.com/talvor/tiled/common"
	tsxm "github.com/talvor/tiled/tsx/manager"
	tsxr "github.com/talvor/tiled/tsx/renderer"
)

var anir *renderer.Renderer

func init() {
	homeDir, _ := os.UserHomeDir()
	tilesetsDir := path.Join(homeDir, "Documents/tilesets/Cute_Fantasy/sprites")
	animationsDir := path.Join(homeDir, "Documents/tilesets/Cute_Fantasy/animations")

	tsxm := tsxm.NewManager([]string{tilesetsDir})

	tsxr := tsxr.NewRenderer(tsxm)
	anim := anim.NewManager([]string{animationsDir})
	anir = renderer.NewRenderer(anim, tsxr)

	anir.SetTilesetResolver(tsxr.TilesetManager)
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

	// Draw animation
	collider := anir.GetCollider("simple_player", "walking", "collider")
	cx := float64(collider.Min.X) + dx
	cy := float64(collider.Min.Y) + dy
	ebitenutil.DrawRect(screen, cx, cy, float64(collider.Dx()), float64(collider.Dy()), color.RGBA{255, 0, 0, 25})
	panicOnError(anir.Draw("simple_player", "walking", &common.DrawOptions{
		Screen: screen,
		Op:     op,
	}))

	moveRight(48, 0)

	panicOnError(anir.Draw("simple_enemy", "running", &common.DrawOptions{
		Screen: screen,
		Op:     op,
	}))

	moveRight(48, 0)

	panicOnError(anir.Draw("timed_player", "idle", &common.DrawOptions{
		Screen: screen,
		Op:     op,
	}))

	moveRight(48, 0)

	collider = anir.GetCollider("complex_player", "chop", "collider")
	cx = float64(collider.Min.X) + dx
	cy = float64(collider.Min.Y) + dy - 8
	ebitenutil.DrawRect(screen, cx, cy, float64(collider.Dx()), float64(collider.Dy()), color.RGBA{255, 0, 0, 25})
	panicOnError(anir.Draw("complex_player", "chop", &common.DrawOptions{
		Screen: screen,
		Op:     op,
	}))

	moveRight(48, 0)

	nextLine(48)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Render animations")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
