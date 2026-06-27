package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


type mode int

const (
	modeNormal mode = iota
	modeAdding
	modePickingStatus
	modeAddingActivity
)

type model struct {
	store         *TaskStore
	activityStore *ActivityStore
	tasks         []Task
	activities    map[int64][]Activity
	cursor        int

	mode  mode
	input textinput.Model

	statusCursor int

	width  int
	height int

	err error
}

type tasksLoadedMsg struct {
	tasks      []Task
	activities map[int64][]Activity
}

type errMsg struct {
	err error
}

func newModel(store *TaskStore, activityStore *ActivityStore) model {
	ti := textinput.New()
	ti.CharLimit = 200

	return model{
		store:         store,
		activityStore: activityStore,
		input:         ti,
	}
}

func (m model) Init() tea.Cmd {
	return m.loadTasks
}

func (m model) loadTasks() tea.Msg {
	tasks, err := m.store.List()
	if err != nil {
		return errMsg{err}
	}
	activities := make(map[int64][]Activity)
	for _, t := range tasks {
		acts, err := m.activityStore.ListForTask(t.ID)
		if err != nil {
			return errMsg{err}
		}
		activities[t.ID] = acts
	}
	return tasksLoadedMsg{tasks, activities}
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
		if m.cursor >= len(m.tasks) {
			m.cursor = len(m.tasks) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case modeAdding:
			return m.updateAdding(msg)
		case modePickingStatus:
			return m.updatePickingStatus(msg)
		case modeAddingActivity:
			return m.updateAddingActivity(msg)
		default:
			return m.updateNormal(msg)
		}
	}

	return m, nil
}

func (m model) updateAdding(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeNormal
		m.input.Reset()
		m.input.Blur()
		return m, nil

	case "enter":
		title := strings.TrimSpace(m.input.Value())
		m.mode = modeNormal
		m.input.Reset()
		m.input.Blur()
		if title == "" {
			return m, nil
		}
		return m, func() tea.Msg {
			if _, err := m.store.Create(title); err != nil {
				return errMsg{err}
			}
			return m.loadTasks()
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) updateAddingActivity(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeNormal
		m.input.Reset()
		m.input.Blur()
		return m, nil

	case "enter":
		note := strings.TrimSpace(m.input.Value())
		m.mode = modeNormal
		m.input.Reset()
		m.input.Blur()
		if note == "" {
			return m, nil
		}
		taskID := m.tasks[m.cursor].ID
		return m, func() tea.Msg {
			if _, err := m.activityStore.Create(taskID, note); err != nil {
				return errMsg{err}
			}
			return m.loadTasks()
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "j", "down":
		if m.cursor < len(m.tasks)-1 {
			m.cursor++
		}

	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}

	case "a":
		m.mode = modeAdding
		m.input.Placeholder = "Task title"
		m.input.Focus()
		return m, textinput.Blink

	case "n":
		if len(m.tasks) == 0 {
			return m, nil
		}
		m.mode = modeAddingActivity
		m.input.Placeholder = "Activity note"
		m.input.Focus()
		return m, textinput.Blink

	case "s":
		if len(m.tasks) == 0 {
			return m, nil
		}
		m.mode = modePickingStatus
		m.statusCursor = statusIndex(m.tasks[m.cursor].Status)
		return m, nil

	case "d":
		if len(m.tasks) == 0 {
			return m, nil
		}
		task := m.tasks[m.cursor]
		return m, func() tea.Msg {
			if err := m.store.Delete(task.ID); err != nil {
				return errMsg{err}
			}
			return m.loadTasks()
		}
	}

	return m, nil
}

func (m model) updatePickingStatus(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeNormal
		return m, nil

	case "j", "down":
		if m.statusCursor < len(AllStatuses)-1 {
			m.statusCursor++
		}

	case "k", "up":
		if m.statusCursor > 0 {
			m.statusCursor--
		}

	case "enter":
		m.mode = modeNormal
		task := m.tasks[m.cursor]
		status := AllStatuses[m.statusCursor]
		return m, func() tea.Msg {
			if err := m.store.SetStatus(task.ID, status); err != nil {
				return errMsg{err}
			}
			return m.loadTasks()
		}
	}

	return m, nil
}

func statusIndex(s Status) int {
	for i, st := range AllStatuses {
		if st == s {
			return i
		}
	}
	return 0
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString("did — tasks\n\n")

	if m.err != nil {
		fmt.Fprintf(&b, "error: %v\n\n", m.err)
	}

	if len(m.tasks) == 0 {
		b.WriteString("No tasks yet. Press 'a' to add one.\n\n")
	}

	for i, t := range m.tasks {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		line := fmt.Sprintf("%s%s %s", cursor, t.Status.Icon(), t.Title)
		if i == m.cursor {
			line = selectedItemStyle.Render(line)
		}
		fmt.Fprintf(&b, "%s\n", line)

		for _, a := range m.activities[t.ID] {
			date := a.CreatedAt.Local().Format("Jan 2 3:04pm")
			fmt.Fprintf(&b, "%s\n", secondaryTextStyle.Render(fmt.Sprintf("    · %s  %s", a.Note, date)))
		}
	}

	body := b.String()
	footerStr := newFooter(m).View()

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
