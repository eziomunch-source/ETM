package ETEhelper

import (
	"encoding/json"
	"fmt"
	"image"
	"os"

	ETEStruct "github.com/Try-si/ETM/ETEStruct"
	Math "github.com/Try-si/MathHelper/Math"
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

func IsInCollision(box [6]float64, Elements []*ETEStruct.Sprite, CellSize float64) bool { // box[width, height, x, y]
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

	if box[1] == 0 {
		for _, element := range World {
			if element.Box[1] == 0 {
				if Math.V2Distance([2]float64{box[2], box[3]}, [2]float64{element.Pos[0], element.Pos[1]}) < element.Box[0]+box[0] {
					return true
				}
			} else {
				if Math.V2Distance([2]float64{box[2], box[3]}, [2]float64{element.Pos[0], element.Pos[1]}) < Math.V2Length([2]float64{element.Box[0], element.Box[1]})+box[0] {
					return true
				}
			}
		}
	} else {
		for _, element := range World {
			if element.Box[1] == 0 {
				// 1. Trouver le point le plus proche sur le rectangle (clamp)
				closestX := box[2]
				if element.Pos[0] < box[2] {
					closestX = box[2] // cercle à gauche
				} else if element.Pos[0] > box[2]+box[0] {
					closestX = box[2] + box[0] // cercle à droite
				} else {
					closestX = element.Pos[0] // cercle dans le rectangle (en X)
				}

				closestY := box[3]
				if element.Pos[1] < box[3] {
					closestY = box[3] // cercle en haut
				} else if element.Pos[1] > box[3]+box[1] {
					closestY = box[3] + box[1] // cercle en bas
				} else {
					closestY = element.Pos[1] // cercle dans le rectangle (en Y)
				}

				// 2. Calculer la distance au carré (évite sqrt)
				distanceX := element.Pos[0] - closestX
				distanceY := element.Pos[1] - closestY
				distanceSquared := distanceX*distanceX + distanceY*distanceY

				// 3. Comparer au rayon au carré
				if distanceSquared < element.Box[0]*element.Box[0] {
					return true
				}
			} else {
				if box[2] < element.Pos[0]+element.Box[0] && // bord gauche box < bord droit element
					box[2]+box[0] > element.Pos[0] && // bord droit box > bord gauche element
					box[3] < element.Pos[1]+element.Box[1] && // bord haut box < bord bas element
					box[3]+box[1] > element.Pos[1] { // bord bas box > bord haut element
					return true
				}
			}
		}
	}

	return false
}
