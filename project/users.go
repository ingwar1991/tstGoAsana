package main

import (
    "bitbucket.org/mikehouston/asana-go"
    "context"
    "golang.org/x/time/rate"
    "sync"
)

func saveUserAsJsonFile(us *asana.User, filePath string, wg *sync.WaitGroup, errChan chan *error) {
    defer wg.Done()

    err := saveAsJsonFile(us, filePath + "user_" + us.ID + ".json")
    if err != nil {
        errChan <- &err
    }
}

func getUsers(workspace *asana.Workspace, client *asana.Client, 
  limiter *rate.Limiter, ctx context.Context,
  wg *sync.WaitGroup, usChan chan []*asana.User, errChan chan *error) {
    defer wg.Done()

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
            errChan <- &err
			return
		}

		users, nextPage, err = workspace.Users(client, []*asana.Options{options}...)
		if err != nil {
			errChan <- &err
            return
		}

		allUsers = append(allUsers, users...)
	}

	usChan <- allUsers
}
