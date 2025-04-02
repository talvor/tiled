package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	tmxmanager "github.com/talvor/tiled/tmx/manager"
	tmxrenderer "github.com/talvor/tiled/tmx/renderer"
	tsxmanager "github.com/talvor/tiled/tsx/manager"
	tsxrenderer "github.com/talvor/tiled/tsx/renderer"
)

var tmxr *tmxrenderer.Renderer

func init() {
	tm, err := tsxmanager.NewManager("./assets/")
	if err != nil {
		log.Fatal(err)
	}
	mm, err := tmxmanager.NewManager("./assets/")
	if err != nil {
		log.Fatal(err)
	}

	tsxr := tsxrenderer.NewRenderer(tm)
	tmxr = tmxrenderer.NewRenderer(mm, tsxr)
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := &ebiten.DrawImageOptions{}
	op := &tmxrenderer.DrawOptions{Screen: screen, Op: opts}

	tmxr.DrawMapLayer("StartScene", "background", op)
	tmxr.DrawMapLayer("StartScene", "bottom", op)
	tmxr.DrawMapLayer("StartScene", "top", op)
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
