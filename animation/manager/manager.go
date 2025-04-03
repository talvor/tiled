package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/talvor/tiled/animation"
)

var ErrAnimationNotFound = errors.New("animation not found")

type AnimationManager struct {
	baseDir    string
	Animations map[string]*animation.Animation
}

func (am *AnimationManager) GetAnimation(class string, action string) (*animation.Animation, error) {
	name := makeAnimationName(class, action)
	ani, ok := am.Animations[name]
	if !ok {
		return nil, fmt.Errorf("class:%s action:%s %w", class, action, ErrAnimationNotFound)
	}

	return ani, nil
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
			name := makeAnimationName(animation.Class, animation.Action)
			am.Animations[name] = animation
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

func makeAnimationName(class string, action string) string {
	return fmt.Sprintf("%s_%s", class, action)
}
