package renderer

import (
	"errors"
	"time"
)

var ErrFrameTimingMismatch = errors.New("frame and timing slices must be the same length")

type SpriteAnimator interface {
	DrawAnimation(opts *DrawOptions) error
}

type AnimationDefaults struct {
	FlipHorizontal bool
	FlipVertical   bool
}

type Frame struct {
	id       int
	duration uint32
}

type animation struct {
	sprite        SpriteDrawer
	frames        []Frame
	defaults      *AnimationDefaults
	currentFrame  int
	nextFrameTime int64
}

func (a *animation) DrawAnimation(opts *DrawOptions) error {
	a.applyDefaults(opts)
	a.determineFrame()

	frame := a.frames[a.currentFrame]
	if frame.id == -1 {
		return nil
	}
	return a.sprite.Draw(uint32(frame.id), opts)
}

func (a *animation) applyDefaults(opts *DrawOptions) {
	if a.defaults != nil {
		opts.FlipHorizontal = a.defaults.FlipHorizontal
		opts.FlipVertical = a.defaults.FlipVertical
	}
}

func (a *animation) determineFrame() {
	if a.nextFrameTime == 0 {
		a.currentFrame = 0
		a.setNextFrameTime()
		return
	}
	if time.Now().UnixMilli() >= a.nextFrameTime {
		a.currentFrame++
		if a.currentFrame >= len(a.frames) {
			a.currentFrame = 0
		}
		a.setNextFrameTime()
	}
}

func (a *animation) setNextFrameTime() {
	frame := a.frames[a.currentFrame]
	a.nextFrameTime = time.Now().UnixMilli() + int64(frame.duration)
}

type SimpleAnimation struct {
	animation
	duration uint32
}

func NewSimpleAnimation(sprite SpriteDrawer, duration uint32, frames []int, defaults *AnimationDefaults) *SimpleAnimation {
	timedFrames := makeTimedFrames(frames, duration)
	return &SimpleAnimation{
		animation: animation{
			sprite:       sprite,
			frames:       timedFrames,
			defaults:     defaults,
			currentFrame: 0,
		},
		duration: duration,
	}
}

func (sa *SimpleAnimation) SetFrames(frames []int) {
	sa.frames = makeTimedFrames(frames, sa.duration)
}

type TimedAnimation struct {
	animation
}

func NewTimedAnimation(sprite SpriteDrawer, frames []int, timing []uint32, defaults *AnimationDefaults) (*TimedAnimation, error) {
	if len(frames) != len(timing) {
		return nil, ErrFrameTimingMismatch
	}

	timedFrames := make([]Frame, len(frames))
	for i, f := range frames {
		timedFrames[i] = Frame{
			id:       f,
			duration: timing[i],
		}
	}

	return &TimedAnimation{
		animation: animation{
			sprite:       sprite,
			frames:       timedFrames,
			defaults:     defaults,
			currentFrame: 0,
		},
	}, nil
}

func makeTimedFrames(frames []int, duration uint32) []Frame {
	timedFrames := make([]Frame, len(frames))
	for i, f := range frames {
		timedFrames[i] = Frame{
			id:       f,
			duration: duration,
		}
	}
	return timedFrames
}
