package timberborn

import (
	"github.com/google/uuid"
)

type Map struct {
	GameVersion string
	Timestamp   MapTime
	Singletons  MapSingletons
	Entities    []MapEntity
}

type MapSingletons struct {
	MapSize               MapSize
	TerrainMap            MapTerrainMap
	CameraStateRestorer   MapCameraStateRestorer
	WaterMap              MapWaterMap
	SoilMoistureSimulator MapSoilMoistureSimulator
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
	// TODO: loads missing
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

// TODO: loads missing
