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
	ChildToParent map[string]string // child name --> parent name

	TasksPerDay []map[string]int64 // Element of slice representing change of tree per day.
	// Key of map is tenant name.
	// Value of map is the number
	// of added entities by this tenant per day.
}

func NewHierarchy(csvNodes string, csvTpd string) *Hierarchy {
	var newHierarchy Hierarchy

	newHierarchy.ChildToParent = readTreeNodes(csvNodes)
	newHierarchy.TasksPerDay = readTasksPerDay(csvTpd)

	return &newHierarchy
}

func readTreeNodes(csvNodes string) map[string]string {
	var r = NewCsvReader(csvNodes)

	var childToParent = make(map[string]string)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		childName, parentsName := record[0], record[1]

		childToParent[childName] = parentsName
	}
	return childToParent
}

func readTasksPerDay(csvTpd string) []map[string]int64 {
	var r = NewCsvReader(csvTpd)
	var tasksPerDay = make([]map[string]int64, 0)

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
		createdTasks, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			fmt.Println(err)
		}

		if date != lastDate {
			tasksPerDay = append(tasksPerDay, make(map[string]int64))
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
