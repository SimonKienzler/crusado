package templates

import (
	"bytes"
	"errors"

	"github.com/yuin/goldmark"
)

var (
	errNoTemplateFoundForName = errors.New("no template found for name")
)

type Service struct {
	TemplatesDirectory string
}

func (s *Service) GetAll() ([]Template, error) {
	// TODO implement

	return nil, errNoTemplateFoundForName
}

func (s *Service) GetByName(name string) (*Template, error) {
	// TODO implement

	return nil, errNoTemplateFoundForName
}

func ConvertMarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
