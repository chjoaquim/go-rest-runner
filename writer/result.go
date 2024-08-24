package writer

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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

	ResponseResult struct {
		Package `json:"package"`
	}

	Package struct {
		Items []interface{} `json:"items"`
	}

	//ShipmentResult struct {
	//	Category           int64  `json:"category"`
	//	CategorySat        string `json:"category_sat"`
	//	Description        string `json:"description"`
	//	UnitCode           string `json:"unit_code"`
	//	DangerousMaterial  string `json:"dangerous_material"`
	//	Packagekey         string `json:"package_key"`
	//	PackageDescription string `json:"package_description"`
	//}
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
		response := ResponseResult{}
		err := json.Unmarshal([]byte(result.Information), &response)
		if err != nil {
			return err
		}

		for _, item := range response.Package.Items {
			value := reflect.ValueOf(item)
			category := value.FieldByName("category")
			description := value.FieldByName("description")

			_, err = file.WriteString(fmt.Sprintf("%d;%s;%s;%s;\n", result.Line, result.Status, category, description))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ORIGINAL:
//func (rw *ResultWriter) Write(filePath string) error {
//	file, err := os.Create(filePath)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	for _, result := range *rw.Results {
//		_, err := file.WriteString(fmt.Sprintf("%d;%s;%s", result.Line, result.Status, result.Information))
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
