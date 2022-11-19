// package data implements functions to load, modidy and save
// data from different file extensions
package data

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

const dataDirLocation = "../data/"
const fieldsFileName = "fields.json"
const planetarySystemsFileName = "planetary-systems.csv"

var loadedCSVCols []string = nil
var loadedPlanetaryData [][]string = nil
var loadedFields map[string]string = nil

func LoadFields(forceReload bool) map[string]string {
	if loadedFields != nil && !forceReload {
		return loadedFields
	} else {
		loadedFields = nil

		bytes, err := os.ReadFile(dataDirLocation + fieldsFileName)
		if err != nil {
			panic(fmt.Sprintf("an error occured loading json data\n%s", err))
		}

		var fields map[string]string
		json.Unmarshal(bytes, &fields)

		loadedFields = fields
		return fields
	}
}

func LoadPlanetarySystems(forceReload bool) [][]string {
	_, data := LoadPlanetarySystemsDataSet(forceReload)
	return data
}

func LoadPlanetarySystemsDataSet(forceReload bool) ([]string, [][]string) {
	if loadedCSVCols != nil && loadedPlanetaryData != nil && !forceReload {
		return loadedCSVCols, loadedPlanetaryData
	} else {
		loadedCSVCols = nil
		loadedPlanetaryData = nil

		f, err := os.Open(dataDirLocation + planetarySystemsFileName)
		if err != nil {
			panic(fmt.Sprintf("an error occured loading csv data\n%s", err))
		}

		r := csv.NewReader(f)
		records, err := r.ReadAll()
		if err != nil {
			panic(fmt.Sprintf("an error occured processing CSV file\n%s", err))
		}

		loadedCSVCols = records[0]
		loadedPlanetaryData = records[1:]
		return loadedCSVCols, loadedPlanetaryData
	}
}

func FilterCSVColumns(data [][]string, selectedColIndices []int) [][]string {
	filteredRows := [][]string{}
	for _, row := range data {
		filteredRow := []string{}
		for _, colIndex := range selectedColIndices {
			filteredRow = append(filteredRow, row[colIndex])
		}
		filteredRows = append(filteredRows, filteredRow)
	}

	return filteredRows
}
