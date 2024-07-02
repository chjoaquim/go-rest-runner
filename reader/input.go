package reader

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type Input struct {
	Name          string `yaml:"name"`
	MaxGoroutines int    `yaml:"max-goroutines"`
	Steps         []Step `yaml:"steps"`
}

type Step struct {
	Name    string   `yaml:"name"`
	Path    string   `yaml:"path"`
	Method  string   `yaml:"method"`
	Body    string   `yaml:"body,omitempty"`
	Headers []Header `yaml:"headers,omitempty"`
}

type Header struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func (i *Input) ReadFile() {
	filePath := flag.Lookup("file").Value.String()
	ymlFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Error("Error when trying to read yaml file. %s", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(ymlFile, &i)
	if err != nil {
		log.Error("Error when trying to unmarshal yaml file.", err)
	}
}

func ToHeadersMap(headers []Header) map[string]interface{} {
	headersMap := make(map[string]interface{})
	for _, h := range headers {
		headersMap[h.Name] = h.Value
	}
	return headersMap
}
