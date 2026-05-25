package ETEStruct

import (
	"image/color"
	"sort"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type AssetsProvider interface {
	GetImage(name string) (*ebiten.Image, bool)
	GetMap(name string) (*ebiten.Image, bool)
}

type Game struct {
	Elements []*Sprite
	Map      *Sprite

	UpdateFunc func(float64) error
	Assets     AssetsProvider
	Maps       map[string]Map
	Conf       Config
}

type Sprite struct {
	Image     string
	Pos, Size [2]float64
	Rotation  float64
	Layer     int
	Box       [4]float64
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

func GetElementByHashMap(Elements []*Sprite, cellSize float64) map[string][]*Sprite {
	result := make(map[string][]*Sprite)
	for _, element := range Elements {
		key := string(int(element.Pos[0]/cellSize)) + "_" + string(int(element.Pos[1]/cellSize))
		result[key] = append(result[key], element)
	}
	return result
}

func drawCircle(screen *ebiten.Image, x, y, radius float64, clr color.Color) {
	centerX, centerY := int(x+radius), int(y+radius)
	r := int(radius)

	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			if dx*dx+dy*dy <= r*r {
				screen.Set(centerX+dx, centerY+dy, clr)
			}
		}
	}
}

func (g *Game) SetScene(Map string, MapImage *ebiten.Image) {
	g.Map = &Sprite{
		Image:    g.Maps[Map].Map,
		Pos:      [2]float64{0, 0},
		Size:     [2]float64{float64(MapImage.Bounds().Dx()), float64(MapImage.Bounds().Dy())},
		Rotation: 0,
		Layer:    -100,
	}

	MapData := g.Maps[Map]
	for i := range MapData.Elements {
		g.Elements = append(g.Elements, &MapData.Elements[i])
	}
}

func (g *Game) Update() error {
	return g.UpdateFunc(ebiten.ActualFPS())
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

func (g *Game) Draw(screen *ebiten.Image) {
	for _, element := range g.GetElementByLayer() {
		img, exists := g.Assets.GetImage(element.Image)
		if !exists {
			img, exists = g.Assets.GetMap(element.Image)

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

	for _, element := range g.Elements {
		posX := (element.Pos[0] - g.Conf.CameraOffset[0])
		posY := (element.Pos[1] + g.Conf.CameraOffset[1])

		whith := element.Box[0]
		height := element.Box[1]

		if element.Box[2] != 0 {
			posX += element.Box[2]
		}
		if element.Box[3] != 0 {
			posY += element.Box[3]
		}

		if whith == 0 && height == 0 {
			continue
		} else if height == 0 {
			//Draw circle
			drawCircle(screen, posX, posY, whith, color.RGBA{255, 255, 255, 128})
		} else {
			//Draw rectangle
			ebitenutil.DrawRect(screen, posX, posY, whith, height, color.RGBA{255, 255, 255, 128})
		}
	}
}
