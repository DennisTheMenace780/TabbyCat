package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func init() {
    // This is required for CI to pass. See https://charm.sh/blog/teatest/
	lipgloss.SetColorProfile(termenv.Ascii)
}

func TestFullOutput(t *testing.T) {
	model := initialModel()
	tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(300, 100))

	teatest.WaitFor(t, tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("1. JOB-62131/JOB-76475/add-location-timers-to-fms"))
		},
	)

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("j"),
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

	finalView := tm.FinalModel(t).View()
	// Want to do a contains assertion here to ignore the way output is displayed.
	// Just care about having the right sub-string in the outputs
	expectedOutput := "checking out JOB-62131/JOB-76477/store-feature-enablement"
	assert.Contains(t, finalView, expectedOutput)

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func initialModel() model {

	branches := []list.Item{
		item("JOB-62131/JOB-76475/add-location-timers-to-fms"),
		item("JOB-62131/JOB-76477/store-feature-enablement"),
		item("JOB-62131/JOB-77400/show-modal-dialogue-on-disablement"),
	}

	const defaultWidth = 20
	const listHeight = 14

	l := list.New(branches, itemDelegate{}, defaultWidth, listHeight)
	return model{list: l}
}
