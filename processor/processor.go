package processor

import (
	"github.com/chjoaquim/go-rest-runner/processor/strategy"
	"github.com/chjoaquim/go-rest-runner/reader"
	log "github.com/sirupsen/logrus"
)

func Run(step reader.Step) strategy.Output {
	log.Info("Running step: %v", step.Name)
	factory := buildFactory()
	st := factory.Find(step.Method)
	output := st.DoRequest(step.Path, step.Body, reader.ToHeadersMap(step.Headers))

	return output
}

func buildFactory() strategy.Factory {
	factory := strategy.Factory{}
	factory.GetStrategy = strategy.NewGetStrategy()
	return factory
}
