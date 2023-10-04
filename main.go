package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
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

    fmt.Println(status)

	if _, err := tea.NewProgram(Model{list: l, repo: repo}).Run(); err != nil {
		log.Fatal(err)
	}
	// if !status.IsClean() {
	// 	fmt.Println("Uncommited Changes")
	// 	tea.Quit()
	// } else {
	// 	if _, err := tea.NewProgram(Model{list: l, repo: repo}).Run(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
}
