# animation

A go module to read ANI animation files

Animation files work in conjunction with the animation module and helps define and render sprite animations

## Reading individual animation files

To read in an animation file use the `animation.LoadFile` function

```golang
package main

import (
    "encoding/json"
    "fmt"

    "github.com/talvor/tiled/animation"
)

func main() {
  ani, err := animation.LoadFile("~/Documents/animations/player.ani")
  if err != nil {
    panic(err)
  }

  aniJSON, _ := json.Marshal(ani)
  fmt.Println(string(aniJSON))
}
```

## Managing multiple animations using AnimationManager

To read in bulk animation files use the `AnimationManager` struct.

```golang
package main

import "github.com/talvor/tiled/animation"

func main() {
	// Create a new animation manager and load all animations from the directory
	aman:= animation.NewManager("~/Documents/animations")

	animation, _ := aman.GetAnimation(className, actionName)
}
```

## Using the renderer

The `animation.renderer` works with the `TilesetRenderer` and the [ebitenengine](https://ebitengine.org/) 2D game engine to provide convenient methods for rendering
animations into the ebiten screen.

See `cmd/animation/renderer/main.go` for an example of using the renderer

## Animation File Format

The _animation_ file format use YAML to define a set of animations.

The yaml file must have a top level fields

- tileset_groups: A list of tileset groups
- animations: This is a list of animations

### tileset_group

Each tileset group has the following properties

- name: Name of the tileset group. Used to reference the tileset group inside an animation, so must be unique.
- tilesets: An ordered list of tilesets required to render the animation.

### animation

Each animation has the following properties

- class: A key used to group a number of animations together. Eg. all animations for the "player" sprite
- action: A secondary key used to signify what the animation is for. The combination of `class` and `action` should be unique

And one of, "tileset_group" and/or "tilesets". If both are provided, the "tilesets" list is appened to the end of the
tileset list referenced by the tileset_group

- tileset_group: A group of tilesets required to render the animation.
- tilesets: An ordered list of tilesets required to render the animation.

And then one of the animation types "simple", "timed" or "complex"

#### simple animation

- simple: Defines an animation where each frame renders the same tile from each tileset. The duration is set for the animation.

  - duration: Defines the total duration for all frames. ie each frames duration = duration / len(frames)
  - defaults:
    - flip_horizontally: Will cause each frame to be rendered flipped horizontally
    - flip_vertically: Will cause each frame to be rendered flipped vertically
    - x_offset: Will cause each frame to be rendered offset on the x axis
    - y_offset: Will cause each frame to be rendered offset on the y axis
  - frames: Ordered list of tile id's

#### complex animation

- timed: Defines an animation where each frame renders the same tile from each tileset. The duration is set for each frame.

  - defaults:
    - flip_horizontally: Will cause each frame to be rendered flipped horizontally
    - flip_vertically: Will cause each frame to be rendered flipped vertically
    - x_offset: Will cause each frame to be rendered offset on the x axis
    - y_offset: Will cause each frame to be rendered offset on the y axis
  - frames: Ordered list of frames
    - id: tile id
    - duration: duration this frame will be rendered

#### complex animation

- complex: Defines an animation where each frame is made up of parts.
  - frames: Ordered list of frames
    - duration: duration this frame will be rendered
    - parts: Ordered list of parts that will be rendered for this frame
      - id: tile id
      - tileset: The index in the "tilesets" list for the tileset to render from
      - flip_horizontally: Will cause each frame to be rendered flipped horizontally
      - flip_vertically: Will cause each frame to be rendered flipped vertically
      - x_offset: Will cause each frame to be rendered offset on the x axis
      - y_offset: Will cause each frame to be rendered offset on the y axis

Sample animation files are available in `/animation/example`
