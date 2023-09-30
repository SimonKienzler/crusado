package workitems

import (
	"context"
	"errors"

	"github.com/simonkienzler/crusado/pkg/templates"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
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
	DryRun         bool

	ProjectName   string
	AreaPath      string
	IterationPath string
}

// Create is responsible for creating arbitrary workitems of the specified type.
func (s *Service) Create(ctx context.Context, title, description string, templateType templates.Type) (*workitemtracking.WorkItem, error) {
	project := s.ProjectName
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
	project := s.ProjectName
	workItemType := TaskType
	validateOnly := s.DryRun
	document := s.buildBasicWorkItemJSONPatchDocument(title, description, templates.TaskType)

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

func (s *Service) buildBasicWorkItemJSONPatchDocument(title, description string, templateType templates.Type) []webapi.JsonPatchOperation {
	fieldPathForDescription := ""
	switch templateType {
	// Bug doesn't use a description, but rather Repro Steps
	case templates.BugType:
		fieldPathForDescription = "/fields/Microsoft.VSTS.TCM.ReproSteps"
	case templates.UserStoryType:
		fieldPathForDescription = "/fields/System.Description"
	case templates.TaskType:
		fieldPathForDescription = "/fields/System.Description"
	default:
		fieldPathForDescription = "/fields/System.Description"
	}
	return []webapi.JsonPatchOperation{
		buildJSONPatchOperation(addOp, "/fields/System.Title", title),
		buildJSONPatchOperation(addOp, fieldPathForDescription, description),
		buildJSONPatchOperation(addOp, "/fields/System.AreaPath", s.AreaPath),
		buildJSONPatchOperation(addOp, "/fields/System.IterationPath", s.IterationPath),
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

func getWorkItemTypeForTemplateType(templateType templates.Type) string {
	switch templateType {
	case templates.UserStoryType:
		return UserStoryType
	case templates.BugType:
		return BugType
	case templates.TaskType:
		return TaskType
	default:
		return ""
	}
}
