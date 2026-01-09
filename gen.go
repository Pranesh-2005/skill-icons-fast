//go:build tools
// +build tools

package main

import (
	"encoding/json"
	"os"
	"strings"
)

func main() {
	files, err := os.ReadDir("./assets")
	if err != nil {
		panic(err)
	}

	icons := map[string]string{}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		data, err := os.ReadFile("./assets/" + f.Name())
		if err != nil {
			panic(err)
		}

		name := strings.TrimSuffix(strings.ToLower(f.Name()), ".svg")
		icons[name] = string(data)
	}

	out, err := json.Marshal(icons)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("icons.json", out, 0644)
	if err != nil {
		panic(err)
	}

	println("icons.json generated")
}
