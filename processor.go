package main

import "strings"

type Processor interface {
	AppendString(str string)

	AppendTag(output string) error

	GetString() string
}

type SimpleProcessor struct {
	output string
}

func (s *SimpleProcessor) AppendString(str string) { s.output = strings.TrimSpace(s.output + " " + str) }

func (s *SimpleProcessor) AppendTag(body string) error {
	s.output += strings.TrimSpace(s.output + " " + strings.TrimSpace(body))

	return nil
}

func (s *SimpleProcessor) GetString() string { return s.output }
