package workitems

import (
	"context"
	"errors"
	"fmt"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/work"
)

var (
	ErrIterationPathNotSet = errors.New("iteration path could not be retrieved")
)

func GetIterationPathFromOffset(ctx context.Context, client work.Client, project string, offset int) (string, error) {
	iteration, err := getIterationRelativeToCurrent(ctx, client, project, offset)
	if err != nil {
		return "", err
	}

	if iteration.Path == nil || *iteration.Path == "" {
		return "", ErrIterationPathNotSet
	}

	return *iteration.Path, nil
}

// getCurrentIteration returns a single Iteration object which represents the
// iteration that is currently in progress in the configured project.
func getCurrentIteration(ctx context.Context, client work.Client, project string) (*work.TeamSettingsIteration, error) {
	iterations, err := client.GetTeamIterations(ctx, work.GetTeamIterationsArgs{
		Project:   &project,
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

// listIterations returns a list of all iterations (past, current and future) in
// the configured project.
func listIterations(ctx context.Context, client work.Client, project string) (*[]work.TeamSettingsIteration, error) {
	iterations, err := client.GetTeamIterations(ctx, work.GetTeamIterationsArgs{
		Project: &project,
	})
	if iterations == nil {
		return nil, err
	}

	return iterations, nil
}

// getIterationRelativeToCurrent takes an integer as offset and will return the
// iteration relative to the current one, if the iteration at the specified
// offset does exist. Will return an error if the offset is too far in the past
// or too far in the future. Using 0 as offset will return the current
// iteration, using 1 will return the next iteration. Use -1 to get the previous
// iteration and so on.
func getIterationRelativeToCurrent(ctx context.Context, client work.Client, project string, offset int) (*work.TeamSettingsIteration, error) {
	current, err := getCurrentIteration(ctx, client, project)
	if err != nil {
		return nil, err
	}

	all, err := listIterations(ctx, client, project)
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
