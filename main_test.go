package main

import (
	"bytes"
	"io"
	"log"
	"os"
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

type WorktreeSuite struct {
	BaseSuite
}

func ExamplePlainClone() {
	// Tempdir to clone the repository
	dir, err := os.MkdirTemp("", "clone-example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir) // clean up

	// Clones the repository into the given dir, just as a normal git clone does
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL: "https://github.com/git-fixtures/basic.git",
	})

	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
	// Prints the content of the CHANGELOG file from the cloned repository
	// changelog, err := os.Open(filepath.Join(dir, "CHANGELOG"))
	changelog, err := os.Open("CHANGELOG")
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(os.Stdout, changelog)
	// Output: Initial changelog
}

func TestOutput(t *testing.T) {

	ExamplePlainClone()

	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}

	model := initialModel(repo)

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

	t.Run("Raises checkout error when branches are modified", func(t *testing.T) {
		tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(300, 100))

		teatest.WaitFor(t, tm.Output(),
			func(bts []byte) bool {
				return bytes.Contains(
					bts,
					[]byte("1. JOB-62131/JOB-76475/add-location-timers-to-fms"),
				)
			},
		)
	})
}

func initialModel(repo *git.Repository) Model {

	branches := []list.Item{
		Item("JOB-62131/JOB-76475/add-location-timers-to-fms"),
		Item("JOB-62131/JOB-76477/store-feature-enablement"),
		Item("JOB-62131/JOB-77400/show-modal-dialogue-on-disablement"),
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
