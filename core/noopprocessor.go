package core

import (
	"errors"
	"fmt"
	"os"
)

// NoOpProcessor is a Processor that literally spits out the contents of the file itself
// with no transformations besides wrapping it in div/pre block
type NoOpProcessor struct{}

var _ Processor = (*NoOpProcessor)(nil)

func (NoOpProcessor) CanHandle(fileLocation string) bool {
	return true
}

func (n NoOpProcessor) Handle(fileLocation string, outputFile string) error {
	if !n.CanHandle(fileLocation) {
		return errors.New("Cannot handle file: " + fileLocation)
	}

	data, err := os.ReadFile(fileLocation)
	if err != nil {
		return err
	}

	htmlified := fmt.Sprintf(
		"<div><pre style=\"font-family: ui-monospace, monospace; margin:0;\">%s</pre></div>",
		string(data))

	err = os.WriteFile(outputFile, []byte(htmlified), 0o644)
	if err != nil {
		return err
	}

	return nil
}
