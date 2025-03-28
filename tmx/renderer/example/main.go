package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/talvor/tiled/tmx"
	tmxrenderer "github.com/talvor/tiled/tmx/renderer"
	"github.com/talvor/tiled/tsx"
	tsxrenderer "github.com/talvor/tiled/tsx/renderer"
)

var tmxr *tmxrenderer.Renderer

func init() {
	tm := tsx.NewTilesetManager("./assets/")
	mm := tmx.NewMapManager("./assets/")

	tsxr := tsxrenderer.NewRenderer(tm)
	tmxr = tmxrenderer.NewRenderer(mm, tsxr)
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	tmxr.DrawMapLayer("StartScene", "background", screen)
	tmxr.DrawMapLayer("StartScene", "bottom", screen)
	tmxr.DrawMapLayer("StartScene", "top", screen)
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
