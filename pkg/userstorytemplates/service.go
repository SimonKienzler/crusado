package userstorytemplates

import (
	"context"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/workitems"
)

// Service deals with user story templates and
// creates them using a workitems service.
type Service struct {
	WorkitemsService workitems.Service
}

func (s *Service) CreateWorkitemsFromUserStoryTemplate(ctx context.Context, userStoryTemplate config.UserStoryTemplate) error {
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
