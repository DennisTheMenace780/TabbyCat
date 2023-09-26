package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func main() {
	slc := GetBranches()
	branchItems := BuildItems(slc)
	l := ListBuilder(branchItems)

    if _, err := tea.NewProgram(Model{list: l}).Run(); err != nil {
		log.Fatal(err)
	}
}
