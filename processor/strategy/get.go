package strategy

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type GetStrategy struct {
	client *http.Client
}

func NewGetStrategy() GetStrategy {
	return GetStrategy{
		client: &http.Client{},
	}
}

func (g GetStrategy) DoRequest(path string, body string, headers map[string]interface{}, vars map[string]interface{}) Output {
	log.Info("Doing GET request")
	req, err := http.NewRequest(http.MethodGet, replaceVars(path, vars), nil)
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
		log.Error("Error when trying to GET: %s", err)
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
