package srgs

import (
	"fmt"
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

func (s *SISRProcessor) AppendString(str string) {
	s.SimpleProcessor.AppendString(str)
	s.AppendTag(fmt.Sprintf("scopes[scopes.length-1]['raw'] = scopes[scopes.length-1]['raw'] ? scopes[scopes.length-1]['raw'] + ' %s' : '%s';", str, str))
}

func (s *SISRProcessor) GetInstance() (string, error) {
	js := s.script

	vm := otto.New()
	_, err := vm.Run("var root;\n" + js)

	if err != nil {
		return "", err
	}

	output, err := vm.Run("root ? root.out : 'No Match Found'")
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
