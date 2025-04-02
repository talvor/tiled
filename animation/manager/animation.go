package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/talvor/tiled/animation"
	tsxm "github.com/talvor/tiled/tsx/manager"
	renderer "github.com/talvor/tiled/tsx/renderer"
	tsxr "github.com/talvor/tiled/tsx/renderer"
)

type AnimationManager struct {
	baseDir    string
	Animations map[string]*animation.Animation
	TSManager  *tsxm.TilesetManager
	TSRenderer *tsxr.Renderer
}

func (am *AnimationManager) SetTilesetManager(tsManager *tsxm.TilesetManager) {
	am.TSManager = tsManager
}

func (am *AnimationManager) SetRenderer(renderer *renderer.Renderer) {
	am.TSRenderer = renderer
}

func (am *AnimationManager) DebugPrintAnimations() {
	for name := range am.Animations {
		fmt.Println(name)
	}
}

func NewManager(baseDir string) (*AnimationManager, error) {
	am := &AnimationManager{
		baseDir:    baseDir,
		Animations: make(map[string]*animation.Animation),
	}

	if err := loadAnimations(am, baseDir); err != nil {
		return nil, err
	}

	return am, nil
}

func loadAnimations(am *AnimationManager, baseDir string) error {
	aniFiles, err := findANIFiles(baseDir)
	if err != nil {
		return fmt.Errorf("error loading animations: %s %w", baseDir, err)
	}
	for _, aniFile := range aniFiles {
		animations, err := animation.LoadFile(aniFile)
		if err != nil {
			return fmt.Errorf("error loading animations: %s %w", baseDir, err)
		}

		for _, animation := range animations.Animations {
			am.Animations[animation.Class] = animation
		}
	}

	return nil
}

func findANIFiles(dir string) ([]string, error) {
	var aniFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".ani" {
			aniFiles = append(aniFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return aniFiles, nil
}
