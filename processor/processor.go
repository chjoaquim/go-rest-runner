package processor

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/chjoaquim/go-rest-runner/processor/strategy"
	"github.com/chjoaquim/go-rest-runner/reader"
	"github.com/chjoaquim/go-rest-runner/writer"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
	"time"
)

func Run() *strategy.Output {
	init := time.Now()
	var input reader.Input
	input.ReadFile()

	filePath := flag.Lookup("data").Value.String()
	csvFile, err := os.Open(filePath)
	if err != nil {
		log.Errorf("Error when trying to read csv file. %s", err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	resultWriter := writer.NewResultWriter()
	semaphore := make(chan struct{}, input.MaxGoroutines)
	var wg sync.WaitGroup
	line := 0

	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Errorf("Erro ao ler o arquivo CSV:", err)
			return nil
		}

		if len(row) < 1 {
			log.Errorf("Linha inválida no arquivo CSV")
			continue
		}

		log.Infoln("----")
		log.Infoln(fmt.Sprintf("Processando Linha %d", line))
		line++

		vars := mapCSVToVariables(row)

		wg.Add(1)

		go func(input reader.Input, resultWriter writer.Writer, vars map[string]interface{}, line int) {
			defer wg.Done()
			semaphore <- struct{}{}

			outputs := runStep(input, vars)
			appendOutputs(*outputs, resultWriter, line)

			<-semaphore
		}(input, resultWriter, vars, line)

		wg.Wait()
	}

	end := time.Now()

	err = resultWriter.Write("output.csv")
	if err != nil {
		log.Errorf("Errow when trying to write output file. %s", err)
	}

	log.Infoln("Time elapsed: ", end.Sub(init))
	return nil
}

func appendOutputs(outputs []strategy.Output, resultWriter writer.Writer, line int) {
	for _, output := range outputs {
		result := writer.Result{
			Line:        line,
			Status:      output.Status,
			Information: output.Message,
		}

		resultWriter.AppendResult(result)
	}
}

func runStep(input reader.Input, vars map[string]interface{}) *[]strategy.Output {
	outputs := make([]strategy.Output, 0)
	for _, s := range input.Steps {
		log.Infoln(s.Name)
		output := runRequest(s, vars)
		outputs = append(outputs, *output)
	}

	return &outputs
}

func runRequest(step reader.Step, vars map[string]interface{}) *strategy.Output {
	factory := buildFactory()
	st := factory.Find(step.Method)
	result := st.DoRequest(step.Path, step.Body, reader.ToHeadersMap(step.Headers), vars)
	return &result
}

func mapCSVToVariables(row []string) map[string]interface{} {
	vars := make(map[string]interface{})

	for i, value := range row {
		vars[fmt.Sprintf("$%d", i+1)] = value
	}

	return vars
}

func writeOutput(output strategy.Output) {

}

func buildFactory() strategy.Factory {
	factory := strategy.Factory{}
	factory.GetStrategy = strategy.NewGetStrategy()
	return factory
}
