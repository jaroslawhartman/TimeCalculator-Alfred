package main

import (
	"encoding/json"
	"fmt"
)

// Structure defining output filtering JSON for Alfred
type Items struct {
	Skipknowldedge bool   `json:"skipknowldedge"`
	Items          []Item `json:"items"`
}

type Action struct {
	URL string `json:"url,omitempty"`
}

type Icon struct {
	Path string `json:"path,omitempty"`
}

type Item struct {
	Uid          string `json:"uid,omitempty"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	Arg          string `json:"arg"`
	Action       Action `json:"action,omitempty"`
	QuickLookUrl string `json:"quicklookurl,omitempty"`
	Icon         Icon   `json:"icon,omitempty"`
}

type outputItemFormat struct {
	title      string
	format     string
	formatFunc func(f string, dt datetime) string
}

func getItems(dt datetime, err error) Items {

	outputItemFormats := []outputItemFormat{
		{
			title:  "Gap",
			format: "%d days, %d hours, %d minutes and %d seconds",
			formatFunc: func(f string, dt datetime) string {
				return fmt.Sprintf(f, dt.day, dt.hour, dt.minute, dt.second)
			},
		},
		{
			title:  "In days",
			format: "%.2f days",
			formatFunc: func(f string, dt datetime) string {
				return fmt.Sprintf(f, dt.days)
			},
		},
		{
			title:  "In hours",
			format: "%.2f hours",
			formatFunc: func(f string, dt datetime) string {
				return fmt.Sprintf(f, dt.hours)
			},
		},
		{
			title:  "In minutes",
			format: "%.2f minutes",
			formatFunc: func(f string, dt datetime) string {
				return fmt.Sprintf(f, dt.minutes)
			},
		},
		{
			title:  "In seconds",
			format: "%.0f seconds",
			formatFunc: func(f string, dt datetime) string {
				return fmt.Sprintf(f, dt.seconds)
			},
		},
	}

	items := Items{
		Skipknowldedge: true,
	}

	// Skip any output if error
	if err == nil {
		for _, v := range outputItemFormats {
			item := Item{
				Title:    fmt.Sprintf("%s", v.title),
				Subtitle: v.formatFunc(v.format, dt),
				Arg:      v.formatFunc(v.format, dt),
			}

			items.Items = append(items.Items, item)
		}
	} else {
		item := Item{
			Uid:      "Error",
			Title:    "Input error!",
			Subtitle: err.Error(),
			Arg:      "error",
			Action: Action{
				URL: buymeacoffee,
			},
		}
		items.Items = append(items.Items, item)

	}

	item := Item{
		Uid:      "_XBuy a Coffee",
		Title:    "Support my work and Buy me a Coffee!",
		Subtitle: "TimeDiff workflow by Jarek Hartman",
		Arg:      "open",
		Action: Action{
			URL: buymeacoffee,
		},
		QuickLookUrl: buymeacoffee,
		Icon: Icon{
			Path: "bin/bmc_icon_black.png",
		},
	}
	items.Items = append(items.Items, item)

	return items
}

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
