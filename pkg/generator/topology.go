package generator

import (
	"math"

	"github.com/LouisBrunner/timberborn-map-generator/pkg/timberborn"
	"github.com/aquilax/go-perlin"
)

const (
	perlinAlpha = 1.8
	perlinBeta  = 2.1
	perlinN     = 3

	mapRatio = 4

	topologyMaxHeight  = 16
	baseLayer          = 4
	minimumRiverLength = 5
)

func (me *generator) generateTopology(options MapOptions) (*timberborn.MapArray[int], error) {
	topology := timberborn.NewMapArray[int](options.Width, options.Height)

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
