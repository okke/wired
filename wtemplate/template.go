package wtemplate

import "bytes"

// Context provides dynamic values for temlates
//
type Context interface {
	Solve(string) string
}

// Template represents a parsed template
//
type Template interface {
	Solve(Context) string
}

type literal struct {
	text string
}

func (literal *literal) Solve(ctx Context) string {
	return literal.text
}

func newLiteral(text string) Template {
	return &literal{text: text}
}

type variable struct {
	name         string
	defaultValue string
}

func (variable *variable) Solve(ctx Context) string {
	solved := ctx.Solve(variable.name)
	if solved == "" {
		return variable.defaultValue
	}
	return solved
}

func newVariable(name string, defaultValue string) Template {
	return &variable{name: name, defaultValue: defaultValue}
}

type template struct {
	parts []Template
}

func (template *template) Solve(ctx Context) string {
	var buffer bytes.Buffer

	for _, part := range template.parts {
		buffer.WriteString(part.Solve(ctx))
	}

	return buffer.String()
}

func newTemplate(parts []Template) Template {
	return &template{parts: parts}
}
