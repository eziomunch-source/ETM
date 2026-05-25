package assets

import (
	ETMhelper "github.com/Try-si/ETM/Helper"
	"github.com/hajimehoshi/ebiten"
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
		var err error
		a.Images[fileName], err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

		if err != nil {
			panic(err)
		}
	}

	for fileName, mapFile := range ETMhelper.GetAllFilesInDirectoryToMap(a.MapsPath) {
		var err error
		a.Maps[fileName], err = ebiten.NewImageFromImage(ETMhelper.TiledMapToImage(&mapFile), ebiten.FilterDefault)

		if err != nil {
			panic(err)
		}
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
