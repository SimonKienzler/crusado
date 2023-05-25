package userstorytemplates

import (
	"context"
	"errors"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/workitems"
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

	userStory, err := s.WorkitemsService.CreateUserStory(ctx, userStoryTemplate.StoryTitle, userStoryTemplate.StoryDescription)
	if err != nil {
		return err
	}

	for i := range userStoryTemplate.Tasks {
		task := userStoryTemplate.Tasks[i]
		_, err = s.WorkitemsService.CreateTaskUnderneathUserStory(ctx, task.Title, task.Description, userStory)
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
