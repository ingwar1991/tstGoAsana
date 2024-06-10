package main

import (
    "io/ioutil"
    "gopkg.in/yaml.v3"
)

type AsanaApp struct {
    PAT string `yaml:"asanaPAT"`
    GetRateLimit int `yaml:"asanaGetRateLimit"`
}

func newAsanaApp() (*AsanaApp, error) {
    buf, err := ioutil.ReadFile("conf.yml") 
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
