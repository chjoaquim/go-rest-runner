package writer

import (
	"fmt"
	"os"
)

type (
	Result struct {
		Line        int
		Status      string
		Information string
	}

	Writer interface {
		Write(filePath string) error
		AppendResult(result Result)
	}

	ResultWriter struct {
		Results *[]Result
	}
)

func NewResultWriter() *ResultWriter {
	results := make([]Result, 0)
	return &ResultWriter{
		Results: &results,
	}
}

func (rw *ResultWriter) AppendResult(result Result) {
	*rw.Results = append(*rw.Results, result)
}

func (rw *ResultWriter) Write(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, result := range *rw.Results {
		_, err := file.WriteString(fmt.Sprintf("%d;%s;%s\n", result.Line, result.Status, result.Information))
		if err != nil {
			return err
		}
	}

	return nil
}
