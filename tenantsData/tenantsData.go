package tenantsData

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type TenantsData struct {
	TenantToParent map[string]string
	Tpd            []TasksPerDay
}

type TasksPerDay struct {
	Tenant string
	Day string
	CreatedTasks float64
}

func NewTenantsData( csvNodes string, csvTpd string) TenantsData  {
	newTd := TenantsData{readNodes(csvNodes),
		                         readCsv( csvTpd)  }

	return newTd
}

func readNodes ( csvNodes string ) map[string]string {
	r := NewCsvReader( csvNodes )

	tenantsToParent := make(map[string]string)
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

func readCsv( csvTpd string ) []TasksPerDay {
	r := NewCsvReader( csvTpd )
	tpd := make([]TasksPerDay, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		tenantUuid, date := record[0], record[1]
		createdTasks, _ := strconv.ParseFloat(record[2], 64)
		tpd = append(tpd, TasksPerDay{tenantUuid, date,
			createdTasks})
	}
	return tpd
}

func NewCsvReader( filepath string ) *csv.Reader {
	csvFile, err := os.Open(filepath)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	r := csv.NewReader(csvFile)
	_, _ = r.Read() // skip columns names

	return r
}
