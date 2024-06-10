package main 

import (
    "bitbucket.org/mikehouston/asana-go"
    "context"
    "golang.org/x/time/rate"
)

func getProjects(workspace *asana.Workspace, client *asana.Client, limiter *rate.Limiter, ctx context.Context) ([]*asana.Project, error) {
	var allProjects []*asana.Project
	nextPage := &asana.NextPage{}

	var projects []*asana.Project
	var err error

	for nextPage != nil {
		options := &asana.Options{
			Limit:  100,
			Offset: nextPage.Offset,
            Fields: []string{
                "gid", "resouce_type", "archived", "color", "created_at",
                "current_status", "custom_field_settings",
            },
		}

		if err := limiter.Wait(ctx); err != nil {
			return nil, err
		}

		projects, nextPage, err = workspace.Projects(client, []*asana.Options{options}...)
		if err != nil {
			return nil, err
		}

		allProjects = append(allProjects, projects...)
	}

	return allProjects, nil
}
