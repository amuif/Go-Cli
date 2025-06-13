package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Todo struct {
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}
type Todos []Todo

type model struct {
	cursor    int
	table     table.Model
	selected  map[int]struct{}
	showPopup bool
	popupText string
}

func initialModel() model {
	todos := Todos{
		{Title: "Learn Bubble Tea", Completed: false, CreatedAt: time.Now()},
	}

	columns := []table.Column{
		{Title: "Title", Width: 30},
		{Title: "Completed", Width: 10},
		{Title: "Created At", Width: 20},
		{Title: "Completed At", Width: 20},
	}

	rows := todos.toTableRows()

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true)
	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("#f8f8f2")).
		Background(lipgloss.Color("#282a36")).
		Bold(true)

	t.SetStyles(styles)

	return model{
		selected: make(map[int]struct{}),
		table:    t,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check if msg is a keypress to quit
	var cmd tea.Cmd
	var quitCmd tea.Cmd

	var quitModel tea.Model
	quitModel, quitCmd = m.quit(msg)
	m = quitModel.(model)
	if quitCmd != nil {
		return m, quitCmd
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "E", "e":
			m.showPopup = true
			m.popupText = "Enter the new name of todo"
			return m, nil
		case "A", "a":
			m.showPopup = true
			m.popupText = "Add a new todo"
			return m, nil

		case "D", "d":
			m.showPopup = true
			m.popupText = "Are you sure you want to delete?"
			return m, nil
		case "T", "t":
			return m, nil
		case "esc":
			m.showPopup = false
			m.popupText = ""
			return m, nil
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) quit(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	tableView := m.table.View()
	footer := "\n\nPress E to edit , A to add , D to delete , T to toggle"

	baseStyle := lipgloss.NewStyle().Margin(1, 2)

	if m.showPopup {
		popup := fmt.Sprintf(" %s\n(Press Esc to close)", m.popupText)

		popupBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Align(lipgloss.Center).
			Width(50).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("238")).
			Render(popup)

		overlay := lipgloss.Place(
			80, 15,
			lipgloss.Center, lipgloss.Center,
			popupBox,
		)

		return baseStyle.Render(tableView + "\n" + overlay + footer)
	}

	return baseStyle.Render(tableView + footer)
}

func (todos *Todos) add(title string) {
	todo := Todo{
		Title:       title,
		Completed:   false,
		CompletedAt: nil,
		CreatedAt:   time.Now(),
	}
	*todos = append(*todos, todo)
}

func (todos *Todos) validateIndex(index int) error {
	if index < 0 || index >= len(*todos) {
		return errors.New("invalid index")
	}
	return nil
}

func (todos *Todos) delete(index int) error {
	if err := todos.validateIndex(index); err != nil {
		return err
	}
	*todos = append((*todos)[:index], (*todos)[index+1:]...)
	return nil
}

func (todos *Todos) toggle(index int) error {
	if err := todos.validateIndex(index); err != nil {
		return err
	}
	isCompleted := (*todos)[index].Completed
	if !isCompleted {
		completionTime := time.Now()
		(*todos)[index].CompletedAt = &completionTime
	}
	(*todos)[index].Completed = !isCompleted
	return nil
}

func (todos *Todos) edit(index int, title string) error {
	if err := todos.validateIndex(index); err != nil {
		return err
	}
	(*todos)[index].Title = title
	return nil
}

func (todos Todos) toTableRows() []table.Row {
	var rows []table.Row
	for _, t := range todos {
		completed := "❌"
		if t.Completed {
			completed = "✔️"
		}
		completedAt := "-"
		if t.CompletedAt != nil {
			completedAt = t.CompletedAt.Format("2006-01-02 15:04")
		}
		rows = append(rows, table.Row{
			t.Title,
			completed,
			t.CreatedAt.Format("2006-01-02 15:04"),
			completedAt,
		})
	}
	return rows
}

func (todos Todos) print() {
	columns := []table.Column{
		{Title: "Title", Width: 30},
		{Title: "Completed", Width: 10},
		{Title: "Created At", Width: 20},
		{Title: "Completed At", Width: 20},
	}

	rows := todos.toTableRows()

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Add a visible border style
	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true)
	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("#f8f8f2")).
		Background(lipgloss.Color("#282a36")).
		Bold(true)
	t.SetStyles(styles)

	m := model{table: t}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running TUI:", err)
		os.Exit(1)
	}

}
