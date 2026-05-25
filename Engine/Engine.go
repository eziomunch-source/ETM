package ETEngine

import (
	"log"

	assets "github.com/Try-si/ETM/Assets"
	ETMhelper "github.com/Try-si/ETM/Helper"
	ETEStruct "github.com/Try-si/ETM/Struct"
	"github.com/hajimehoshi/ebiten"
)

var (
	Gam *ETEStruct.Game
)

func Init(update func(float64) error) {
	Gam = &ETEStruct.Game{}

	Gam.Conf = ETMhelper.Jsontostruct[ETEStruct.Config]("config.json")
	Gam.UpdateFunc = update

	Gam.Assets = assets.Assets{}
	Gam.Assets.SpritePath = Gam.Conf.SpritePath
	Gam.Assets.MapsPath = Gam.Conf.MapsPath
	Gam.Assets.Init()

	Gam.Maps = make(map[string]ETEStruct.Map)

	Map := ETMhelper.GetAllFilesInDirectoryToStruct[ETEStruct.Map](Gam.Conf.JsonMapPath)

	for _, mapData := range Map {
		Gam.Maps[mapData.Name] = mapData.Obj
	}

	Gam.SetScene(Gam.Conf.StartMap)
}

func GameLoop() {
	ebiten.SetWindowSize(Gam.Conf.ScreenWidth, Gam.Conf.ScreenHeight)
	ebiten.SetWindowTitle(Gam.Conf.Title)
	if err := ebiten.RunGame(Gam); err != nil {
		log.Fatal(err)
	}
}
