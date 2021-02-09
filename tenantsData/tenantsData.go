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
	CreatedTasks int
}

func NewTenantsData( csvNodes string, csvTpd string) TenantsData  {
	new_td := TenantsData{ readNodes(csvNodes),
		                         readCsv( csvTpd)  }

	return new_td
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
		child_uuid, parents_uuid := record[0], record[1]

		tenantsToParent[child_uuid] = parents_uuid
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
		tenant_uuid, date := record[0], record[1]
		created_tasks, _ := strconv.Atoi(record[2])
		tpd = append(tpd, TasksPerDay{ tenant_uuid, date,
										created_tasks})
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
