package hierarchy

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type Hierarchy struct {
	ChildToParent map[string]string // child name --> parent name

	WeightsPerEpoch []map[string]int64 // Element of slice representing change of tree per epoch.
	// Key of map is node name.
	// Value of map is the number
	// of weight by this node per epoch.
}

func NewHierarchy(csvChildParent string, csvWeightsPerEpoch string) *Hierarchy {
	var newHierarchy Hierarchy

	newHierarchy.ChildToParent = readTreeNodes(csvChildParent)
	newHierarchy.WeightsPerEpoch = readWeightsPerEpoch(csvWeightsPerEpoch)

	return &newHierarchy
}

func readTreeNodes(csvChildParent string) map[string]string {
	var r = newCsvReader(csvChildParent)

	var childToParent = make(map[string]string)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		childName, parentName := record[0], record[1]

		childToParent[childName] = parentName
	}
	return childToParent
}

func readWeightsPerEpoch(csvWeightPerEpoch string) []map[string]int64 {
	var r = newCsvReader(csvWeightPerEpoch)
	var weightsPerEpoch = make([]map[string]int64, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var nodeName = record[0]
		var epoch, _ = strconv.Atoi(record[1])
		var weight, _ = strconv.ParseInt(record[2], 10, 64)

		if epoch == len(weightsPerEpoch) {
			weightsPerEpoch = append(weightsPerEpoch, make(map[string]int64))
		}

		weightsPerEpoch[epoch][nodeName] += weight
	}
	return weightsPerEpoch
}

func newCsvReader(filepath string) *csv.Reader {
	csvFile, err := os.Open(filepath)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	r := csv.NewReader(csvFile)
	_, _ = r.Read() // skip columns names

	return r
}
