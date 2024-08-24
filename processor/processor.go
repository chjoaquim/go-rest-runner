package processor

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/chjoaquim/go-rest-runner/processor/strategy"
	"github.com/chjoaquim/go-rest-runner/reader"
	"github.com/chjoaquim/go-rest-runner/writer"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	bar *progressbar.ProgressBar
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
	lines, err := csvReader.ReadAll()
	if err != nil {
		log.Errorf("Error reading CSV file: %s\n", err)
		os.Exit(1)
	}

	lineCount := len(lines)
	bar = progressbar.Default(int64(lineCount), "Processing ...")
	csvFile.Seek(0, 0)
	csvReader = csv.NewReader(csvFile)

	resultWriter := writer.NewResultWriter()
	jobs := make(chan []string)
	var wg sync.WaitGroup
	line := 0

	for w := 1; w <= input.MaxGoroutines; w++ {
		wg.Add(1)
		go worker(input, jobs, resultWriter, &wg)
	}

	go func() {
		for {
			row, err := csvReader.Read()
			if err != nil {
				break
			}

			if len(row) < 1 {
				log.Errorf("Linha inválida no arquivo CSV")
				continue
			}
			line++
			row = append(row, fmt.Sprintf("%d", line))
			jobs <- row
		}
		close(jobs)
	}()

	wg.Wait()
	end := time.Now()

	err = resultWriter.Write(input.OutputFile)
	if err != nil {
		log.Errorf("Errow when trying to write output file. %s", err)
	}

	log.Infoln("Time elapsed: ", end.Sub(init))
	return nil
}

func worker(input reader.Input, line <-chan []string, resultWriter writer.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
	for l := range line {
		size := len(l)
		num, _ := strconv.Atoi(l[size-1])
		vars := mapCSVToVariables(l)
		outputs := runStep(input, vars)
		appendOutputs(outputs, resultWriter, num)
		bar.Add(1)
	}
}

func appendOutputs(outputs []strategy.Output, resultWriter writer.Writer, lineRef int) {
	for _, output := range outputs {
		result := writer.Result{
			Line:        lineRef,
			Status:      output.Status,
			Information: fmt.Sprintf("%s", output.Message),
		}

		resultWriter.AppendResult(result)
	}
}

func runStep(input reader.Input, vars map[string]interface{}) []strategy.Output {
	outputs := make([]strategy.Output, 0)
	for _, s := range input.Steps {
		output := runRequest(s, vars)
		outputs = append(outputs, *output)
	}

	return outputs
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

func buildFactory() strategy.Factory {
	factory := strategy.Factory{}
	factory.GetStrategy = strategy.NewGetStrategy()
	factory.PostStrategy = strategy.NewPostStrategy()
	return factory
}
