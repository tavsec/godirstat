package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"path/filepath"
	"tavsec/godirstat/services/walker"
)

type model struct {
	textInput textinput.Model
}

func initialModel() model {

	ti := textinput.New()
	currentPath, _ := os.Executable()
	ti.SetValue(filepath.Dir(currentPath))
	ti.Focus()
	ti.CharLimit = 512
	ti.Width = 512

	return model{
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		switch msg.Type {
		case tea.KeyEnter:
			walker.Walk(m.textInput.Value())
		}

		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := "Welsome to GoDirStat\n\n"
	s += "Where would you like to perform directory list?"
	s += m.textInput.View()
	s += "\nPress q to quit.\n"

	return s
}
