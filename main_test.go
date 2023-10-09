package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/muesli/termenv"
)

func init() {
	// This is required for CI to pass. See https://charm.sh/blog/teatest/
	lipgloss.SetColorProfile(termenv.Ascii)
}

func TestRepo(t *testing.T) {
	// Bless. This seems to work really well! The next steps here will be to try
	// load the model with these branch names
	repo, _ := git.PlainOpen("./testdata/TestRepo/")

	var branches []string
	branchIter, err := repo.Branches()

	if err != nil {
		log.Print("Error: ", err)
	}

	myFunc := func(ref *plumbing.Reference) error {
		trim := strings.TrimPrefix(ref.Name().String(), "refs/heads/")
		fmt.Println(trim)
		branches = append(branches, trim)
		return nil
	}

	branchIter.ForEach(myFunc)

	fmt.Println(branches)
	model := initModel(repo, branches)

	tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(300, 100))

	// Assert that the program, at some point, has the following byte string ... make a helper function?
	teatest.WaitFor(t, tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(
				bts,
				[]byte("1. branch-1"),
			)
		},
	)

	moveNDownAndSelectBranch(tm, 1)

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)

}

func TestOutput(t *testing.T) {

	// Tempdir to clone the repository
	dir, err := os.MkdirTemp("", "test-directory")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	repo := initTmpGitRepository(dir)

	branchNames := []string{
		"JOB-62131/JOB-76475/add-location-timers-to-fms",
		"JOB-62131/JOB-76477/store-feature-enablement",
		"JOB-62131/JOB-77400/show-modal-dialogue-on-disablement",
	}
	createBranches(repo, branchNames)

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

		moveNDownAndSelectBranch(tm, 1)

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

		moveNDownAndSelectBranch(tm, 2)

		tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

		out, err := io.ReadAll(tm.FinalOutput(t))
		if err != nil {
			t.Error(err)
		}
		teatest.RequireEqualOutput(t, out)

	})
}

func initModel(repo *git.Repository, branches []string) Model {
	var items []list.Item
	for _, branch := range branches {
		items = append(items, Item(branch))
	}
	l := list.New(items, ItemDelegate{}, DefaultWidth, ListHeight)
	return Model{list: l, repo: repo}
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

func moveNDownAndSelectBranch(tm *teatest.TestModel, down int) {

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

func createBranches(repo *git.Repository, branchNames []string) {
	for _, b := range branchNames {
		opts := &config.Branch{Name: b}
		err := repo.CreateBranch(opts)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func initTmpGitRepository(dir string) *git.Repository {
	_, err := git.PlainInit(dir, false)
	if err != nil {
		log.Fatal(err)
	}
	repo, err := git.PlainOpen(dir)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
