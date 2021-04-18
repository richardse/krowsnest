package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Namespace_Url string
}

func loadConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		panic(err.Error())
	}

	dec := json.NewDecoder(file)
	var obj Config
	dec.Decode(&obj)
	return obj
}
