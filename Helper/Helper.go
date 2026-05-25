package ETEhelper

import (
	"encoding/json"
	"fmt"
	"image"
	"math"
	"os"

	ETEStruct "github.com/Try-si/ETM/ETEStruct"
	tiled "github.com/lafriks/go-tiled"
	render "github.com/lafriks/go-tiled/render"
)

func Jsontostruct[T any](path string) T {
	t := *new(T)
	data, err := os.ReadFile(path)
	if err != nil {
		return *new(T)
	}
	err = json.Unmarshal(data, &t)
	if err != nil {
		return *new(T)
	}
	return t
}

func GetAllFilesInDirectory(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	result := make([]string, len(files))
	for i, file := range files {
		result[i] = file.Name()
	}
	return result
}

// ne marche pas avec les fichiers qui ne sont pas des json
func GetAllFilesInDirectoryToStruct[T any](path string) []struct {
	Name string
	Obj  T
} {
	files := GetAllFilesInDirectory(path)
	result := make([]struct {
		Name string
		Obj  T
	}, len(files))
	for i, file := range files {
		result[i] = struct {
			Name string
			Obj  T
		}{file, Jsontostruct[T](path + "/" + file)}
	}
	return result
}

func TiledMapToImage(Map *tiled.Map) image.Image {
	// You can also render the map to an in-memory image for direct
	// use with the default Renderer, or by making your own.
	renderer, err := render.NewRenderer(Map)
	if err != nil {
		fmt.Println("map unsupported for rendering: " + err.Error())
		os.Exit(2)
	}

	// Render just layer 0 to the Renderer.
	err = renderer.RenderVisibleLayers()
	if err != nil {
		fmt.Println("layer unsupported for rendering: " + err.Error())
		os.Exit(2)
	}

	return renderer.Result
}

func GetAllFilesInDirectoryToMap(path string) map[string]tiled.Map {
	result := make(map[string]tiled.Map)

	for _, file := range GetAllFilesInDirectory(path) {
		gameMap, err := tiled.LoadFile(path + "/" + file)
		if err != nil {
			fmt.Println("error parsing map: " + err.Error())
			os.Exit(2)
		}
		result[file] = *gameMap
	}
	return result
}

func GetAllImagesInDirectory(path string) map[string]image.Image {
	result := make(map[string]image.Image)

	for _, file := range GetAllFilesInDirectory(path) {
		imgFile, err := os.Open(path + "/" + file)
		if err != nil {
			fmt.Println("error opening image: " + err.Error())
			os.Exit(2)
		}
		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil {
			fmt.Println("error decoding image: " + err.Error())
			os.Exit(2)
		}
		result[file] = img
	}
	return result
}

func IsInCollision(box [6]float64, Elements []*ETEStruct.Sprite, CellSize float64) bool { // box[width, height, offsetX, offsetY, x, y]
	HashWorld := ETEStruct.GetElementByHashMap(Elements, CellSize)

	key := [2]int{
		int(Elements[0].Pos[0] / CellSize),
		int(Elements[0].Pos[1] / CellSize),
	}

	keyUp := string(key[0]) + "_" + string(key[1]+int(CellSize))
	keyDown := string(key[0]) + "_" + string(key[1]-int(CellSize))
	keyLeft := string(key[0]-int(CellSize)) + "_" + string(key[1])
	keyRight := string(key[0]+int(CellSize)) + "_" + string(key[1])
	keyMiddle := string(key[0]) + "_" + string(key[1])

	World := make([]*ETEStruct.Sprite, 0)

	for h, elements := range HashWorld {
		if h == keyUp {
			World = append(World, elements...)
			// TODO: check collision
		}
		if h == keyDown {
			World = append(World, elements...)
			// TODO: check collision
		}
		if h == keyLeft {
			World = append(World, elements...)
			// TODO: check collision
		}
		if h == keyRight {
			World = append(World, elements...)
			// TODO: check collision
		}
		if h == keyMiddle {
			World = append(World, elements...)
			// TODO: check collision
		}
	}

	// Vérifié les collisions
	for _, element := range World {
		elementX := element.Pos[0] + element.Box[2]
		elementY := element.Pos[1] + element.Box[3]

		if box[1] == 0 {
			// Box is a circle
			if element.Box[1] == 0 {
				// Circle vs Circle
				dx := box[4] - elementX
				dy := box[5] - elementY
				distance := math.Sqrt(dx*dx + dy*dy)
				if distance < box[0]+element.Box[0] {
					return true
				}
			} else {
				// Circle vs Rectangle
				closestX := box[4]
				if elementX < box[4] {
					closestX = elementX
				} else if elementX > box[4]+element.Box[0] {
					closestX = box[4] + element.Box[0]
				}

				closestY := box[5]
				if elementY < box[5] {
					closestY = elementY
				} else if elementY > box[5]+element.Box[1] {
					closestY = box[5] + element.Box[1]
				}

				dx := box[4] - closestX
				dy := box[5] - closestY
				distanceSquared := dx*dx + dy*dy

				if distanceSquared < box[0]*box[0] {
					return true
				}
			}
		} else {
			// Box is a rectangle
			if element.Box[1] == 0 {
				// Rectangle vs Circle
				closestX := elementX
				if box[4] < elementX {
					closestX = box[4]
				} else if box[4]+box[0] > elementX {
					closestX = box[4] + box[0]
				}

				closestY := elementY
				if box[5] < elementY {
					closestY = box[5]
				} else if box[5]+box[1] > elementY {
					closestY = box[5] + box[1]
				}

				dx := elementX - closestX
				dy := elementY - closestY
				distanceSquared := dx*dx + dy*dy

				if distanceSquared < element.Box[0]*element.Box[0] {
					return true
				}
			} else {
				// Rectangle vs Rectangle
				if box[4] < elementX+element.Box[0] &&
					box[4]+box[0] > elementX &&
					box[5] < elementY+element.Box[1] &&
					box[5]+box[1] > elementY {
					return true
				}
			}
		}
	}

	return false
}
