package main

import (
	"dbctui/dbc"
	"dbctui/ui"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	args := os.Args[1:]
	dbcFile := args[0]
	data, err := os.ReadFile(dbcFile)
	if err != nil {
		fmt.Println(err)
	} else {
		messages, signals := dbc.Parse(string(data))

		p := tea.NewProgram(ui.InitialModel(messages, signals))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	}
}
