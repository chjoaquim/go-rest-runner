package main

import (
	"flag"
	"github.com/chjoaquim/go-rest-runner/processor"
	"github.com/chjoaquim/go-rest-runner/reader"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("Hello Rest Runner ! ... ")
	configureFlags()

	var input reader.Input
	input.ReadFile()

	for _, s := range input.Steps {
		log.Info(s.Name)
		output := processor.Run(s)
		log.Infoln("Output: ", output)
	}
}

func configureFlags() {
	flag.String("file", "file.yaml", "Data file full path")
	flag.Parse()
}
