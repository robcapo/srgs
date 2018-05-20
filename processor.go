package srgs

import (
	"github.com/robertkrimen/otto"
	"strings"
)

type Processor interface {
	AppendString(str string)
	AppendTag(body string)
	SetRoot(rootId string)
	GetInterpretation() string
	GetInstance() (string, error)
}

type SimpleProcessor struct {
	root   string
	output string
	script string
}

func (s *SimpleProcessor) AppendString(str string)      { s.output = strings.TrimSpace(s.output + " " + str) }
func (s *SimpleProcessor) AppendTag(body string)        { s.script = s.script + "\n" + body }
func (s *SimpleProcessor) SetRoot(root string)          { s.root = root }
func (s *SimpleProcessor) GetInterpretation() string    { return s.output }
func (s *SimpleProcessor) GetInstance() (string, error) { return s.script, nil }

type SISRProcessor struct {
	SimpleProcessor
}

func (s *SISRProcessor) GetInstance() (string, error) {
	js := s.script

	vm := otto.New()
	_, err := vm.Run(js)

	if err != nil {
		return "", err
	}

	output, err := vm.Run("root.out")
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
