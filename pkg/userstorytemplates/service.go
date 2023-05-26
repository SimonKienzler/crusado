package userstorytemplates

import (
	"bytes"
	"context"
	"errors"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/workitems"

	"github.com/yuin/goldmark"
)

var (
	errNoUserStoryTemplateFoundForName = errors.New("No user story template found for name")
)

// Service deals with user story templates and
// creates them using a workitems service.
type Service struct {
	WorkitemsService workitems.Service
	TemplateList     config.TemplateList
}

func (s *Service) CreateWorkitemsFromUserStoryTemplate(ctx context.Context, userStoryTemplateName string) error {

	return nil
}

func (s *Service) GetUserStoryTemplateFromName(userStoryTemplateName string) (*config.UserStoryTemplate, error) {
	for i := range s.TemplateList.Templates {
		template := s.TemplateList.Templates[i]
		if template.Name == userStoryTemplateName {
			return &template, nil
		}
	}

	return nil, errNoUserStoryTemplateFoundForName
}

func ConvertMarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
