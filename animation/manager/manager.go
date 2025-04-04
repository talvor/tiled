package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/talvor/tiled/animation"
	"github.com/talvor/tiled/common"
)

var (
	ErrAnimationNotFound    = errors.New("manager: animation not found")
	ErrTilesetGroupNotFound = errors.New("manager: tileset group not found")
)

type AnimationManager struct {
	baseDir       string
	Animations    map[string]*animation.Animation
	TilesetGroups map[string]*animation.TilesetGroup
}

func (am *AnimationManager) GetAnimation(class string, action string) (*animation.Animation, error) {
	name := makeAnimationName(class, action)
	ani, ok := am.Animations[name]
	if !ok {
		return nil, fmt.Errorf("class:%s action:%s %w", class, action, ErrAnimationNotFound)
	}

	return ani, nil
}

func (am *AnimationManager) GetTilesetGroup(name string) (*animation.TilesetGroup, error) {
	group, ok := am.TilesetGroups[name]
	if !ok {
		return nil, fmt.Errorf("group: %s %w", name, ErrTilesetGroupNotFound)
	}
	return group, nil
}

func (am *AnimationManager) DebugPrintAnimations() {
	for name := range am.Animations {
		fmt.Println(name)
	}
}

func NewManager(baseDirs []string) *AnimationManager {
	am := &AnimationManager{
		Animations:    make(map[string]*animation.Animation),
		TilesetGroups: make(map[string]*animation.TilesetGroup),
	}

	// Load animations from the base directories
	for _, baseDir := range baseDirs {
		if err := common.PathShouldBeDirectory(baseDir); err != nil {
			fmt.Printf("Error: %s %v\n", baseDir, err)
			continue
		}

		if err := loadAnimations(am, baseDir); err != nil {
			fmt.Printf("Error loading animations from %s: %v\n", baseDir, err)
		}
	}

	return am
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

		for _, group := range animations.TilesetGroups {
			am.TilesetGroups[group.Name] = group
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
