package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/LouisBrunner/timberborn-map-generator/pkg/generator"
)

func GenerateMap(options generator.MapOptions, output string) error {
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	generator := generator.NewGenerator()
	return generator.Generate(f, options)
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [opts] filename\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "options:\n")
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

	err := GenerateMap(generator.MapOptions{
		Width:  256,
		Height: 256,
		Seed:   defaultSeed,
	}, flag.Args()[0])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
