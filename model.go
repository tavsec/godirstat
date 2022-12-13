package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"tavsec/godirstat/services/walker"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	textInput textinput.Model
	files     []fs.FileInfo
	table     table.Model
	isDisplay bool
}

func initialModel() model {

	ti := textinput.New()
	currentPath, _ := os.Executable()
	ti.SetValue(filepath.Dir(currentPath))
	ti.Focus()
	ti.CharLimit = 512
	ti.Width = 512

	columns := []table.Column{
		{Title: "Is Directory", Width: 10},
		{Title: "File", Width: 100},
		{Title: "Size", Width: 20},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return model{
		textInput: ti,
		files:     make([]fs.FileInfo, 0),
		table:     t,
		isDisplay: false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.Type {
		case tea.KeyEnter:
			_, m.files = walker.Walk(m.textInput.Value())
			m.table.SetRows(filesToTableRows(m.files))
			m.isDisplay = true
		}

		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	}
	m.table, cmd = m.table.Update(msg)
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var s string
	if !m.isDisplay {
		s = "Welcome to GoDirStat\n\n"
		s += "Where would you like to perform directory list?"
		s += m.textInput.View()
		s += "\nPress q to quit.\n"
	} else {
		s = baseStyle.Render(m.table.View())

	}

	return s
}

func filesToTableRows(files []fs.FileInfo) []table.Row {
	rows := make([]table.Row, 0)
	for _, file := range files {
		isDirIcon := "✔️"
		if !file.IsDir() {
			isDirIcon = ""
		}
		rows = append(rows, table.Row{isDirIcon, file.Name(), strconv.FormatInt(file.Size(), 10)})
	}

	return rows
}
