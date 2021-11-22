package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/akamensky/argparse"
	gonav "github.com/pnxenopoulos/csgonavparse"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	argparser := argparse.NewParser("navtransitiongen", "Generate json connections needed to generate round map control visualisations")

	path := argparser.String("p", "path", &argparse.Options{Required: true, Help: ".nav filepath"})
	outputPath := argparser.String("o", "output", &argparse.Options{Required: true, Help: "where to write the generated .json file"})

	err := argparser.Parse(os.Args)
	if err != nil {
		fmt.Print(argparser.Usage(err))
		os.Exit(1)
	}

	inputPathExists, err := exists(*path)
	if err != nil {
		fmt.Println("Failed to stat input path.")
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	if !inputPathExists {
		fmt.Println("The input path does not exist.")
		os.Exit(1)
	}

	outputPathExists, err := exists(*outputPath)
	if err != nil {
		fmt.Println("Failed to stat output path.")
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	if outputPathExists {
		fmt.Println("The output path already exists")
		os.Exit(1)
	}

	f, _ := os.Open(*path)

	parser := gonav.Parser{Reader: f}
	mesh, _ := parser.Parse()

	output := make(map[string][]string)

	for _, area := range mesh.QuadTreeAreas.Areas {
		connections := make([]string, 0)
		for _, connection := range area.Connections {
			var stringAreaID = strconv.FormatUint(uint64(connection.TargetAreaID), 10)
			connections = append(connections, stringAreaID)
		}

		stringAreaID := strconv.FormatUint(uint64(area.ID), 10)
		output[stringAreaID] = connections
	}

	jsonOutput, _ := json.Marshal(output)

	err = ioutil.WriteFile(*outputPath, []byte(string(jsonOutput)), 0700)

	if err != nil {
		fmt.Println("Failed to write output .json file.")
		fmt.Printf("%v\n", err)
	}

}
