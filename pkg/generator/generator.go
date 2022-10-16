package generator

import (
	"archive/zip"
	"encoding/json"
	"io"
	"time"

	"github.com/LouisBrunner/timberborn-map-generator/pkg/timberborn"
)

type MapOptions struct {
	Width  int
	Height int
	Seed   int64
}

type Generator interface {
	Generate(w io.Writer, options MapOptions) error
}

type generator struct {
}

func NewGenerator() Generator {
	return &generator{}
}

func (me *generator) generateMap(options MapOptions) (*timberborn.Map, error) {
	topology, err := me.generateTopology(options)
	if err != nil {
		return nil, err
	}

	entities, err := me.generateEntities(options, topology)
	if err != nil {
		return nil, err
	}

	return &timberborn.Map{
		GameVersion: timberborn.GameVersion,
		Timestamp: timberborn.MapTime{
			Time: time.Now(),
		},
		Entities: entities,
		Singletons: timberborn.MapSingletons{
			MapSize: timberborn.MapSize{
				Size: timberborn.NewVector2(options.Width, options.Height),
			},
			TerrainMap: timberborn.MapTerrainMap{
				Heights: *topology,
			},
			WaterMap: timberborn.MapWaterMap{
				WaterDepths: timberborn.NewMapArray[int](options.Width, options.Height),
				Outflows:    timberborn.NewMapArray[timberborn.MapOutflow](options.Width, options.Height),
			},
			SoilMoistureSimulator: timberborn.MapSoilMoistureSimulator{
				MoistureLevels: timberborn.NewMapArray[int](options.Width, options.Height),
			},
			CameraStateRestorer: timberborn.MapCameraStateRestorer{
				SavedCameraState: timberborn.MapSavedCameraState{
					Target:          timberborn.NewVector3(0, 0, 0),
					ZoomLevel:       0,
					HorizontalAngle: 30,
					VerticalAngle:   70,
				},
			},
		},
	}, nil
}

func (me *generator) Generate(w io.Writer, options MapOptions) error {
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
