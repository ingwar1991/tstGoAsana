package main

import (
    "bitbucket.org/mikehouston/asana-go"
    "context"
    "golang.org/x/time/rate"
)

func getUsers(workspace *asana.Workspace, client *asana.Client, limiter *rate.Limiter, ctx context.Context) ([]*asana.User, error) {
    var allUsers []*asana.User
	nextPage := &asana.NextPage{}

	var users []*asana.User
	var err error

	for nextPage != nil {
		options := &asana.Options{
			Limit:  50,
			Offset: nextPage.Offset,
            Fields: []string{
                "gid", "name", "email", "photo", "resource_type", "workspaces",
            },
		}

		if err := limiter.Wait(ctx); err != nil {
			return nil, err
		}

		users, nextPage, err = workspace.Users(client, []*asana.Options{options}...)
		if err != nil {
			return nil, err
		}

		allUsers = append(allUsers, users...)
	}

	return allUsers, nil
}
