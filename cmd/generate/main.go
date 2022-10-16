package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/aquilax/go-perlin"
	"github.com/google/uuid"
)

const GameVersion = "0.2.9.1-0b5fdc2-sm"

type Map struct {
	GameVersion string
	Timestamp   MapTime
	Singletons  MapSingletons
	Entities    []MapEntity
}

type MapTime struct {
	Time time.Time
}

type MapSingletons struct {
	MapSize               MapSize
	TerrainMap            MapTerrainMap
	CameraStateRestorer   MapCameraStateRestorer
	WaterMap              MapWaterMap
	SoilMoistureSimulator MapSoilMoistureSimulator
}

type MapArray[T any] struct {
	width   int
	content []T
}

type MapTerrainMap struct {
	Heights MapArray[int]
}

type MapCameraStateRestorer struct {
	SavedCameraState MapSavedCameraState
}

type MapSavedCameraState struct {
	Target          Vector3
	ZoomLevel       int
	HorizontalAngle int
	VerticalAngle   int
}

type MapOutflow struct {
	A int
	B int
	C int
	D int
}

type MapWaterMap struct {
	WaterDepths MapArray[int]
	Outflows    MapArray[MapOutflow]
}

type MapSoilMoistureSimulator struct {
	MoistureLevels MapArray[int]
}

type MapSize struct {
	Size Vector2
}

type MapEntity struct {
	ID         uuid.UUID `json:"Id"`
	Template   MapTemplateID
	Components any
}

type Vector2 struct {
	X int
	Y int
}

type Vector3 struct {
	Vector2
	Z int
}

type MapTemplateID string

const (
	MapTemplateStartingLocation MapTemplateID = "StartingLocation"
	MapTemplateSlope            MapTemplateID = "Slope"
	MapTemplateWaterSource      MapTemplateID = "WaterSource"
	MapTemplateBarrier          MapTemplateID = "Barrier"
	MapTemplateUndergroundRuins MapTemplateID = "UndergroundRuins"
	MapTemplateRuinColumnH8     MapTemplateID = "RuinColumnH8"
	MapTemplateBirch            MapTemplateID = "Birch"
	MapTemplateBlueberryBush    MapTemplateID = "BlueberryBush"
)

type MapBlockObject struct {
	Coordinates Vector3
}

type MapWaterSource struct {
	SpecifiedStrength int
	CurrentStrength   int
}

type MapDryObject struct {
	IsDry bool
}

type MapEntityStartLocation struct {
	BlockObject MapBlockObject
}

type MapEntityBarrier struct {
	BlockObject MapBlockObject
}

type MapEntityWaterSource struct {
	BlockObject MapBlockObject
	WaterSource MapWaterSource
}

type MapEntityUndergroundRuins struct {
	BlockObject MapBlockObject
	DryObject   MapDryObject
}

type MapEntitySlope struct {
	BlockObject MapBlockObject
	// TODO: incomplete
}

type MapEntityRuinColumn struct {
	BlockObject MapBlockObject
	DryObject   MapDryObject
	// TODO: incomplete
}

type MapEntityBirch struct {
	BlockObject MapBlockObject
	DryObject   MapDryObject
	// TODO: incomplete
}

type MapEntityBlueberryBush struct {
	BlockObject MapBlockObject
	DryObject   MapDryObject
	// TODO: incomplete
}

type MapOptions struct {
	Width  int
	Height int
	Seed   int64
}

type Generator struct {
}

func NewMapArray[T any](width, height int) MapArray[T] {
	return MapArray[T]{
		width:   width,
		content: make([]T, width*height),
	}
}

func (me *MapArray[T]) Set(x, y int, value T) error {
	position := y*me.width + x
	if position >= len(me.content) {
		return fmt.Errorf("could not set %v,%v as it is out-of-range", x, y)
	}
	me.content[position] = value
	return nil
}

func (me *MapArray[T]) Get(x, y int) (T, error) {
	position := y*me.width + x
	if position >= len(me.content) {
		return me.content[0], fmt.Errorf("could not set %v,%v as it is out-of-range", x, y)
	}
	return me.content[position], nil
}

func (me MapArray[T]) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(`{"Array":"`)
	for i, v := range me.content {
		raw := fmt.Sprint(v)
		buf.WriteString(raw)
		if i+1 < len(me.content) {
			buf.WriteString(` `)
		}
	}
	buf.WriteString(`"}`)

	return buf.Bytes(), nil
}

func (me MapOutflow) String() string {
	// FIXME: just need a generic ForEach
	return strings.Join([]string{
		fmt.Sprint(me.A),
		fmt.Sprint(me.B),
		fmt.Sprint(me.C),
		fmt.Sprint(me.D),
	}, ":")
}

func (me MapTime) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(`"`)
	buf.WriteString(me.Time.Format("2006-01-02 15:04:05"))
	buf.WriteString(`"`)

	return buf.Bytes(), nil
}

const (
	perlinAlpha = 1.8
	perlinBeta  = 2.1
	perlinN     = 3

	mapRatio = 4

	topologyMaxHeight  = 16
	baseLayer          = 4
	minimumRiverLength = 5
)

func (me *Generator) generateTopology(options MapOptions) (*MapArray[int], error) {
	topology := NewMapArray[int](options.Width, options.Height)

	prln := perlin.NewPerlin(perlinAlpha, perlinBeta, perlinN, options.Seed)

	// minV := -math.Sqrt(float64(perlinN) / 4)
	maxV := math.Sqrt(float64(perlinN)) / 2

	for i := 0; i < options.Width; i += 1 {
		for j := 0; j < options.Height; j += 1 {
			// calculate perlin
			v := prln.Noise2D(float64(i)/float64(options.Width/mapRatio), float64(j)/float64(options.Height/mapRatio))
			// move to a -1,1 range
			v = (v + maxV) / (maxV * 2)
			// move to 0,MAX range and round
			v = math.Round(v * topologyMaxHeight)
			// check range
			v = math.Max(v, 0)
			v = math.Min(v, topologyMaxHeight)

			err := topology.Set(i, j, int(v))
			if err != nil {
				return nil, err
			}
		}
	}

	return &topology, nil
}

func (me *Generator) findStart(options MapOptions, topology *MapArray[int]) (*Vector3, error) {
	// TODO: wrong, need some kind of BFS search at the center of the map to find it
	return &Vector3{
		Vector2: Vector2{
			X: options.Width / 2,
			Y: options.Height / 2,
		},
		Z: baseLayer,
	}, nil
}

type getCoords func(i int) Vector2

func (me *Generator) findSource(options MapOptions, topology *MapArray[int], maxDimension, maxElevation, minimumStreak int, getCoords getCoords) ([]Vector3, error) {
	lastElevation := -1
	currentStreak := 0

	checkFound := func(index int) []Vector3 {
		if currentStreak < minimumStreak {
			return nil
		}

		riverBed := make([]Vector3, currentStreak)
		for i := 0; i < currentStreak; i += 1 {
			riverBed[i] = Vector3{
				Vector2: getCoords(index - currentStreak + i),
				Z:       lastElevation,
			}
		}
		return riverBed
	}

	for i := 0; i < maxDimension; i += 1 {
		coords := getCoords(i)
		elevation, err := topology.Get(coords.X, coords.Y)
		if err != nil {
			return nil, err
		}
		if elevation <= maxElevation && lastElevation == elevation {
			currentStreak += 1
		} else {
			if found := checkFound(i); found != nil {
				return found, nil
			}
			currentStreak = 1
		}
		lastElevation = elevation
	}

	if found := checkFound(maxDimension - 1); found != nil {
		return found, nil
	}

	return nil, fmt.Errorf("not found")
}

func (me *Generator) findSources(options MapOptions, topology *MapArray[int]) ([]Vector3, error) {
	maxElevation := baseLayer - 1
	minimumStreak := minimumRiverLength

	cases := []struct {
		name      string
		dimension int
		getCoords getCoords
	}{
		{
			name:      "top",
			dimension: options.Width,
			getCoords: func(i int) Vector2 {
				return Vector2{X: i, Y: 0}
			},
		},
		{
			name:      "bottom",
			dimension: options.Width,
			getCoords: func(i int) Vector2 {
				return Vector2{X: i, Y: options.Height - 1}
			},
		},
		{
			name:      "left",
			dimension: options.Height,
			getCoords: func(i int) Vector2 {
				return Vector2{X: 0, Y: i}
			},
		},
		{
			name:      "right",
			dimension: options.Height,
			getCoords: func(i int) Vector2 {
				return Vector2{X: options.Width - 1, Y: i}
			},
		},
	}

	var soFar []Vector3

	for _, icase := range cases {
		sources, err := me.findSource(options, topology, icase.dimension, maxElevation, minimumStreak, icase.getCoords)
		if err == nil {
			if soFar == nil {
				soFar = sources
				continue
			}
			soFar = append(soFar, sources...)
			return soFar, nil
		}
	}

	return nil, fmt.Errorf("could not generate water sources, try another seed")
}

func (me *Generator) generateEntities(options MapOptions, topology *MapArray[int]) ([]MapEntity, error) {
	entities := []MapEntity{}

	start, err := me.findStart(options, topology)
	if err != nil {
		return nil, err
	}
	entities = append(entities, MapEntity{
		ID:       uuid.New(),
		Template: MapTemplateStartingLocation,
		Components: &MapEntityStartLocation{
			BlockObject: MapBlockObject{
				Coordinates: Vector3{
					Vector2: Vector2{
						X: start.X,
						Y: start.Y,
					},
					Z: start.Z,
				},
			},
		},
	})

	sources, err := me.findSources(options, topology)
	if err != nil {
		return nil, err
	}
	for _, source := range sources {
		entities = append(entities, MapEntity{
			ID:       uuid.New(),
			Template: MapTemplateWaterSource,
			Components: &MapEntityWaterSource{
				BlockObject: MapBlockObject{
					Coordinates: Vector3{
						Vector2: Vector2{
							X: source.X,
							Y: source.Y,
						},
						Z: source.Z,
					},
				},
				WaterSource: MapWaterSource{
					SpecifiedStrength: 8,
					CurrentStrength:   8,
				},
			},
		})
	}

	return entities, nil
}

func (me *Generator) generateMap(options MapOptions) (*Map, error) {
	topology, err := me.generateTopology(options)
	if err != nil {
		return nil, err
	}

	entities, err := me.generateEntities(options, topology)
	if err != nil {
		return nil, err
	}

	return &Map{
		GameVersion: GameVersion,
		Timestamp: MapTime{
			Time: time.Now(),
		},
		Entities: entities,
		Singletons: MapSingletons{
			MapSize: MapSize{
				Size: Vector2{
					X: options.Width,
					Y: options.Height,
				},
			},
			TerrainMap: MapTerrainMap{
				Heights: *topology,
			},
			WaterMap: MapWaterMap{
				WaterDepths: NewMapArray[int](options.Width, options.Height),
				Outflows:    NewMapArray[MapOutflow](options.Width, options.Height),
			},
			SoilMoistureSimulator: MapSoilMoistureSimulator{
				MoistureLevels: NewMapArray[int](options.Width, options.Height),
			},
			CameraStateRestorer: MapCameraStateRestorer{
				SavedCameraState: MapSavedCameraState{
					Target: Vector3{
						Vector2: Vector2{
							X: 0,
							Y: 0,
						},
						Z: 0,
					},
					ZoomLevel:       0,
					HorizontalAngle: 30,
					VerticalAngle:   70,
				},
			},
		},
	}, nil
}

func (me *Generator) Generate(w io.Writer, options MapOptions) error {
	writer := zip.NewWriter(w)
	defer writer.Close()

	world, err := writer.Create("world.json")
	if err != nil {
		return err
	}

	timber, err := me.generateMap(options)
	if err != nil {
		return err
	}

	buffer, err := json.Marshal(timber)
	if err != nil {
		return err
	}

	_, err = world.Write(buffer)
	return err
}

func GenerateMap(options MapOptions, output string) error {
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	generator := Generator{}
	return generator.Generate(f, options)
}

func usage() {
	fmt.Printf("Usage: %s [opts] filename\n", os.Args[0])
	fmt.Printf("options:\n")
	flag.PrintDefaults()
}

func main() {
	var width, height int
	var seed int64
	defaultSeed := time.Now().UnixMilli() * int64(os.Getpid())
	flag.Usage = usage
	flag.IntVar(&width, "width", 256, "width of the map")
	flag.IntVar(&height, "height", 256, "height of the map")
	flag.Int64Var(&seed, "seed", defaultSeed, "seed used for generation")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Printf("error: missing filename\n")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Seed: %v\n", seed)

	err := GenerateMap(MapOptions{
		Width:  256,
		Height: 256,
		Seed:   defaultSeed,
	}, flag.Args()[0])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
