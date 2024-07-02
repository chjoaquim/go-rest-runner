package reader

import (
	"encoding/csv"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

func ReadDataFile() {
	filePath := flag.Lookup("data").Value.String()
	csvFile, err := os.Open(filePath)
	if err != nil {
		log.Errorf("Error when trying to read csv file. %s", err)
	}
	defer csvFile.Close()

	csv.NewReader(csvFile)

}
