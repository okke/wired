package wtemplate

import (
	"bytes"
	"strings"
	"text/scanner"
)

// Parser parses a tempate
//
type Parser interface {
	Parse(string) Template
}

type parser struct {
}

const tokDOLLAR = '$'
const tokLBRACE = '{'
const tokRBRACE = '}'
const tokCOLON = ':'
const tokBACKSLASH = '\\'
const tokVARSEPARATORS = " \n\t/\\:"

func (parser *parser) readIntoBuffer(s *scanner.Scanner, filler func(*bytes.Buffer)) string {
	var buffer bytes.Buffer

	filler(&buffer)

	return buffer.String()
}

func (parser *parser) parseLiteral(s *scanner.Scanner) (string, bool) {

	if s.Peek() == tokDOLLAR || s.Peek() == scanner.EOF {
		return "", false
	}

	return parser.readIntoBuffer(s, func(buffer *bytes.Buffer) {
		for r := s.Peek(); r != tokDOLLAR && r != scanner.EOF; r = s.Peek() {
			if r == tokBACKSLASH {
				s.Next()
				if r = s.Peek(); r != scanner.EOF {
					buffer.WriteRune(s.Next())
				}
			} else {
				buffer.WriteRune(s.Next())
			}
		}
	}), true

}

func (parser *parser) parseDefaultLiteralForVariable(s *scanner.Scanner) string {

	return parser.readIntoBuffer(s, func(buffer *bytes.Buffer) {
		for r := s.Peek(); r != tokRBRACE && r != scanner.EOF; r = s.Peek() {
			buffer.WriteRune(s.Next())
		}
	})

}

func (parser *parser) parseBracedVariable(s *scanner.Scanner) (string, string, bool) {

	defaultValue := ""

	if s.Peek() != tokLBRACE {
		return "", "", false
	}

	s.Next()

	return parser.readIntoBuffer(s, func(buffer *bytes.Buffer) {

		for r := s.Next(); r != tokRBRACE && r != scanner.EOF; r = s.Next() {
			if r == tokCOLON {
				defaultValue = parser.parseDefaultLiteralForVariable(s)
			} else {
				buffer.WriteRune(r)
			}
		}

	}), defaultValue, true
}

func (parser *parser) parseUnBracedVariable(s *scanner.Scanner) (string, string, bool) {

	result := parser.readIntoBuffer(s, func(buffer *bytes.Buffer) {
		for r := s.Peek(); !strings.ContainsRune(tokVARSEPARATORS, r) && r != scanner.EOF; r = s.Peek() {
			buffer.WriteRune(s.Next())
		}
	})

	return result, "", result != ""
}

func (parser *parser) parseVariable(s *scanner.Scanner) (string, string, bool) {
	if s.Peek() != tokDOLLAR {
		return "", "", false
	}
	s.Next()
	if text, defaultValue, parsed := parser.parseBracedVariable(s); parsed {
		return text, defaultValue, parsed
	}
	if text, defaultValue, parsed := parser.parseUnBracedVariable(s); parsed {
		return text, defaultValue, parsed
	}
	return "", "", false
}

func (parser *parser) Parse(template string) Template {
	var s scanner.Scanner
	s.Init(strings.NewReader(template))

	parts := make([]Template, 0, 0)

	for {
		if text, parsed := parser.parseLiteral(&s); parsed {
			parts = append(parts, newLiteral(text))
			continue
		}

		if text, defaultValue, parsed := parser.parseVariable(&s); parsed {
			parts = append(parts, newVariable(text, defaultValue))
			continue
		}

		return newTemplate(parts)
	}

}

// NewParser constructs a parser
//
func NewParser() Parser {
	return &parser{}
}

// Parse parses a template and evaluate it within a given context
//
func Parse(ctx Context, template string) string {
	return NewParser().Parse(template).Solve(ctx)
}
