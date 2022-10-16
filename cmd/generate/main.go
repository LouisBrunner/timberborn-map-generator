package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

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

func (me *Generator) generateTopology(options MapOptions) (*MapArray[int], error) {
	topology := NewMapArray[int](options.Width, options.Height)

	for i := 0; i < options.Width; i += 1 {
		for j := 0; j < options.Height; j += 1 {
			err := topology.Set(i, j, 4)
			if err != nil {
				return nil, err
			}
		}
	}

	return &topology, nil
}

func (me *Generator) generateMap(options MapOptions) (*Map, error) {
	start := MapEntity{
		ID:       uuid.New(),
		Template: MapTemplateStartingLocation,
		Components: &MapEntityStartLocation{
			BlockObject: MapBlockObject{
				Coordinates: Vector3{
					Vector2: Vector2{
						X: options.Width / 2,
						Y: options.Width / 2,
					},
					Z: 4,
				},
			},
		},
	}

	topology, err := me.generateTopology(options)
	if err != nil {
		return nil, err
	}

	return &Map{
		GameVersion: GameVersion,
		Timestamp: MapTime{
			Time: time.Now(),
		},
		Entities: []MapEntity{
			start,
		},
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

func main() {
	err := GenerateMap(MapOptions{
		Width:  16,
		Height: 16,
	}, "/Users/louis/Documents/Timberborn/Maps/Test2.timber")
	if err != nil {
		panic(err.Error())
	}
}
