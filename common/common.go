package common

import (
	"errors"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ErrPathNotFound     = errors.New("path not found")
	ErrPathNotDirectory = errors.New("path is not a directory")
)

type DrawOptions struct {
	Screen         *ebiten.Image
	Op             *ebiten.DrawImageOptions
	FlipHorizontal bool
	FlipVertical   bool
	OffsetX        float64
	OffsetY        float64
}

func PathShouldBeDirectory(path string) error {
	info, err := os.Stat(path)
	// Check if the path exists
	if os.IsNotExist(err) {
		return fmt.Errorf("path: %s %w", path, ErrPathNotFound)
	}
	// Check if the path is a directory
	if !info.IsDir() {
		return fmt.Errorf("path: %s %w", path, ErrPathNotDirectory)
	}
	return nil
}
