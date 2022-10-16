package generator

import "github.com/LouisBrunner/timberborn-map-generator/pkg/timberborn"

func (me *generator) findStart(options MapOptions, topology *timberborn.MapArray[int]) (timberborn.Vector3, error) {
	// TODO: wrong, need some kind of BFS search at the center of the map to find it
	return timberborn.NewVector3(options.Width/2, options.Height/2, baseLayer), nil
}
