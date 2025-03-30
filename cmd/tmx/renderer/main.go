package main

import (
	"image/color"
	"log"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/talvor/tiled/tmx"
	tmxrenderer "github.com/talvor/tiled/tmx/renderer"
	"github.com/talvor/tiled/tsx"
	tsxrenderer "github.com/talvor/tiled/tsx/renderer"
)

var tmxr *tmxrenderer.Renderer

func init() {
	homeDir, _ := os.UserHomeDir()
	tilesetsDir := path.Join(homeDir, "Documents/tilesets/Cute_Fantasy/tiles")
	mapsDir := path.Join(homeDir, "Documents/tilesets/Cute_Fantasy/maps")

	tm, err := tsx.NewTilesetManager(tilesetsDir)
	panicOnError(err)
	tm.LoadTilesetsFromDir(path.Join(homeDir, "Documents/tilesets/Cute_Fantasy/outdoor_decorations/"))

	mm, err := tmx.NewMapManager(mapsDir)
	panicOnError(err)

	tsxr := tsxrenderer.NewRenderer(tm)
	tmxr = tmxrenderer.NewRenderer(mm, tsxr)
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	op := &ebiten.DrawImageOptions{}

	opts := &tmxrenderer.DrawOptions{
		Screen: screen,
		Op:     op,
	}
	panicOnError(tmxr.DrawMapLayer("GameScene", "background", opts))
	panicOnError(tmxr.DrawMapLayer("GameScene", "bottom", opts))
	panicOnError(tmxr.DrawMapLayer("GameScene", "top", opts))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
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
