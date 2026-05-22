package ETEngine

type Game struct {
	Elements []*Sprite
	Map      *Sprite
}

type Sprite struct {
	Image     string
	Pos, Size [2]float64
	Rotation  float64
	Layer     int
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
}

type Map struct {
	Map      string
	Elements []Sprite
}
