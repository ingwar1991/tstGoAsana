package main

import (
    "os"
    "gopkg.in/yaml.v3"
)

type AsanaApp struct {
    PAT string `yaml:"asanaPAT"`
    GetRateLimit int `yaml:"asanaGetRateLimit"`
    FilePath string `yaml:"pathForJsonFiles"`
}

func newAsanaApp() (*AsanaApp, error) {
    buf, err := os.ReadFile("conf.yml") 
    if err != nil {
        return nil, err
    }

    aApp := &AsanaApp{}
    err = yaml.Unmarshal(buf, aApp)
    if err != nil {
        return nil, err
    }

   return aApp, nil
}
