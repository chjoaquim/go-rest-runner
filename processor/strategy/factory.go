package strategy

type RequestStrategy interface {
	DoRequest(string string, body string, headers map[string]interface{}) Output
}

type Factory struct {
	GetStrategy
}

type Output struct {
	Succeeded  bool
	Message    string
	StatusCode int
}

func (sf Factory) Find(requestType string) RequestStrategy {
	switch requestType {
	case "GET":
		return sf.GetStrategy
	default:
		return nil
	}
}
