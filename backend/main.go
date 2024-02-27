package main

import (
	"fmt"
	"os"
)

const buymeacoffee = "https://www.buymeacoffee.com/jhartman"

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
