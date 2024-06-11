package main 

import (
    "os"
    "fmt"
    "bitbucket.org/mikehouston/asana-go"
    "context"
    "golang.org/x/time/rate"
	"sync"
	"time"
)

func main() {
    frequency := "5m"
    if len(os.Args[1:]) > 0 {
        if os.Args[1] == "30s" {
            frequency = "30s"
        }
    }

    aApp, err := newAsanaApp() 
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    client := asana.NewClientWithAccessToken(aApp.PAT) 
    limiter := rate.NewLimiter(rate.Every(60 * time.Second), aApp.GetRateLimit)

    STARTOVER:
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	allWorkspaces, err := client.AllWorkspaces()
	if err != nil {
        fmt.Println(err)
        os.Exit(1)
	}

    var wg sync.WaitGroup
    
    prChan := make(chan []*asana.Project)
    usChan := make(chan []*asana.User)
    errChan := make(chan *error)
   
	for _, workspace := range allWorkspaces {
        wg.Add(2)

        go getProjects(workspace, client, limiter, ctx, &wg, prChan, errChan)
        go getUsers(workspace, client, limiter, ctx, &wg, usChan, errChan)
	}

    go func(){
        wg.Wait()
        close(prChan)
        close(usChan)
        close(errChan)
    }()
    
    // wg & chan for saving to .json files
    var wgFiles sync.WaitGroup
    errChanFiles := make(chan *error)

    // code looks simpler if I'm getting all the results into arrays
    // and then save them with for stmt
    //
    // but for the case when we might have a lot of entries, it's better to send them 
    // right away into files to use less memory
    for {
        select {
        case projects, ok := <- prChan:
            if ok {
                for _, project := range projects {
                    wgFiles.Add(1)
                    go saveProjectAsJsonFile(project, aApp.FilePath, &wgFiles, errChanFiles)
                }
            } else {
                prChan = nil
            }
        case users, ok := <- usChan:
            if ok {
                for _, user := range users {
                    for wsKey, ws := range user.Workspaces {
                        for _, wsFull := range allWorkspaces {
                            if ws.ID == wsFull.ID {
                                user.Workspaces[wsKey] = wsFull
                                break
                            }
                        }
                    }

                    wgFiles.Add(1)
                    go saveUserAsJsonFile(user, aApp.FilePath, &wgFiles, errChanFiles)
                }
            } else {
                usChan = nil
            }
        case err, ok := <- errChan:
            if ok {
                fmt.Println(err)
            } else {
                errChan = nil
            }
        }

        if prChan == nil && usChan == nil && errChan == nil {
            break
        }
    }

    go func() {
        wgFiles.Wait()
        close(errChanFiles)
    }()

    for err := range errChanFiles {
        fmt.Println(err)
    }
    
    fmt.Println("Finished")

    if frequency == "5m" {
        fmt.Println("Sleeping for 5m")
        time.Sleep(5 * time.Minute)
    } else {
        fmt.Println("Sleeping for 30s")
        time.Sleep(30 * time.Second)
    }

    goto STARTOVER
}
