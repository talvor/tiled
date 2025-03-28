# tmx

A go module to read TMX maos from tiled map editor (https://www.mapeditor.org/)

## Reading individual Maps

To read in a map use the `tmx.LoadFile` function

```golang
package main

import (
    "encoding/json"
    "fmt"

    "github.com/talvor/tmx"
)

func main() {
  map1, err := tmx.LoadFile("~/Documents/maps/scene1.tmx")
  if err != nil {
    panic(err)
  }

  mapJSON, _ := json.Marshal(map1)
  fmt.Println(string(mapJSON))
}
```

## Managing multiple maps using MapManager

To read in bulk maps use the `MapManager` struct.

```golang
package main

import "github.com/talvor/tmx"

func main() {
	// Create a new map manager and load all maps from the directory
	mm := tmx.NewMapManager("~/Documents/maps")

	map, err := mm.GetMapByName("scene1")
  if err != nil {
    panic(err)
  }
}
```

## Using the renderer

The `tmx.renderer` works with the `MapManager` and `TsxRenderer` to provide convenient methods for rendering
maps into the ebiten screen.

See `renderer/examples/main.go` for an example of using the renderer
