package main 

import (
    "bitbucket.org/mikehouston/asana-go"
    "context"
    "golang.org/x/time/rate"
    "sync"
)

func saveProjectAsJsonFile(pr *asana.Project, filePath string, wg *sync.WaitGroup, errChan chan *error) {
    defer wg.Done()

    err := saveAsJsonFile(pr, filePath + "project_" + pr.ID + ".json")
    if err != nil {
        errChan <- &err
    }
}

func getProjects(workspace *asana.Workspace, client *asana.Client, 
  limiter *rate.Limiter, ctx context.Context,
  wg *sync.WaitGroup, prChan chan []*asana.Project, errChan chan *error) {
    defer wg.Done()

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
            errChan <- &err
			return
		}

		projects, nextPage, err = workspace.Projects(client, []*asana.Options{options}...)
		if err != nil {
			errChan <- &err
            return
		}

		allProjects = append(allProjects, projects...)
	}

	prChan <- allProjects
}
