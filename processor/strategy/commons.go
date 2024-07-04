package strategy

import (
	"fmt"
	"net/http"
	"strings"
)

func replaceVars(value string, vars map[string]interface{}) string {
	for key, val := range vars {
		value = strings.ReplaceAll(value, key, fmt.Sprintf("%v", val))
	}

	return value
}

func setHeaders(headers map[string]interface{}, req *http.Request) {
	for key, value := range headers {
		req.Header.Set(key, value.(string))
	}
}
