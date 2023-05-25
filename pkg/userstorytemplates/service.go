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
	Profile          config.Profile
}

func (s *Service) CreateWorkitemsFromUserStoryTemplate(ctx context.Context, userStoryTemplateName string) error {
	userStoryTemplate, err := s.GetUserStoryTemplateFromName(userStoryTemplateName)
	if err != nil {
		return err
	}

	storyDescriptionHTML, err := convertMarkdownToHTML(userStoryTemplate.StoryDescription)
	if err != nil {
		return err
	}

	userStory, err := s.WorkitemsService.CreateUserStory(ctx, userStoryTemplate.StoryTitle, storyDescriptionHTML)
	if err != nil {
		return err
	}

	for i := range userStoryTemplate.Tasks {
		task := userStoryTemplate.Tasks[i]
		taskDescriptionHTML, err := convertMarkdownToHTML(userStoryTemplate.StoryDescription)
		if err != nil {
			return err
		}

		_, err = s.WorkitemsService.CreateTaskUnderneathUserStory(ctx, task.Title, taskDescriptionHTML, userStory)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetUserStoryTemplateFromName(userStoryTemplateName string) (*config.UserStoryTemplate, error) {
	for i := range s.Profile.Templates {
		template := s.Profile.Templates[i]
		if template.Name == userStoryTemplateName {
			return &template, nil
		}
	}

	return nil, errNoUserStoryTemplateFoundForName
}

func convertMarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
