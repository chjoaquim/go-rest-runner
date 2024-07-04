package main

import (
	"flag"
	"github.com/chjoaquim/go-rest-runner/processor"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("Hello Rest Runner ! ... ")
	configureFlags()
	processor.Run()
}

func configureFlags() {
	flag.String("file", "file.yaml", "Steps file full path")
	flag.String("data", "data.cvs", "Data file full path")
	flag.Parse()
}
