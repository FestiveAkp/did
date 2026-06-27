package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type viewKind int

const (
	viewActive viewKind = iota
	viewTimeline
	viewTodo
	viewDone
)

type model struct {
	active      activeModel
	tasks       []Task
	activities  map[int64][]Activity
	currentView viewKind
	width       int
	height      int
	err         error
}

type tasksLoadedMsg struct {
	tasks      []Task
	activities map[int64][]Activity
}

type errMsg struct {
	err error
}

func newModel(taskStore *TaskStore, activityStore *ActivityStore) model {
	return model{
		active: newActiveModel(taskStore, activityStore),
	}
}

func (m model) Init() tea.Cmd {
	return m.active.loadTasks
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tasksLoadedMsg:
		m.tasks = msg.tasks
		m.activities = msg.activities
		if m.active.cursor >= len(m.tasks) {
			m.active.cursor = len(m.tasks) - 1
		}
		if m.active.cursor < 0 {
			m.active.cursor = 0
		}
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		if m.active.mode != modeNormal {
			var cmd tea.Cmd
			m.active, cmd = m.active.Update(msg, m.tasks)
			return m, cmd
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.currentView = viewActive
		case "2":
			m.currentView = viewTimeline
		case "3":
			m.currentView = viewTodo
		case "4":
			m.currentView = viewDone
		default:
			if m.currentView == viewActive {
				var cmd tea.Cmd
				m.active, cmd = m.active.Update(msg, m.tasks)
				return m, cmd
			}
		}
		return m, nil
	}

	return m, nil
}

func (m model) tabBar() string {
	tabs := []struct {
		label string
		view  viewKind
		key   string
	}{
		{"Active", viewActive, "1"},
		{"Timeline", viewTimeline, "2"},
		{"Todo", viewTodo, "3"},
		{"Done", viewDone, "4"},
	}

	var parts []string
	for _, t := range tabs {
		label := fmt.Sprintf("%s %s", t.key, t.label)
		if t.view == m.currentView {
			parts = append(parts, activeTabStyle.Render(label))
		} else {
			parts = append(parts, inactiveTabStyle.Render(label))
		}
	}
	return strings.Join(parts, "  ")
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(m.tabBar())
	b.WriteString("\n\n")

	if m.err != nil {
		fmt.Fprintf(&b, "error: %v\n\n", m.err)
	}

	switch m.currentView {
	case viewTimeline:
		b.WriteString(m.viewTimelineBody())
	case viewTodo:
		b.WriteString(m.viewTodoBody())
	case viewDone:
		b.WriteString(m.viewDoneBody())
	default:
		b.WriteString(m.active.View(m.tasks, m.activities))
	}

	body := b.String()
	footerStr := newFooter(m.active).View()

	if m.height > 0 {
		bodyLines := strings.Count(body, "\n")
		footerLines := strings.Count(footerStr, "\n")
		padding := m.height - bodyLines - footerLines - 1
		if padding > 0 {
			body += strings.Repeat("\n", padding)
		}
	} else {
		body += "\n"
	}

	return body + footerStr
}
