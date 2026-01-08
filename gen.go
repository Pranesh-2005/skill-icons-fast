package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	iconsDir, err := os.ReadDir("./assets")
	if err != nil {
		panic(err)
	}

	var b strings.Builder

	b.WriteString("package handler\n\n")
	b.WriteString("var icons = map[string]string{\n")

	for _, file := range iconsDir {
		if file.IsDir() {
			continue
		}

		data, err := os.ReadFile("./assets/" + file.Name())
		if err != nil {
			panic(err)
		}

		name := strings.TrimSuffix(strings.ToLower(file.Name()), ".svg")

		b.WriteString(fmt.Sprintf(
			"\t%q: `%s`,\n",
			name,
			data,
		))
	}

	b.WriteString("}\n")

	err = os.WriteFile("./api/icons_data.go", []byte(b.String()), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("icons_data.go generated")
}
