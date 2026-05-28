package ETEStruct

import (
	"image/color"
	"sort"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type AssetsProvider interface {
	GetImage(name string) (*ebiten.Image, bool)
	GetMap(name string) (*ebiten.Image, bool)
}

type Game struct {
	Elements []*Sprite // tous les éléments dans le jeu
	Map      *Sprite   // map actuelle
	Mape     string    // nom de la map actuelle

	UpdateFunc func(float64) error // fonction de mise à jour
	Assets     AssetsProvider      // fournisseur d'assets
	Maps       map[string]Map      // toutes les maps
	Conf       Config              // configuration

	Debug bool // debug mode
}

type Sprite struct {
	Image     string     // image name
	Pos, Size [2]float64 // position and size
	Rotation  float64    // rotation
	Layer     int        // layer
	Box       [4]float64 // bounding box
}

type Config struct {
	ScreenWidth  int     // largeur de l'écran
	ScreenHeight int     // hauteur de l'écran
	Title        string  // titre de la fenêtre
	SpritePath   string  // chemin vers les sprites
	MapsPath     string  // chemin vers les maps
	StartMap     string  // map de départ
	JsonMapPath  string  // chemin vers les maps json
	CellSize     float64 // taille de la cellule
	Unité        int     // taille d'une unité en pixels
	Cam          Caméra  // caméra
}

type Caméra struct {
	Zoom, Offset [2]float64 // Zoom et offset de cam
}

type Map struct {
	Map      string   // nom de la map
	Elements []Sprite // éléments dans la map
	TileSize int      // taille des tuiles
}

func GetElementByHashMap(Elements []*Sprite, cellSize float64) map[string][]*Sprite { // regrouper les éléments par hash map
	result := make(map[string][]*Sprite) // initialiser la map de résultat
	for _, element := range Elements {   // itérer sur tous les éléments
		key := string(int(element.Pos[0]/cellSize)) + "_" + string(int(element.Pos[1]/cellSize)) // clé de la map
		result[key] = append(result[key], element)                                               // ajouter l'élément à la map
	}
	return result // retourner la hashmap
}

func drawCircle(screen *ebiten.Image, x, y, radius float64, clr color.Color) { // dessiner un cercle
	centerX, centerY := int(x), int(y) // centre du cercle
	r := int(radius)                   // rayon du cercle

	for dy := -r; dy <= r; dy++ { // itérer sur tous les pixels du cercle
		for dx := -r; dx <= r; dx++ {
			if dx*dx+dy*dy <= r*r { // vérifier si le pixel est dans le cercle
				screen.Set(centerX+dx, centerY+dy, clr) // dessiner le pixel
			}
		}
	}
}

func (g *Game) SetScene(Map string, MapImage *ebiten.Image) { // définir la scène
	g.Map = &Sprite{ // créer la map
		Image:    g.Maps[Map].Map,
		Pos:      [2]float64{0, 0},
		Size:     [2]float64{float64(MapImage.Bounds().Dx()) / float64(g.Conf.Unité), float64(MapImage.Bounds().Dy()) / float64(g.Conf.Unité)},
		Rotation: 0,
		Layer:    -100,
	}

	MapData := g.Maps[Map] // obtenir les données de la map
	for i := range MapData.Elements {
		g.Elements = append(g.Elements, &MapData.Elements[i]) // ajouter les éléments à la scène
	}
}

func (g *Game) Update() error {
	return g.UpdateFunc(ebiten.ActualFPS())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.Conf.ScreenWidth, g.Conf.ScreenHeight
}

func (g *Game) GetElementByLayer() []*Sprite { // obtenir les éléments par couche
	// Combine Maps and Elements
	allSprites := append([]*Sprite{g.Map}, g.Elements...) // combiner les maps et les éléments

	// Sort by layer (smallest to largest)
	sort.Slice(allSprites, func(i, j int) bool { // trier par couche (du plus petit au plus grand)
		return allSprites[i].Layer <= allSprites[j].Layer // retourner true si le layer de i est plus petit que celui de j
	})

	return allSprites // retourner les éléments triés
}

func (g *Game) Draw(screen *ebiten.Image) { // dessiner la scène

	for _, element := range g.GetElementByLayer() { // itérer sur tous les éléments
		img, exists := g.Assets.GetImage(element.Image) // obtenir l'image
		if !exists {                                    // si l'image n'existe pas
			img, exists = g.Assets.GetMap(element.Image) // obtenir la map

			exists = !exists // inverser la valeur de exists

			if exists { // si la map n'existe pas
				continue
			}
		}
		if g.Debug { // si le mode debug est activé
			posX := (element.Pos[0] - g.Conf.Cam.Offset[0]) * float64(g.Conf.Unité) // calculer la position x en pixels
			posY := (element.Pos[1] + g.Conf.Cam.Offset[1]) * float64(g.Conf.Unité) // calculer la position y en pixels

			whith := element.Box[0] * float64(g.Conf.Unité)  // obtenir la largeur de la hitbox en pixels
			height := element.Box[1] * float64(g.Conf.Unité) // obtenir la hauteur de la hitbox en pixels

			if element.Box[2] != 0 { // si xOffset est différent de 0
				posX += element.Box[2] * float64(g.Conf.Unité)
			}
			if element.Box[3] != 0 { // si yOffset est différent de 0
				posY += element.Box[3] * float64(g.Conf.Unité)
			}

			if whith == 0 && height == 0 { // si la hitbox n'est pas définie
				continue
			} else if height == 0 { // si la hitbox est un cercle
				//Draw circle
				drawCircle(screen, posX, posY, whith, color.RGBA{255, 255, 255, 128})
			} else { // si la hitbox est un rectangle
				//Draw rectangle
				DrawRect(screen, posX, posY, whith, height, color.RGBA{255, 255, 255, 128})
			}
		}

		opts := &ebiten.DrawImageOptions{}

		// 1. Centrer sur l'origine (avant scale)
		width := float64(img.Bounds().Dx())               // largeur de l'image en pixels
		height := float64(img.Bounds().Dy())              // hauteur de l'image en pixels
		if element.Size[0] != 0 && element.Size[1] != 0 { // si la taille est définie
			width = element.Size[0] * float64(g.Conf.Unité)  // largeur en pixels
			height = element.Size[1] * float64(g.Conf.Unité) // hauteur en pixels
		}
		opts.GeoM.Translate(-width/2, -height/2) // centrer sur l'origine

		// 2. Scale (avec zoom)
		if element.Size[0] != 0 && element.Size[1] != 0 { // si la taille est définie
			opts.GeoM.Scale(float64(element.Size[0]*float64(g.Conf.Unité))/float64(img.Bounds().Dx()), float64(element.Size[1]*float64(g.Conf.Unité))/float64(img.Bounds().Dy()))
			// scale with element size : element.Size = taille en unité, * g.Conf.Unité = mettre taille en pixels, / img.Bounds().Dx() = scale
		} else {
			opts.GeoM.Scale(float64(g.Conf.Unité), float64(g.Conf.Unité))
		}

		// 3. Rotate
		opts.GeoM.Rotate(element.Rotation) // rotate

		// 4. Translate vers la position finale (sans zoom dans la translation)
		opts.GeoM.Translate(float64(element.Pos[0])*float64(g.Conf.Unité), float64(element.Pos[1])*float64(g.Conf.Unité)) // translate

		// 5. Camera offset (avec zoom)
		opts.GeoM.Scale(g.Conf.Cam.Zoom[0], g.Conf.Cam.Zoom[1])                                                     // Zoom
		opts.GeoM.Translate(g.Conf.Cam.Offset[0]*float64(g.Conf.Unité), g.Conf.Cam.Offset[1]*float64(g.Conf.Unité)) // Center

		screen.DrawImage(img, opts) // dessiner l'image
	}
}

func DrawRect(screen *ebiten.Image, x, y, width, height float64, clr color.Color) {
	vector.FillRect(screen, float32(x-width/2), float32(y-height/2), float32(width), float32(height), clr, false)
}
