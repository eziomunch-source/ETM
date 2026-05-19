package ETEhelper

import (
	"encoding/json"
	"fmt"
	"image"
	"os"

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
