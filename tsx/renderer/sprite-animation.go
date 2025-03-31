package renderer

import (
	"time"
)

type SpriteAnimator interface {
	DrawAnimation(opts *DrawOptions) error
}

type SimpleAnimation struct {
	sprite   SpriteDrawer
	frames   []int
	duration uint32
	defaults *AnimationDefaults
}

type AnimationDefaults struct {
	FlipHorizontal bool
	FlipVertical   bool
}

func NewSimpleAnimation(sprite SpriteDrawer, duration uint32, frames []int, defaults *AnimationDefaults) *SimpleAnimation {
	return &SimpleAnimation{
		sprite:   sprite,
		frames:   frames,
		duration: duration,
		defaults: defaults,
	}
}

func (sa *SimpleAnimation) SetFrames(frames []int) {
	sa.frames = frames
}

func (sa *SimpleAnimation) DrawAnimation(opts *DrawOptions) error {
	if sa.defaults != nil {
		opts.FlipHorizontal = sa.defaults.FlipHorizontal
		opts.FlipVertical = sa.defaults.FlipVertical
	}

	animationIdx := int(time.Now().UnixMilli()) / int(sa.duration) % len(sa.frames)
	frame := sa.frames[animationIdx]
	if frame == -1 {
		return nil
	}
	return sa.sprite.Draw(uint32(frame), opts)
}
