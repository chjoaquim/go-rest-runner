package strategy

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type PostStrategy struct {
	client *http.Client
}

func NewPostStrategy() PostStrategy {
	return PostStrategy{
		client: &http.Client{},
	}
}

func (g PostStrategy) DoRequest(path string, body string, headers map[string]interface{}, vars map[string]interface{}) Output {
	bodyReader, err := StructToReader(body)
	if err != nil {
		log.Error("Error when trying to create body reader: %s", err)
		return Output{
			Succeeded: false,
			Message:   err.Error(),
			Status:    "Request Error",
		}
	}
	req, err := http.NewRequest(http.MethodPost, replaceVars(path, vars), bodyReader)
	if err != nil {
		log.Error("Error when trying to create a REQUEST: %s", err)
		return Output{
			Succeeded: false,
			Message:   err.Error(),
			Status:    "Error creating request",
		}
	}
	setHeaders(headers, req)
	resp, err := g.client.Do(req)

	if err != nil {
		log.Error("Error when trying to Post: %s", err)
		return Output{
			Succeeded: false,
			Message:   err.Error(),
			Status:    "Request Error",
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error to read response body: %v", err)
		return Output{
			Succeeded: false,
			Message:   string(bodyBytes),
			Status:    "Read body Error",
		}
	}

	return Output{
		Succeeded: true,
		Message:   string(bodyBytes),
		Status:    http.StatusText(resp.StatusCode),
	}
}

func StructToReader(s interface{}) (io.Reader, error) {
	marched, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(marched), nil
}
