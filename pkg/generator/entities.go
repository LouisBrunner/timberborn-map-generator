package generator

import (
	"github.com/LouisBrunner/timberborn-map-generator/pkg/timberborn"
	"github.com/google/uuid"
)

func (me *generator) generateEntities(options MapOptions, topology *timberborn.MapArray[int]) ([]timberborn.MapEntity, error) {
	entities := []timberborn.MapEntity{}

	start, err := me.findStart(options, topology)
	if err != nil {
		return nil, err
	}
	entities = append(entities, timberborn.MapEntity{
		ID:       uuid.New(),
		Template: timberborn.MapTemplateStartingLocation,
		Components: &timberborn.MapEntityStartLocation{
			BlockObject: timberborn.MapBlockObject{
				Coordinates: start,
			},
		},
	})

	sources, err := me.findSources(options, topology)
	if err != nil {
		return nil, err
	}
	for _, source := range sources {
		entities = append(entities, timberborn.MapEntity{
			ID:       uuid.New(),
			Template: timberborn.MapTemplateWaterSource,
			Components: &timberborn.MapEntityWaterSource{
				BlockObject: timberborn.MapBlockObject{
					Coordinates: source,
				},
				WaterSource: timberborn.MapWaterSource{
					SpecifiedStrength: 8,
					CurrentStrength:   8,
				},
			},
		})
	}

	return entities, nil
}
