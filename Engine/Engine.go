package ETEngine

import (
	"log"

	assets "github.com/Try-si/ETM/Assets"
	ETEStruct "github.com/Try-si/ETM/ETEStruct"
	ETMhelper "github.com/Try-si/ETM/Helper"
	"github.com/hajimehoshi/ebiten"
)

var (
	Gam    *ETEStruct.Game
	Assets assets.Assets
)

func Init(update func(float64) error) {
	Gam = &ETEStruct.Game{}

	Gam.Conf = ETMhelper.Jsontostruct[ETEStruct.Config]("config.json")
	Gam.UpdateFunc = update

	Assets = assets.Assets{}
	Assets.SpritePath = Gam.Conf.SpritePath
	Assets.MapsPath = Gam.Conf.MapsPath
	Assets.Init()

	Gam.Assets = &Assets

	Gam.Maps = make(map[string]ETEStruct.Map)

	Map := ETMhelper.GetAllFilesInDirectoryToStruct[ETEStruct.Map](Gam.Conf.JsonMapPath)

	for _, mapData := range Map {
		Gam.Maps[mapData.Name] = mapData.Obj
	}

	Gam.SetScene(Gam.Conf.StartMap, Assets.Maps[Gam.Conf.StartMap])
}

func GameLoop() {
	ebiten.SetWindowSize(Gam.Conf.ScreenWidth, Gam.Conf.ScreenHeight)
	ebiten.SetWindowTitle(Gam.Conf.Title)
	if err := ebiten.RunGame(Gam); err != nil {
		log.Fatal(err)
	}
}
