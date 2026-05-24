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

func GetElementByHashMap(Elements []*Sprite) map[string][]*Sprite {
	result := make(map[string][]*Sprite)
	for _, element := range Elements {
		key := string(int(element.Pos[0]/Conf.CameraOffset[0])) + "_" + string(int(element.Pos[1]/Conf.CameraOffset[1]))
		result[key] = append(result[key], element)
	}
	return result
}
