package srgs

import "strings"

type Processor interface {
	AppendString(str string)
	AppendTag(body string)
	GetString() (string, error)
}

type SimpleProcessor struct {
	output string
}

func (s *SimpleProcessor) AppendString(str string)    { s.output = strings.TrimSpace(s.output + " " + str) }
func (s *SimpleProcessor) AppendTag(body string)      { s.output = strings.TrimSpace(s.output + " " + body) }
func (s *SimpleProcessor) GetString() (string, error) { return s.output, nil }
