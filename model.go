package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type viewKind int

const (
	viewTasks viewKind = iota
	viewTimeline
)

type taskFilter int

const (
	filterAll taskFilter = iota
	filterInProgress
	filterTodo
	filterDone
)

var taskFilters = []struct {
	label  string
	status Status
}{
	{"All", ""},
	{"In Progress", StatusInProgress},
	{"To Do", StatusTodo},
	{"Done", StatusDone},
}

type model struct {
	tasksView   tasksModel
	tasks       []Task
	activities  map[int64][]Activity
	currentView viewKind
	filter      taskFilter
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
		tasksView: newTasksModel(taskStore, activityStore),
	}
}

func (m model) filteredTasks() []Task {
	if m.filter == filterAll {
		return m.tasks
	}
	status := taskFilters[m.filter].status
	var result []Task
	for _, t := range m.tasks {
		if t.Status == status {
			result = append(result, t)
		}
	}
	return result
}

func (m model) Init() tea.Cmd {
	return m.tasksView.loadTasks
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
		filtered := m.filteredTasks()
		if m.tasksView.cursor >= len(filtered) {
			m.tasksView.cursor = len(filtered) - 1
		}
		if m.tasksView.cursor < 0 {
			m.tasksView.cursor = 0
		}
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		if m.tasksView.mode != modeNormal {
			var cmd tea.Cmd
			m.tasksView, cmd = m.tasksView.Update(msg, m.filteredTasks())
			return m, cmd
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.currentView = viewTasks
		case "2":
			m.currentView = viewTimeline
		case "h", "left":
			if m.currentView == viewTasks && m.filter > 0 {
				m.filter--
				m.tasksView.cursor = 0
			}
		case "l", "right":
			if m.currentView == viewTasks && int(m.filter) < len(taskFilters)-1 {
				m.filter++
				m.tasksView.cursor = 0
			}
		default:
			if m.currentView == viewTasks {
				var cmd tea.Cmd
				m.tasksView, cmd = m.tasksView.Update(msg, m.filteredTasks())
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
		{"Tasks", viewTasks, "1"},
		{"Timeline", viewTimeline, "2"},
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

func (m model) filterBar() string {
	var parts []string
	for i, f := range taskFilters {
		if taskFilter(i) == m.filter {
			parts = append(parts, activeTabStyle.Render(f.label))
		} else {
			parts = append(parts, inactiveTabStyle.Render(f.label))
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
	default:
		b.WriteString(m.filterBar())
		b.WriteString("\n\n")
		b.WriteString(m.tasksView.View(m.filteredTasks(), m.activities))
	}

	body := b.String()
	footerStr := newFooter(m.tasksView, m.filter).View()

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
