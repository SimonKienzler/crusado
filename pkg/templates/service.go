package templates

import (
	"bytes"
	"errors"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/workitems"

	"github.com/yuin/goldmark"
)

var (
	errNoTemplateFoundForName = errors.New("No template found for name")
)

// Service deals with templates and creates them using a workitems service.
type Service struct {
	WorkitemsService workitems.Service
	TemplateList     config.TemplateList
}

func (s *Service) GetTemplateFromName(templateName string) (*config.Template, error) {
	for i := range s.TemplateList.Templates {
		template := s.TemplateList.Templates[i]
		if template.Name == templateName {
			return &template, nil
		}
	}

	return nil, errNoTemplateFoundForName
}

func ConvertMarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
