package main

import (
	"image/color"
	"log"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/talvor/tiled/tsx"
	"github.com/talvor/tiled/tsx/renderer"
)

var (
	r                *renderer.Renderer
	simpleSprite     *renderer.SimpleSprite
	compoundSprite   *renderer.CompoundSprite
	attackAnimation  *renderer.SimpleAnimation
	runningAnimation *renderer.SimpleAnimation
	complexSprite    *renderer.ComplexSprite
	complexAnimation *renderer.SimpleAnimation
)

func init() {
	homeDir, _ := os.UserHomeDir()
	tilesetsDir := path.Join(homeDir, "Documents/examples/tilesets")
	otherTilesetDir := path.Join(homeDir, "Documents/tilesets/Cute_Fantasy/")
	tm := tsx.NewTilesetManager(tilesetsDir)
	tm.LoadTilesetsFromDir(otherTilesetDir)
	// tm.AddTileset(path.Join(otherTilesetDir, "sprites/Player", "Player_Base_Running.tsx"))
	// tm.AddTileset(path.Join(otherTilesetDir, "sprites/Player", "Pants_Blue_Running.tsx"))
	// tm.AddTileset(path.Join(otherTilesetDir, "sprites/Player", "Shirt_Green_Running.tsx"))
	// tm.AddTileset(path.Join(otherTilesetDir, "sprites/Player", "Medium_Hair_Brown_Running.tsx"))
	// tm.AddTileset(path.Join(otherTilesetDir, "sprites/Player", "Shoes_Brown_Running.tsx"))
	// tm.AddTileset(path.Join(otherTilesetDir, "sprites/Player", "Hands_Bare_Running.tsx"))
	// tm.AddTileset(path.Join(otherTilesetDir, "sprites/Player", "Player_Base_Attack.tsx"))
	// tm.AddTileset(path.Join(otherTilesetDir, "sprites/Player", "Pants_Blue_Attack.tsx"))
	// tm.AddTileset(path.Join(otherTilesetDir, "sprites/Player", "Shirt_Green_Attack.tsx"))
	// tm.AddTileset(path.Join(otherTilesetDir, "sprites/Player", "Medium_Hair_Brown_Attack.tsx"))
	// tm.AddTileset(path.Join(otherTilesetDir, "sprites/Player", "Shoes_Brown_Attack.tsx"))
	// tm.AddTileset(path.Join(otherTilesetDir, "sprites/Player", "Hands_Bare_Attack.tsx"))
	// tm.AddTileset(path.Join(otherTilesetDir, "sprites/Player", "Iron_Sword.tsx"))
	// tm.AddTileset(path.Join(otherTilesetDir, "tilesets/Waterfall", "Waterfall_1.tsx"))

	r = renderer.NewRenderer(tm)
	simpleSprite = renderer.NewSimpleSprite("NinjaDark", r)
	compoundSprite = renderer.NewCompoundSprite([]string{"Player_Base_Running", "Pants_Blue_Running", "Shirt_Green_Running", "Medium_Hair_Brown_Running", "Shoes_Brown_Running", "Hands_Bare_Running"}, r)
	runningAnimation = renderer.NewSimpleAnimation(compoundSprite, 100, []int{})
	attackSprite := renderer.NewCompoundSprite([]string{"Player_Base_Attack", "Pants_Blue_Attack", "Shirt_Green_Attack", "Medium_Hair_Brown_Attack", "Shoes_Brown_Attack", "Hands_Bare_Attack", "Iron_Sword"}, r)
	attackAnimation = renderer.NewSimpleAnimation(attackSprite, 150, []int{0, 1, 2, 4})

	complexSprite = renderer.NewComplexSprite("Waterfall_1", 3, r)
	complexSprite.AddPart(0, []uint32{0, 1, 2, 18, 19, 20, 36, 37, 38, 54, 55, 56, 72, 73, 74})
	complexSprite.AddPart(1, []uint32{3, 4, 5, 21, 22, 23, 39, 40, 41, 57, 58, 59, 75, 76, 77})
	complexSprite.AddPart(2, []uint32{6, 7, 8, 24, 25, 26, 42, 43, 44, 60, 61, 62, 78, 79, 80})
	complexSprite.AddPart(3, []uint32{9, 10, 11, 27, 28, 29, 45, 46, 47, 63, 64, 65, 81, 82, 83})
	complexSprite.AddPart(4, []uint32{12, 13, 14, 30, 31, 32, 48, 49, 50, 66, 67, 68, 84, 85, 86})
	complexSprite.AddPart(5, []uint32{15, 16, 17, 33, 34, 35, 51, 52, 53, 69, 70, 71, 87, 88, 89})

	complexAnimation = renderer.NewSimpleAnimation(complexSprite, 200, []int{0, 1, 2, 3, 4, 5})
	// complexAnimation = renderer.NewSimpleAnimation(complexSprite, 200, []int{0, 2, 3})
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

	// Draw simple sprite by id then Flip Vertical
	simpleSprite.Draw(0, &renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(32, 0)
	simpleSprite.Draw(0, &renderer.DrawOptions{Screen: screen, Op: op, FlipVertical: true})
	nextLine(32)

	// Draw simple sprite by name then Flip Horizontal
	simpleSprite.Draw("WalkRight", &renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(32, 0)
	simpleSprite.Draw("WalkRight", &renderer.DrawOptions{Screen: screen, Op: op, FlipHorizontal: true})
	nextLine(32)

	// Draw simple sprite by name with animation then Flip Horizontal
	simpleSprite.DrawWithAnimation("WalkLeft", 100, &renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(32, 0)
	simpleSprite.DrawWithAnimation("WalkLeft", 100, &renderer.DrawOptions{Screen: screen, Op: op, FlipHorizontal: true})
	nextLine(32)

	// Draw compound sprites
	// Moving down
	for i := 0; i < 6; i++ {
		compoundSprite.Draw(i, &renderer.DrawOptions{Screen: screen, Op: op})
		moveRight(48, 0)

	}
	nextLine(48)
	// Moving right
	for i := 6; i < 12; i++ {
		compoundSprite.Draw(i, &renderer.DrawOptions{Screen: screen, Op: op})
		moveRight(48, 0)
	}
	nextLine(48)
	// Moving left
	for i := 6; i < 12; i++ {
		compoundSprite.Draw(i, &renderer.DrawOptions{Screen: screen, Op: op, FlipHorizontal: true})
		moveRight(48, 0)
	}
	nextLine(48)
	// Moving down
	for i := 12; i < 18; i++ {
		compoundSprite.Draw(i, &renderer.DrawOptions{Screen: screen, Op: op})
		moveRight(48, 0)
	}
	nextLine(48)

	// Draw simple animation
	attackAnimation.SetFrames([]int{0, 1, 2, 3})
	attackAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(48, 0)
	attackAnimation.SetFrames([]int{4, 5, 6, 7})
	attackAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(48, 0)
	attackAnimation.SetFrames([]int{8, 9, 10, 11})
	attackAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(64, 0)
	attackAnimation.SetFrames([]int{12, 13, 14, 15})
	attackAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(48, 0)
	attackAnimation.SetFrames([]int{16, 17, 18, 19})
	attackAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(48, 0)
	attackAnimation.SetFrames([]int{20, 21, 22, 23})
	attackAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(64, 0)
	attackAnimation.SetFrames([]int{12, 13, 14, 15})
	attackAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op, FlipHorizontal: true})
	moveRight(48, 0)
	attackAnimation.SetFrames([]int{16, 17, 18, 19})
	attackAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op, FlipHorizontal: true})
	moveRight(48, 0)
	attackAnimation.SetFrames([]int{20, 21, 22, 23})
	attackAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op, FlipHorizontal: true})
	moveRight(64, 0)
	attackAnimation.SetFrames([]int{24, 25, 26, 27})
	attackAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(48, 0)
	attackAnimation.SetFrames([]int{28, 29, 30, 31})
	attackAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(48, 0)
	attackAnimation.SetFrames([]int{32, 33, 34, 35})
	attackAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op})

	// Running anumations
	nextLine(48)
	runningAnimation.SetFrames([]int{0, 1, 2, 3, 4, 5})
	runningAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(48, 0)
	runningAnimation.SetFrames([]int{6, 7, 8, 9, 10, 11})
	runningAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(48, 0)
	runningAnimation.SetFrames([]int{6, 7, 8, 9, 10, 11})
	runningAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op, FlipHorizontal: true})
	moveRight(48, 0)
	runningAnimation.SetFrames([]int{12, 13, 14, 15, 16, 17})
	runningAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op})

	nextLine(48)
	ts, _ := r.TilesetManager.GetTilesetByName("TilesetFloor")
	r.DrawTile(ts, 0, &renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(16, 0)
	r.DrawTile(ts, 1, &renderer.DrawOptions{Screen: screen, Op: op})
	moveRight(16, 0)
	r.DrawTile(ts, 3, &renderer.DrawOptions{Screen: screen, Op: op})

	nextLine(48)
	// complexSprite.Draw(0, &renderer.DrawOptions{Screen: screen, Op: op})
	//
	// moveRight(64, 0)
	complexAnimation.DrawAnimation(&renderer.DrawOptions{Screen: screen, Op: op})

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
