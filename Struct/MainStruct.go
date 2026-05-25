package ETEStruct

import (
	"sort"

	assets "github.com/Try-si/ETM/Assets"
	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	Elements []*Sprite
	Map      *Sprite

	UpdateFunc func(float64) error
	Assets     assets.Assets
	Maps       map[string]Map
	Conf       Config
}

type Sprite struct {
	Image     string
	Pos, Size [2]float64
	Rotation  float64
	Layer     int
	Box       [2]float64
}

type Config struct {
	ScreenWidth  int
	ScreenHeight int
	Title        string
	SpritePath   string
	MapsPath     string
	StartMap     string
	JsonMapPath  string
	CameraOffset [2]float64
	CellSize     float64
}

type Map struct {
	Map      string
	Elements []Sprite
}

func GetElementByHashMap(Elements []*Sprite, CameraOffset [2]float64) map[string][]*Sprite {
	result := make(map[string][]*Sprite)
	for _, element := range Elements {
		key := string(int(element.Pos[0]/CameraOffset[0])) + "_" + string(int(element.Pos[1]/CameraOffset[1]))
		result[key] = append(result[key], element)
	}
	return result
}

func (g *Game) SetScene(Map string) {
	mapImg := g.Assets.Maps[g.Maps[Map].Map]
	g.Map = &Sprite{
		Image:    g.Maps[Map].Map,
		Pos:      [2]float64{0, 0},
		Size:     [2]float64{float64(mapImg.Bounds().Dx()), float64(mapImg.Bounds().Dy())},
		Rotation: 0,
		Layer:    -100,
	}

	MapData := g.Maps[Map]
	for i := range MapData.Elements {
		g.Elements = append(g.Elements, &MapData.Elements[i])
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	return g.UpdateFunc(ebiten.CurrentFPS())
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, element := range g.GetElementByLayer() {
		img, exists := g.Assets.Images[element.Image]
		if !exists {
			img, exists = g.Assets.Maps[element.Image]

			if !exists {
				continue
			}
		}

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(element.Pos[0]-g.Conf.CameraOffset[0]), float64(element.Pos[1]+g.Conf.CameraOffset[1]))
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
	return g.Conf.ScreenWidth, g.Conf.ScreenHeight
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
