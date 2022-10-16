package generator

import (
	"fmt"

	"github.com/LouisBrunner/timberborn-map-generator/pkg/timberborn"
)

type getCoords func(i int) timberborn.Vector2

func (me *generator) findSource(options MapOptions, topology *timberborn.MapArray[int], maxDimension, maxElevation, minimumStreak int, getCoords getCoords) ([]timberborn.Vector3, error) {
	lastElevation := -1
	currentStreak := 0

	checkFound := func(index int) []timberborn.Vector3 {
		if currentStreak < minimumStreak {
			return nil
		}

		riverBed := make([]timberborn.Vector3, currentStreak)
		for i := 0; i < currentStreak; i += 1 {
			riverBed[i] = timberborn.Vector3{
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

func (me *generator) findSources(options MapOptions, topology *timberborn.MapArray[int]) ([]timberborn.Vector3, error) {
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
			getCoords: func(i int) timberborn.Vector2 {
				return timberborn.NewVector2(i, 0)
			},
		},
		{
			name:      "bottom",
			dimension: options.Width,
			getCoords: func(i int) timberborn.Vector2 {
				return timberborn.NewVector2(i, options.Height-1)
			},
		},
		{
			name:      "left",
			dimension: options.Height,
			getCoords: func(i int) timberborn.Vector2 {
				return timberborn.NewVector2(0, i)
			},
		},
		{
			name:      "right",
			dimension: options.Height,
			getCoords: func(i int) timberborn.Vector2 {
				return timberborn.NewVector2(options.Width-1, i)
			},
		},
	}

	var soFar []timberborn.Vector3

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
