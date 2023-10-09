package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Model struct {
	list     list.Model
	repo     *git.Repository
	choice   string
	quitting bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	w, err := m.repo.Worktree()
	if err != nil {
		log.Print(err)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(Item)
			if ok {
				m.choice = string(i)

				refName := plumbing.NewBranchReferenceName(m.choice)
				opts := git.CheckoutOptions{
					Branch: refName,
					Create: false,
					Force:  false,
                    Keep: false,
				}

				e := w.Checkout(&opts)
				if e != nil {
                    log.Print("Checkout error: ", e)
				}
			}
			return m, tea.Quit
		}

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.choice != "" {
		return QuitTextStyle.Render(fmt.Sprintf("switched to branch '%s'", m.choice))
	}
	if m.quitting {
		return QuitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + m.list.View()
}
