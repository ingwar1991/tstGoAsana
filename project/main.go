package main 

import (
    "os"
    "fmt"
    "bitbucket.org/mikehouston/asana-go"
    "context"
    "golang.org/x/time/rate"
//	"sync"
	"time"
)

func main() {
    aApp, err := newAsanaApp() 
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    client := asana.NewClientWithAccessToken(aApp.PAT) 
    limiter := rate.NewLimiter(rate.Every(60 * time.Second), aApp.GetRateLimit)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

//	var wg sync.WaitGroup

	workspaces, err := client.AllWorkspaces()
	if err != nil {
        fmt.Println(err)
        os.Exit(1)
	}
   
    allProjects := []*asana.Project{}
    allUsers := []*asana.User{}
	for _, workspace := range workspaces {
//        wg.Add(1)
        projects, err := getProjects(workspace, client, limiter, ctx)
        if err != nil {
            fmt.Println("Pr error:", err)
        } else {
            allProjects = append(allProjects, projects...)
        }

//        wg.Add(1)
        users, err := getUsers(workspace, client, limiter, ctx)
        if err != nil {
            fmt.Println("Us error:", err)
        } else {
            allUsers = append(allUsers, users...)
        }
	}

    fmt.Println("All projects:", allProjects)
    fmt.Println("All users:", allUsers)

//	wg.Wait()
//	fmt.Println("All API calls completed")
}
