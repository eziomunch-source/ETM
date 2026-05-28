package assets

import (
	ETMhelper "github.com/Try-si/ETM/Helper"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Assets struct {
	Images     map[string]*ebiten.Image // toutes les images
	Maps       map[string]*ebiten.Image // toutes les maps
	SpritePath string                   // chemin vers tous les fichiers sprites
	MapsPath   string                   // chemin vers tous les fichiers maps
}

func (a *Assets) Init() {
	a.Images = make(map[string]*ebiten.Image) // initialize images map
	a.Maps = make(map[string]*ebiten.Image)   // initialize maps map

	for fileName, img := range ETMhelper.GetAllImagesInDirectory(a.SpritePath) { // charger toutes les images
		a.Images[fileName] = ebiten.NewImageFromImage(img) // convertir l'image en image ebiten
	}

	for fileName, mapFile := range ETMhelper.GetAllFilesInDirectoryToMap(a.MapsPath) { // charger toutes les maps
		a.Maps[fileName] = ebiten.NewImageFromImage(ETMhelper.TiledMapToImage(&mapFile)) // convertir la map en image ebiten
	}
}

func (a *Assets) GetImage(name string) (*ebiten.Image, bool) {
	img, exists := a.Images[name] // obtenir l'image par son nom
	return img, exists
}

func (a *Assets) GetMap(name string) (*ebiten.Image, bool) {
	img, exists := a.Maps[name] // obtenir la map par son nom
	return img, exists
}
