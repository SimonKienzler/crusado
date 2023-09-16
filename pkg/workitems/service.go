package workitems

import (
	"context"
	"errors"
	"fmt"

	"github.com/simonkienzler/crusado/pkg/config"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/work"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
)

const (
	UserStoryType = "User Story"
	BugType       = "Bug"
	TaskType      = "Task"
)

var (
	ErrCurrentIterationUnidentifiable = errors.New("search for current iteration returned unexpected number of results")
	ErrCouldNotGetIterations          = errors.New("could not properly get the current or all iterations")
	ErrOffsetTooFarInFuture           = errors.New("offset points to a non-existent iteration in the future")
	ErrOffsetTooFarInPast             = errors.New("offset points to a non-existent iteration in the past")
	ErrTaskWithoutParent              = errors.New("cannot create task underneath work item without parent")
	ErrCouldNotAssertLinks            = errors.New("could not assert the expected type from the workItems' Links field")
)

// addOp is a shortcut variable for the Add operation.
var addOp = webapi.OperationValues.Add

// Service deals with workitems and provides functions
// for creating user stories, bugs or tasks.
type Service struct {
	WorkitemClient workitemtracking.Client
	WorkClient     work.Client
	DryRun         bool
	ProjectConfig  config.ProjectConfig
}

// GetCurrentIteration returns a single Iteration object which represents the
// iteration that is currently in progress in the configured project.
func (s *Service) GetCurrentIteration(ctx context.Context) (*work.TeamSettingsIteration, error) {
	iterations, err := s.WorkClient.GetTeamIterations(ctx, work.GetTeamIterationsArgs{
		Project:   &s.ProjectConfig.Name,
		Timeframe: stringPointer("current"),
	})
	if iterations == nil {
		return nil, err
	}

	if len(*iterations) != 1 {
		return nil, fmt.Errorf("%w: %d iterations found", ErrCurrentIterationUnidentifiable, len(*iterations))
	}

	currentIteration := (*iterations)[0]

	return &currentIteration, nil
}

// ListIterations returns a list of all iterations (past, current and future) in
// the configured project.
func (s *Service) ListIterations(ctx context.Context) (*[]work.TeamSettingsIteration, error) {
	iterations, err := s.WorkClient.GetTeamIterations(ctx, work.GetTeamIterationsArgs{
		Project: &s.ProjectConfig.Name,
	})
	if iterations == nil {
		return nil, err
	}

	return iterations, nil
}

// GetIterationRelativeToCurrent takes an integer as offset and will return the
// iteration relative to the current one, if the iteration at the specified
// offset does exist. Will return an error if the offset is too far in the past
// or too far in the future. Using 0 as offset will return the current
// iteration, using 1 will return the next iteration. Use -1 to get the previous
// iteration and so on.
func (s *Service) GetIterationRelativeToCurrent(ctx context.Context, offset int) (*work.TeamSettingsIteration, error) {
	current, err := s.GetCurrentIteration(ctx)
	if err != nil {
		return nil, err
	}

	all, err := s.ListIterations(ctx)
	if err != nil {
		return nil, err
	}

	if current == nil || all == nil {
		return nil, ErrCouldNotGetIterations
	}

	currentIndex := 0

	for i := range *all {
		if *(*all)[i].Id == *current.Id {
			currentIndex = i
			break
		}
	}

	relativePos := currentIndex + offset

	if len(*all) <= relativePos {
		return nil, fmt.Errorf("%w: %d", ErrOffsetTooFarInFuture, offset)
	}

	if relativePos < 0 {
		return nil, fmt.Errorf("%w: %d", ErrOffsetTooFarInPast, offset)
	}

	return &(*all)[relativePos], nil
}

// Create is responsible for creating arbitrary workitems of the specified type.
func (s *Service) Create(ctx context.Context, title, description string, templateType config.TemplateType) (*workitemtracking.WorkItem, error) {
	project := s.ProjectConfig.Name
	workItemType := getWorkItemTypeForTemplateType(templateType)
	validateOnly := s.DryRun
	document := s.buildBasicWorkItemJSONPatchDocument(title, description, templateType)

	return s.WorkitemClient.CreateWorkItem(ctx, workitemtracking.CreateWorkItemArgs{
		Document:     &document,
		Project:      &project,
		Type:         &workItemType,
		ValidateOnly: &validateOnly,
	})
}

func (s *Service) CreateTaskUnderneath(ctx context.Context, title, description string, parent *workitemtracking.WorkItem) (*workitemtracking.WorkItem, error) {
	project := s.ProjectConfig.Name
	workItemType := TaskType
	validateOnly := s.DryRun
	document := s.buildBasicWorkItemJSONPatchDocument(title, description, config.TemplateTypeTask)

	if parent == nil {
		return nil, ErrTaskWithoutParent
	}

	// if we're in dry-run mode, don't specify the parent-child relationship,
	// because this would trigger an existence check on the parent. This fails
	// and the command would error.
	if !validateOnly {
		document = append(document, buildJSONPatchOperation(
			addOp, "/relations/-", workitemtracking.WorkItemRelation{
				Url: parent.Url,
				Rel: stringPointer("System.LinkTypes.Hierarchy-Reverse"),
			},
		))
	}

	return s.WorkitemClient.CreateWorkItem(ctx, workitemtracking.CreateWorkItemArgs{
		Document:     &document,
		Project:      &project,
		Type:         &workItemType,
		ValidateOnly: &validateOnly,
	})
}

// GetWorkItemHTMLRef returns the URL pointing to the Azure DevOps link that
// shows the HTML view of the passed work item. That's the URL the user would
// want to visit in a browser. Returns an error if the necessary type assertion
// on the Links field fails or the expected map key isn't present. Returns nil
// in dry-run mode.
func (s *Service) GetWorkItemHTMLRef(workItem *workitemtracking.WorkItem) (*string, error) {
	if s.DryRun {
		return nil, nil
	}

	links, ok := workItem.Links.(map[string]interface{})
	if !ok {
		return nil, ErrCouldNotAssertLinks
	}

	html, ok := links["html"].(map[string]interface{})
	if !ok {
		return nil, ErrCouldNotAssertLinks
	}

	href, ok := html["href"].(string)
	if !ok {
		return nil, ErrCouldNotAssertLinks
	}

	return &href, nil
}

func (s *Service) buildBasicWorkItemJSONPatchDocument(title, description string, templateType config.TemplateType) []webapi.JsonPatchOperation {
	fieldPathForDescription := ""
	switch templateType {
	// Bug doesn't use a description, but rather Repro Steps
	case config.TemplateTypeBug:
		fieldPathForDescription = "/fields/Microsoft.VSTS.TCM.ReproSteps"
	case config.TemplateTypeUserStory:
		fieldPathForDescription = "/fields/System.Description"
	case config.TemplateTypeTask:
		fieldPathForDescription = "/fields/System.Description"
	default:
		fieldPathForDescription = "/fields/System.Description"
	}
	return []webapi.JsonPatchOperation{
		buildJSONPatchOperation(addOp, "/fields/System.Title", title),
		buildJSONPatchOperation(addOp, fieldPathForDescription, description),
		buildJSONPatchOperation(addOp, "/fields/System.AreaPath", s.ProjectConfig.AreaPath),
		buildJSONPatchOperation(addOp, "/fields/System.IterationPath", s.ProjectConfig.IterationPath),
	}
}

func buildJSONPatchOperation(op webapi.Operation, path string, value interface{}) webapi.JsonPatchOperation {
	return webapi.JsonPatchOperation{
		Op:    &op,
		Path:  &path,
		Value: value,
	}
}

func stringPointer(s string) *string {
	return &s
}

func getWorkItemTypeForTemplateType(templateType config.TemplateType) string {
	switch templateType {
	case config.TemplateTypeUserStory:
		return UserStoryType
	case config.TemplateTypeBug:
		return BugType
	case config.TemplateTypeTask:
		return TaskType
	default:
		return ""
	}
}
