package main

import (
    "encoding/json"
    "os"
)

func toJson(entity interface{}) ([]byte, error) {
	res, err := json.MarshalIndent(entity, "", "  ")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func saveAsJsonFile(entity interface{}, filePath string) error {
    jsonData, err := toJson(entity)
    if err != nil {
        return err
    }

    err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

    return nil
}
