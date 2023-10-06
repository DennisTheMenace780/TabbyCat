package main

import (
	"bytes"
	"io"
	"log"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/go-git/go-git/v5"
	"github.com/muesli/termenv"
)

func init() {
	// This is required for CI to pass. See https://charm.sh/blog/teatest/
	lipgloss.SetColorProfile(termenv.Ascii)
}

func TestOutput(t *testing.T) {

	model := initialModel()

	t.Run("Moving down once and selecting a branch", func(t *testing.T) {
		tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(300, 100))

        // Assert that the program, at some point, has the following byte string ... make a helper function?
		teatest.WaitFor(t, tm.Output(),
			func(bts []byte) bool {
				return bytes.Contains(
					bts,
					[]byte("1. JOB-62131/JOB-76475/add-location-timers-to-fms"),
				)
			},
		)

		moveDownAndSelectBranch(tm, 1)

		tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

		out, err := io.ReadAll(tm.FinalOutput(t))
		if err != nil {
			t.Error(err)
		}
		teatest.RequireEqualOutput(t, out)

	})

	t.Run("Moving down twice and selecting a branch", func(t *testing.T) {
		tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(300, 100))

		teatest.WaitFor(t, tm.Output(),
			func(bts []byte) bool {
				return bytes.Contains(
					bts,
					[]byte("1. JOB-62131/JOB-76475/add-location-timers-to-fms"),
				)
			},
		)

		moveDownAndSelectBranch(tm, 2)

		tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

		out, err := io.ReadAll(tm.FinalOutput(t))
		if err != nil {
			t.Error(err)
		}
		teatest.RequireEqualOutput(t, out)

	})
}

func initialModel() Model {

	branches := []list.Item{
		Item("JOB-62131/JOB-76475/add-location-timers-to-fms"),
		Item("JOB-62131/JOB-76477/store-feature-enablement"),
		Item("JOB-62131/JOB-77400/show-modal-dialogue-on-disablement"),
	}

	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}

	l := list.New(branches, ItemDelegate{}, DefaultWidth, ListHeight)
	return Model{list: l, repo: repo}
}

func moveDownAndSelectBranch(tm *teatest.TestModel, down int) {

	for i := 0; i < down; i++ {
		tm.Send(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune("j"),
		})
	}

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})
}
