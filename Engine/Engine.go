package ETEngine

import (
	"log"
	"sort"

	assets "github.com/Try-si/ETM/Assets"
	ETMhelper "github.com/Try-si/ETM/Helper"
	"github.com/hajimehoshi/ebiten"
)

var (
	Conf       Config
	Maps       map[string]Map
	Assets     assets.Assets
	Gam        *Game
	UpdateFunc func(float64) error
)

func Init(update func(float64) error) {
	Conf = ETMhelper.Jsontostruct[Config]("config.json")
	UpdateFunc = update

	Assets = assets.Assets{}
	Assets.SpritePath = Conf.SpritePath
	Assets.MapsPath = Conf.MapsPath
	Assets.Init()

	Maps = make(map[string]Map)

	Map := ETMhelper.GetAllFilesInDirectoryToStruct[Map](Conf.JsonMapPath)

	for _, mapData := range Map {
		Maps[mapData.Name] = mapData.Obj
	}

	Gam = &Game{}

	Gam.SetScene(Conf.StartMap)
}

func GameLoop() {
	ebiten.SetWindowSize(Conf.ScreenWidth, Conf.ScreenHeight)
	ebiten.SetWindowTitle(Conf.Title)
	if err := ebiten.RunGame(Gam); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	return UpdateFunc(ebiten.CurrentFPS())
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, element := range g.GetElementByLayer() {
		img, exists := Assets.Images[element.Image]
		if !exists {
			img, exists = Assets.Maps[element.Image]

			if !exists {
				continue
			}
		}

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(element.Pos[0]-Conf.CameraOffset[0]), float64(element.Pos[1]+Conf.CameraOffset[1]))
		opts.GeoM.Rotate(element.Rotation)
		if element.Size[0] != 0 && element.Size[1] != 0 {
			opts.GeoM.Scale(float64(element.Size[0])/float64(img.Bounds().Dx()), float64(element.Size[1])/float64(img.Bounds().Dy()))
		} else {
			opts.GeoM.Scale(1, 1)
		}
		screen.DrawImage(img, opts)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return Conf.ScreenWidth, Conf.ScreenHeight
}

func (g *Game) GetElementByLayer() []*Sprite {
	// Combine Maps and Elements
	allSprites := append([]*Sprite{g.Map}, g.Elements...)

	// Sort by layer (smallest to largest)
	sort.Slice(allSprites, func(i, j int) bool {
		return allSprites[i].Layer <= allSprites[j].Layer
	})

	return allSprites
}

func (g *Game) SetScene(Map string) {
	mapImg := Assets.Maps[Maps[Map].Map]
	g.Map = &Sprite{
		Image:    Maps[Map].Map,
		Pos:      [2]float64{0, 0},
		Size:     [2]float64{float64(mapImg.Bounds().Dx()), float64(mapImg.Bounds().Dy())},
		Rotation: 0,
		Layer:    -100,
	}

	MapData := Maps[Map]
	for i := range MapData.Elements {
		g.Elements = append(g.Elements, &MapData.Elements[i])
	}
}
