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

func Jsontostruct[T any](path string) T { // convertir un fichier json en struct
	t := *new(T)                   // créer une nouvelle instance de T
	data, err := os.ReadFile(path) // lire le fichier
	if err != nil {                // si une erreur est survenue
		return *new(T)
	}
	err = json.Unmarshal(data, &t) // déserialiser le json
	if err != nil {                // si une erreur est survenue
		return *new(T)
	}
	return t
}

func GetAllFilesInDirectory(path string) []string { // obtenir tous les fichiers dans un dossier
	files, err := os.ReadDir(path) // lire le dossier
	if err != nil {                // si une erreur est survenue
		return nil
	}
	result := make([]string, len(files)) // créer un slice de string
	for i, file := range files {         // pour chaque fichier
		result[i] = file.Name() // ajouter le nom du fichier au slice
	}
	return result
}

// ne marche pas avec les fichiers qui ne sont pas des json
func GetAllFilesInDirectoryToStruct[T any](path string) []struct {
	Name string
	Obj  T
} {
	files := GetAllFilesInDirectory(path) // obtenir tous les fichiers dans le dossier
	result := make([]struct {
		Name string
		Obj  T
	}, len(files)) // créer un slice de struct
	for i, file := range files { // pour chaque fichier
		result[i] = struct {
			Name string
			Obj  T
		}{file, Jsontostruct[T](path + "/" + file)}
	}
	return result
}

func TiledMapToImage(Map *tiled.Map) image.Image { // convertir une carte tiled en image
	renderer, err := render.NewRenderer(Map) // créer un renderer
	if err != nil {                          // si une erreur est survenue
		fmt.Println("map unsupported for rendering: " + err.Error())
		os.Exit(2)
	}

	// Render just layer 0 to the Renderer.
	err = renderer.RenderVisibleLayers()
	if err != nil { // si une erreur est survenue
		fmt.Println("layer unsupported for rendering: " + err.Error())
		os.Exit(2)
	}

	return renderer.Result
}

func GetAllFilesInDirectoryToMap(path string) map[string]tiled.Map { // obtenir toutes les cartes tiled dans un dossier
	result := make(map[string]tiled.Map) // créer une map de string vers tiled.Map

	for _, file := range GetAllFilesInDirectory(path) { // pour chaque fichier
		gameMap, err := tiled.LoadFile(path + "/" + file)
		if err != nil {
			fmt.Println("error parsing map: " + err.Error())
			os.Exit(2)
		}
		result[file] = *gameMap
	}
	return result
}

func GetAllImagesInDirectory(path string) map[string]image.Image { // obtenir toutes les images dans un dossier
	result := make(map[string]image.Image) // créer une map de string vers image.Image

	for _, file := range GetAllFilesInDirectory(path) { // pour chaque fichier
		imgFile, err := os.Open(path + "/" + file)
		if err != nil { // si une erreur est survenue
			fmt.Println("error opening image: " + err.Error())
			os.Exit(2)
		}
		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil { // si une erreur est survenue
			fmt.Println("error decoding image: " + err.Error())
			os.Exit(2)
		}
		result[file] = img
	}
	return result
}

func CheckCollision(box [4]float64, elements []*ETEStruct.Sprite, maxDist float64, Unité int) bool { // vérifier si une boîte collisionne avec des éléments (maxDist en unités)
	maxDist *= float64(Unité)
	for _, element := range elements { // pour chaque élément
		Box := ElementToBox(element)                                                            // obtenir la boîte de l'élément
		if math.V2Distance([2]float64{box[2], box[3]}, [2]float64{Box[2], Box[3]}) <= maxDist { // si la distance entre les centres est inférieure à la distance maximale (en unités)
			if CheckIntersection([2][4]float64{box, Box}, Unité) { // vérifier si les boîtes se chevauchent
				return true
			}
		}
	}

	return false
}

func CheckIntersection(boxs [2][4]float64, Unité int) bool { // {witdh height  PosX, PosY}
	unit := float64(Unité)
	if boxs[0][1] == 0 { // si la hauteur est 0, c'est un cercle
		if boxs[1][1] == 0 { // si la hauteur est 0, c'est un cercle
			// Cercle vs Cercle

			Obj0 := resolv.NewCircle(boxs[0][2]*unit, boxs[0][3]*unit, boxs[0][0]*unit)
			Obj1 := resolv.NewCircle(boxs[1][2]*unit, boxs[1][3]*unit, boxs[1][0]*unit)

			if !Obj0.Intersection(Obj1).IsEmpty() { // si les cercles se chevauchent
				return true
			}
		} else { // sinon c'est un rectangle
			// Cercle vs Rectangle

			Obj0 := resolv.NewCircle(boxs[0][2]*unit, boxs[0][3]*unit, boxs[0][0]*unit)
			Obj1 := resolv.NewRectangle(boxs[1][2]*unit, boxs[1][3]*unit, boxs[1][0]*unit, boxs[1][1]*unit)

			if !Obj0.Intersection(Obj1).IsEmpty() { // si le cercle et le rectangle se chevauchent
				return true
			}
		}
	} else { // sinon c'est un rectangle
		if boxs[1][1] == 0 { // si la hauteur est 0, c'est un cercle
			//Rectangle vs Cercle

			Obj0 := resolv.NewRectangle(boxs[0][2]*unit, boxs[0][3]*unit, boxs[0][0]*unit, boxs[0][1]*unit)
			Obj1 := resolv.NewCircle(boxs[1][2]*unit, boxs[1][3]*unit, boxs[1][0]*unit)

			if !Obj0.Intersection(Obj1).IsEmpty() { // si le rectangle et le cercle se chevauchent
				return true
			}
		} else { // sinon c'est un rectangle
			// Rectangle vs Rectangle

			Obj0 := resolv.NewRectangle(boxs[0][2]*unit, boxs[0][3]*unit, boxs[0][0]*unit, boxs[0][1]*unit)
			Obj1 := resolv.NewRectangle(boxs[1][2]*unit, boxs[1][3]*unit, boxs[1][0]*unit, boxs[1][1]*unit)

			if !Obj0.Intersection(Obj1).IsEmpty() { // si les rectangles se chevauchent
				return true
			}
		}
	}
	return false
}

func ElementToBox(element *ETEStruct.Sprite) [4]float64 { // convertir un sprite en boîte
	return [4]float64{element.Box[0], element.Box[1], element.Pos[0] + element.Box[2], element.Pos[1] + element.Box[3]}
}
