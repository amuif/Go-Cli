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
	table table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check if msg is a keypress to quit
	var cmd tea.Cmd
	var quitCmd tea.Cmd

	var quitModel tea.Model
	quitModel, quitCmd = m.Quit(msg)
	m = quitModel.(model)
	if quitCmd != nil {
		return m, quitCmd
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) Quit(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	style := lipgloss.NewStyle().Margin(1, 2)
	return style.Render(m.table.View())
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
		table.WithHeight(20),
	)
	// Add a visible border style
	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true)
	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)
	t.SetStyles(styles)

	m := model{table: t}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running TUI:", err)
		os.Exit(1)
	}
}
