package animation

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type Part struct {
	Sprite         string `json:"sprite"`
	TileID         int    `json:"tile_id"`
	XOffset        int    `json:"x_offset"`
	YOffset        int    `json:"y_offset"`
	FlipHorizontal bool   `json:"flip_horizontal"`
	FlipVertical   bool   `json:"flip_vertical"`
}

type Frame struct {
	Duration int    `json:"duration"`
	Parts    []Part `json:"parts"`
}

type Animation struct {
	Class   string   `yaml:"class" json:"class"`
	Action  string   `yaml:"action" json:"action"`
	Sprites []string `yaml:"sprites" json:"sprites"`
	Frames  []Frame  `json:"frames"`
	Simple  *simple  `yaml:"simple,omitempty" json:"simple,omitempty"`
	Timed   *timed   `yaml:"timed,omitempty" json:"timed,omitempty"`
	Complex *complex `yaml:"complex,omitempty" json:"complex,omitempty"`

	currentFrame  int
	nextFrameTime int64
}

func (a *Animation) GetCurrentFrame() Frame {
	a.determineFrame()
	return a.Frames[a.currentFrame]
}

func (a *Animation) determineFrame() {
	if a.nextFrameTime == 0 {
		a.currentFrame = 0
		a.setNextFrameTime()
		return
	}
	if time.Now().UnixMilli() >= a.nextFrameTime {
		a.currentFrame++
		if a.currentFrame >= len(a.Frames) {
			a.currentFrame = 0
		}
		a.setNextFrameTime()
	}
}

func (a *Animation) setNextFrameTime() {
	frame := a.Frames[a.currentFrame]
	a.nextFrameTime = time.Now().UnixMilli() + int64(frame.Duration)
}

func (a *Animation) decodeSimple() {
	if a.Simple == nil {
		return
	}
	seq := func(yield func(Frame) bool) {
		for idx := range a.Simple.Frames {
			part := Part{
				TileID:         a.Simple.Frames[idx],
				FlipHorizontal: a.Simple.Defaults.FlipHorizontal,
				FlipVertical:   a.Simple.Defaults.FlipVertical,
			}
			frame := Frame{
				Duration: a.Simple.Duration,
				Parts:    []Part{part},
			}
			if !yield(frame) {
				return
			}
		}
	}
	a.Frames = slices.AppendSeq(a.Frames, seq)
	a.Simple = nil
}

func (a *Animation) decodeTimed() {
	if a.Timed == nil {
		return
	}
	seq := func(yield func(Frame) bool) {
		for _, frame := range a.Timed.Frames {
			part := Part{
				TileID:         frame.ID,
				FlipHorizontal: a.Timed.Defaults.FlipHorizontal,
				FlipVertical:   a.Timed.Defaults.FlipVertical,
			}
			f := Frame{
				Duration: frame.Duration,
				Parts:    []Part{part},
			}
			if !yield(f) {
				return
			}
		}
	}
	a.Frames = slices.AppendSeq(a.Frames, seq)
	a.Timed = nil
}

func (a *Animation) decodeComplex() {
	if a.Complex == nil {
		return
	}

	frameSeq := func(yield func(Frame) bool) {
		for _, frame := range a.Complex.Frames {

			partSeq := func(yield func(Part) bool) {
				for _, part := range frame.Parts {
					p := Part{
						TileID:         part.ID,
						Sprite:         a.Sprites[part.Sprite-1],
						XOffset:        part.XOffset,
						YOffset:        part.YOffset,
						FlipHorizontal: part.FlipHorizontal,
						FlipVertical:   part.FlipVertical,
					}
					if !yield(p) {
						return
					}
				}
			}

			parts := slices.AppendSeq([]Part{}, partSeq)

			f := Frame{
				Duration: frame.Duration,
				Parts:    parts,
			}
			if !yield(f) {
				return
			}
		}
	}

	a.Frames = slices.AppendSeq(a.Frames, frameSeq)
	a.Complex = nil
}

type Animations struct {
	baseDir    string
	Source     string       `json:"source"`
	Animations []*Animation `yaml:"animations" json:"animations"`
}

// LoadReader function loads tileset in TSX format from io.Reader
// baseDir is used for loading additional tile data, current directory is used if empty
func fileReader(source string, r io.Reader) (*Animations, error) {
	d := yaml.NewDecoder(r)

	baseDir := filepath.Dir(source)
	animations := &Animations{
		baseDir: baseDir,
		Source:  source,
	}

	if err := d.Decode(animations); err != nil {
		return nil, err
	}

	for _, animation := range animations.Animations {
		animation.Frames = []Frame{}

		animation.decodeSimple()
		animation.decodeTimed()
		animation.decodeComplex()
	}

	return animations, nil
}

// LoadFile function loads tileset in TSX format from file
func LoadFile(fileName string) (*Animations, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return fileReader(fileName, f)
}

type defaults struct {
	FlipHorizontal bool `yaml:"flip_horizontal"`
	FlipVertical   bool `yaml:"flip_vertical"`
}
type simple struct {
	Duration int      `yaml:"duration"`
	Defaults defaults `yaml:"defaults"`
	Frames   []int    `yaml:"frames"`
}
type timed struct {
	Defaults defaults `yaml:"defaults"`
	Frames   []struct {
		ID       int `yaml:"id"`
		Duration int `yaml:"duration"`
	} `yaml:"frames"`
}
type complex struct {
	Frames []struct {
		Duration int `yaml:"duration"`
		Parts    []struct {
			ID             int  `yaml:"id"`
			Sprite         int  `yaml:"sprite"`
			XOffset        int  `yaml:"x_offset"`
			YOffset        int  `yaml:"y_offset"`
			FlipHorizontal bool `yaml:"flip_horizontal"`
			FlipVertical   bool `yaml:"flip_vertical"`
		} `yaml:"parts"`
	} `yaml:"frames"`
}
