package main

import (
	"github.com/charmbracelet/bubbles/spinner"
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

var (
	isLoading = false
	isDisplay = false
	files     = make([]fs.FileInfo, 0)
	sub       = make(chan struct{})
)

type loadingMsg struct{}

type model struct {
	textInput textinput.Model
	table     table.Model
	spinner   spinner.Model
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

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		textInput: ti,
		table:     t,
		spinner:   sp,
	}
}

func waitForEndLoading() tea.Cmd {
	return func() tea.Msg {
		return loadingMsg(<-sub)
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
			walker.WG.Add(1)
			walker.WalkDir(m.textInput.Value())
			isLoading = true
			cmd = m.spinner.Tick

			go func() {
				walker.WG.Wait()
				files = walker.Files
				isLoading = false
				isDisplay = true
				sub <- struct{}{}
			}()

			return m, tea.Batch(cmd, waitForEndLoading())

		}

		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case loadingMsg:
		m.table.SetRows(filesToTableRows(files))
		return m, nil

	}
	m.table, cmd = m.table.Update(msg)
	m.textInput, cmd = m.textInput.Update(msg)

	if isLoading {
		m.spinner, cmd = m.spinner.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	var s string
	if isLoading {
		s = "Loading " + m.spinner.View() + "\n"
	} else if !isDisplay {
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
