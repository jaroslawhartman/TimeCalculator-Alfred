package main

import (
	"fmt"
	"os"
	"strconv"
)

// TODO
// - JSON marshalling/unmarshalling

var output_header = `
{
	"skipknowldedge": true,
	"items": [
`

var output_item = `
{
	"uid": "Zzz Adv",
	"title": "Days, hours, minutes, and seconds",
	"subtitle": "%d days, 4 hours, <m> minutes, and <s> seconds",
	"arg": "%d days, 4 hours, <m> minutes, and <s> seconds",
},
`

var output_adv = `
{
	"uid": "Buy a Coffee",
	"title": "Support my work and Buy me a Coffee!",
	"subtitle": "TimeDiff workflow by Jarek Hartman",

	"action": {
		"url": "https://www.buymeacoffee.com/jhartman",
	},
	"arg": "open",
	"quicklookurl": "https://www.buymeacoffee.com/jhartman",
	"icon": {
		"path": "bin/bmc_icon_black.png"
	}
},
`

var output_footer = `
]}
`

func main() {
	var input string
	var output string

	if len(os.Args) != 2 {
		// No parameters
		input = ""
	} else {
		input = os.Args[1]
	}

	output = output_header

	if number, err := strconv.Atoi(input); err == nil {
		for i := 0; i < number; i++ {
			output += fmt.Sprintf(output_item, i, i)
		}
	}

	output += output_adv
	output += output_footer

	fmt.Println(output)
}
