package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const buymeacoffee = "https://www.buymeacoffee.com/jhartman"

func getAlfredJson(p string) string {
	var items Items

	dt, err := parse(p)
	items = getItems(dt, err)

	b, err := json.MarshalIndent(items, "", "  ")
	if err == nil {
		return string(b)
	} else {
		outputWithError := `
		{
		"skipknowldedge": true,
		"items": [
			{
				"uid": "Error",
				"title": "JSON Marshalling error",
				"subtitle": "%s",
			},
			]
		}
		`
		return fmt.Sprintf(outputWithError, err)
	}
}

func main() {
	var input string

	if len(os.Args) != 2 {
		// No parameters
		input = ""
	} else {
		input = os.Args[1]
	}

	fmt.Println(getAlfredJson(input))
}
