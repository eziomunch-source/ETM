package ETEhelper

import (
	"encoding/json"
	"fmt"
	"image"
	"os"

	ETEStruct "github.com/Try-si/ETM/ETEStruct"
	math "github.com/Try-si/MathHelper/Math"
	tiled "github.com/lafriks/go-tiled"
	render "github.com/lafriks/go-tiled/render"
	"github.com/solarlune/resolv"
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

func CheckCollision(box [6]float64, elements []*ETEStruct.Sprite, maxDist int) bool {
	for _, element := range elements {
		Box := ElementToBox(element)
		if math.V2Distance([2]float64{box[4], box[5]}, [2]float64{Box[4], Box[5]}) <= float64(maxDist) {
			if CheckIntersection([2][6]float64{box, Box}) {
				return true
			}
		}
	}

	return false
}

func CheckIntersection(boxs [2][6]float64) bool { // {witdh height  offsetX offsetY  PosX, PosY}
	if boxs[0][1] == 0 {
		if boxs[1][1] == 0 {
			// Cercle vs Cercle

			Obj0 := resolv.NewCircle(boxs[0][4], boxs[0][5], boxs[0][0])
			Obj1 := resolv.NewCircle(boxs[1][4], boxs[1][5], boxs[1][0])

			if !Obj0.Intersection(Obj1).IsEmpty() {
				return true
			}
		} else {
			// Cercle vs Rectangle

			Obj0 := resolv.NewCircle(boxs[0][4], boxs[0][5], boxs[0][0])
			Obj1 := resolv.NewRectangle(boxs[1][4], boxs[1][5], boxs[1][0], boxs[1][1])

			if !Obj0.Intersection(Obj1).IsEmpty() {
				return true
			}
		}
	} else {
		if boxs[1][1] == 0 {
			//Rectangle vs Cercle

			Obj0 := resolv.NewRectangle(boxs[0][4], boxs[0][5], boxs[0][0], boxs[0][1])
			Obj1 := resolv.NewCircle(boxs[1][4], boxs[1][5], boxs[1][0])

			if !Obj0.Intersection(Obj1).IsEmpty() {
				return true
			}
		} else {
			// Rectangle vs Rectangle

			Obj0 := resolv.NewRectangle(boxs[0][4], boxs[0][5], boxs[0][0], boxs[0][1])
			Obj1 := resolv.NewRectangle(boxs[1][4], boxs[1][5], boxs[1][0], boxs[1][1])

			if !Obj0.Intersection(Obj1).IsEmpty() {
				return true
			}
		}
	}
	return false
}

func ElementToBox(element *ETEStruct.Sprite) [6]float64 {
	return [6]float64{element.Box[0], element.Box[1], element.Box[2], element.Box[3], element.Pos[0] + element.Box[2], element.Pos[1] + element.Box[3]}
}
