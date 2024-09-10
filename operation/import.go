package operation

import (
	"fmt"
	"encoding/csv"
	"os"
)

type ImportCSV struct {
	CsvPath string
}

type CSV struct {
	SourceOrg []string
	SourceProject []string
	TargetOrg []string
	TargetProject []string
}


func (m ImportCSV) Exec() (*CSV, error) {
	// Read the file
	file, err := os.Open(m.CsvPath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Parse the file
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error readying CSV: %v", err)
	}

	// Initialize the CSV struct to store the data
	parsedCSV := &CSV{
		SourceOrg:     []string{},
		SourceProject: []string{},
		TargetOrg:     []string{},
		TargetProject: []string{},
	}

	// Loops through each line of the CSV and adds to the parsedCSV struct
	for i, row := range records {
		if i == 0 {
			continue
			// Skip header row
		}

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid CSV format. Expected 4 columns, but got %d", len(row))
		}

		parsedCSV.SourceOrg = append(parsedCSV.SourceOrg, row[0])
		parsedCSV.SourceProject = append(parsedCSV.SourceProject, row[1])
		parsedCSV.TargetOrg = append(parsedCSV.TargetOrg, row[2])
		parsedCSV.TargetProject = append(parsedCSV.TargetProject, row[3])
	}

	return parsedCSV, nil
}