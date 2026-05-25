package assets

import (
	ETMhelper "github.com/Try-si/ETM/Helper"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Assets struct {
	Images     map[string]*ebiten.Image
	Maps       map[string]*ebiten.Image
	SpritePath string
	MapsPath   string
}

func (a *Assets) Init() {
	a.Images = make(map[string]*ebiten.Image)
	a.Maps = make(map[string]*ebiten.Image)

	for fileName, img := range ETMhelper.GetAllImagesInDirectory(a.SpritePath) {
		a.Images[fileName] = ebiten.NewImageFromImage(img)
	}

	for fileName, mapFile := range ETMhelper.GetAllFilesInDirectoryToMap(a.MapsPath) {
		a.Maps[fileName] = ebiten.NewImageFromImage(ETMhelper.TiledMapToImage(&mapFile))
	}

}

func (a *Assets) GetImage(name string) (*ebiten.Image, bool) {
	img, exists := a.Images[name]
	return img, exists
}

func (a *Assets) GetMap(name string) (*ebiten.Image, bool) {
	img, exists := a.Maps[name]
	return img, exists
}
