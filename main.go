package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-git/go-git/v5"
)

func main() {
	slc := GetBranches()
	branchItems := BuildItems(slc)
	l := ListBuilder(branchItems)

	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}

	w, err := repo.Worktree()
	if err != nil {
		log.Print(err)
	}

	status, err := w.Status()
	if err != nil {
		log.Print(nil)
	}

	var stylez = lipgloss.NewStyle().PaddingLeft(8).Foreground(lipgloss.Color("213"))

	if !status.IsClean() {
		fmt.Println(
			"error: Your local changes to the following files would be overwritten by checkout:",
		)
		for k := range status {
			fmt.Println(stylez.Render(k))
		}
		fmt.Println("Please commit your changes or stash them before you switch branches.")
		fmt.Println("Aborting")
	} else {
		if _, err := tea.NewProgram(Model{list: l, repo: repo}).Run(); err != nil {
			log.Fatal(err)
		}
	}
}
