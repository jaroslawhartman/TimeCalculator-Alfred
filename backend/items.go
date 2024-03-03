package main

import (
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
	formatFunc func(dt datetime) string
}

func getItems(dt datetime, err error) Items {

	outputItemFormatsDuration := []outputItemFormat{
		{
			title: "Result",
			formatFunc: func(dt datetime) string {
				format := "%d days, %d hours, %d minutes and %d seconds"
				return fmt.Sprintf(format, dt.day, dt.hour, dt.minute, dt.second)
			},
		},
		{
			title: "Result (hh:mm:ss)",
			formatFunc: func(dt datetime) string {
				if dt.day == 0 {
					format := "%02d:%02d:%02d"
					return fmt.Sprintf(format, dt.hour, dt.minute, dt.second)
				} else {
					format := "%dd, %02d:%02d:%02d"
					return fmt.Sprintf(format, dt.day, dt.hour, dt.minute, dt.second)
				}
			},
		},
		{
			title: "In days",
			formatFunc: func(dt datetime) string {
				format := "%.2f days"
				return fmt.Sprintf(format, dt.days)
			},
		},
		{
			title: "In hours",
			formatFunc: func(dt datetime) string {
				format := "%.2f hours"
				return fmt.Sprintf(format, dt.hours)
			},
		},
		{
			title: "In minutes",
			formatFunc: func(dt datetime) string {
				format := "%.2f minutes"
				return fmt.Sprintf(format, dt.minutes)
			},
		},
		{
			title: "In seconds",
			formatFunc: func(dt datetime) string {
				format := "%.0f seconds"
				return fmt.Sprintf(format, dt.seconds)
			},
		},
	}

	outputItemFormatsNumber := []outputItemFormat{
		{
			title: "Result",
			formatFunc: func(dt datetime) string {
				format := "%d"
				return fmt.Sprintf(format, dt.ts)
			},
		},
	}

	outputItemFormatsTimestamp := []outputItemFormat{
		{
			title: "Result",
			formatFunc: func(dt datetime) string {
				format := "%v"
				return fmt.Sprintf(format, dt.dt)
			},
		},
	}

	items := Items{
		Skipknowldedge: true,
	}

	// Skip any output if error
	if err == nil {
		if dt.kind == number {
			for _, v := range outputItemFormatsNumber {
				item := Item{
					Title:    fmt.Sprintf("%s", v.title),
					Subtitle: v.formatFunc(dt),
					Arg:      v.formatFunc(dt),
				}
				items.Items = append(items.Items, item)
			}
		} else if dt.kind == duration {

			for _, v := range outputItemFormatsDuration {
				item := Item{
					Title:    fmt.Sprintf("%s", v.title),
					Subtitle: v.formatFunc(dt),
					Arg:      v.formatFunc(dt),
				}
				items.Items = append(items.Items, item)
			}
		} else if dt.kind == timestamp {
			for _, v := range outputItemFormatsTimestamp {
				item := Item{
					Title:    fmt.Sprintf("%s", v.title),
					Subtitle: v.formatFunc(dt),
					Arg:      v.formatFunc(dt),
				}
				items.Items = append(items.Items, item)
			}
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
