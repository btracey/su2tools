package main

import (
	"encoding/csv"
	"flag"

	"github.com/btracey/ransuq/dataloader"
)

func main() {
	var baseFile string
	flag.StringVar(&baseFile, "base", "", "base file off of which the old is based")
	var deltaFile string
	flag.StringVar(&deltaFile, "delta", "", "new file to be compared to the base file")
	var diffField string
	flag.StringVar(&diffField, "field", "", "which field should the difference be computed")
	var configFile string
	flag.StringVar(&config, "config", "", "config file for running SU2 Sol")

	flag.Parse()

	// Turn the files into dataloaders
	baseLoader := &dataloader.Dataset{
		Name:     "base",
		Filename: baseFile,
		Format:   dataloader.SU2_restart_2dturb{},
	}
	deltaLoader := &dataloader.Dataset{
		Name:     "delta",
		Filename: baseFile,
		Format:   dataloader.SU2_restart_2dturb{},
	}

	fields := []string{
		"PointID"
		"XLoc",
		"YLoc",
		diffField,
	}

	baseData := dataloader.LoadFromDataset(fields, baseLoader)
	deltaData := dataloader.LoadFromDataset(fields, deltaLoader)

	newData := make([][]float64, len(baseData))
	for i := range newData {
		newData[i] = make([]float64, len(baseData[i]))
	}

	for i := range newData {
		for j := range newData[i] {
			// check that the x and y locations are the same
		}
	}
}
