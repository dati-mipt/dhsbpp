package hierarchy

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Hierarchy struct {
	TenantToParent map[string]string
	TasksPerDay    []map[string]float64
}

func NewHierarchy(csvNodes string, csvTpd string) *Hierarchy {
	var newHierarchy = Hierarchy{readNodes(csvNodes),
		readCsv(csvTpd)}

	return &newHierarchy
}

func readNodes(csvNodes string) map[string]string {
	var r = NewCsvReader(csvNodes)

	var tenantsToParent = make(map[string]string)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		childUuid, parentsUuid := record[0], record[1]

		tenantsToParent[childUuid] = parentsUuid
	}
	return tenantsToParent
}

func readCsv(csvTpd string) []map[string]float64 {
	var r = NewCsvReader(csvTpd)
	var tasksPerDay = make([]map[string]float64, 0)

	var lastDate string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var tenantUuid, date = record[0], record[1]
		createdTasks, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			fmt.Println(err)
		}

		if date != lastDate {
			tasksPerDay = append(tasksPerDay, make(map[string]float64))
		}
		tasksPerDay[len(tasksPerDay)-1][tenantUuid] += createdTasks

		lastDate = date
	}
	return tasksPerDay
}

func NewCsvReader(filepath string) *csv.Reader {
	csvFile, err := os.Open(filepath)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	r := csv.NewReader(csvFile)
	_, _ = r.Read() // skip columns names

	return r
}
