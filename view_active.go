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

type activeModel struct {
	taskStore     *TaskStore
	activityStore *ActivityStore
	cursor        int
	mode          mode
	input         textinput.Model
	statusCursor  int
}

func newActiveModel(taskStore *TaskStore, activityStore *ActivityStore) activeModel {
	ti := textinput.New()
	ti.CharLimit = 200
	return activeModel{
		taskStore:     taskStore,
		activityStore: activityStore,
		input:         ti,
	}
}

func (a activeModel) loadTasks() tea.Msg {
	tasks, err := a.taskStore.List()
	if err != nil {
		return errMsg{err}
	}
	activities := make(map[int64][]Activity)
	for _, t := range tasks {
		acts, err := a.activityStore.ListForTask(t.ID)
		if err != nil {
			return errMsg{err}
		}
		activities[t.ID] = acts
	}
	return tasksLoadedMsg{tasks, activities}
}

func (a activeModel) Update(msg tea.KeyMsg, tasks []Task) (activeModel, tea.Cmd) {
	switch a.mode {
	case modeAdding:
		return a.updateAdding(msg)
	case modePickingStatus:
		return a.updatePickingStatus(msg, tasks)
	case modeAddingActivity:
		return a.updateAddingActivity(msg, tasks)
	default:
		return a.updateNormal(msg, tasks)
	}
}

func (a activeModel) updateAdding(msg tea.KeyMsg) (activeModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		a.mode = modeNormal
		a.input.Reset()
		a.input.Blur()
		return a, nil
	case "enter":
		title := strings.TrimSpace(a.input.Value())
		a.mode = modeNormal
		a.input.Reset()
		a.input.Blur()
		if title == "" {
			return a, nil
		}
		return a, func() tea.Msg {
			if _, err := a.taskStore.Create(title); err != nil {
				return errMsg{err}
			}
			return a.loadTasks()
		}
	}
	var cmd tea.Cmd
	a.input, cmd = a.input.Update(msg)
	return a, cmd
}

func (a activeModel) updateAddingActivity(msg tea.KeyMsg, tasks []Task) (activeModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		a.mode = modeNormal
		a.input.Reset()
		a.input.Blur()
		return a, nil
	case "enter":
		note := strings.TrimSpace(a.input.Value())
		a.mode = modeNormal
		a.input.Reset()
		a.input.Blur()
		if note == "" {
			return a, nil
		}
		taskID := tasks[a.cursor].ID
		return a, func() tea.Msg {
			if _, err := a.activityStore.Create(taskID, note); err != nil {
				return errMsg{err}
			}
			return a.loadTasks()
		}
	}
	var cmd tea.Cmd
	a.input, cmd = a.input.Update(msg)
	return a, cmd
}

func (a activeModel) updateNormal(msg tea.KeyMsg, tasks []Task) (activeModel, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if a.cursor < len(tasks)-1 {
			a.cursor++
		}
	case "k", "up":
		if a.cursor > 0 {
			a.cursor--
		}
	case "a":
		a.mode = modeAdding
		a.input.Placeholder = "Task title"
		a.input.Focus()
		return a, textinput.Blink
	case "n":
		if len(tasks) == 0 {
			return a, nil
		}
		a.mode = modeAddingActivity
		a.input.Placeholder = "Activity note"
		a.input.Focus()
		return a, textinput.Blink
	case "s":
		if len(tasks) == 0 {
			return a, nil
		}
		a.mode = modePickingStatus
		a.statusCursor = statusIndex(tasks[a.cursor].Status)
		return a, nil
	case "d":
		if len(tasks) == 0 {
			return a, nil
		}
		task := tasks[a.cursor]
		return a, func() tea.Msg {
			if err := a.taskStore.Delete(task.ID); err != nil {
				return errMsg{err}
			}
			return a.loadTasks()
		}
	}
	return a, nil
}

func (a activeModel) updatePickingStatus(msg tea.KeyMsg, tasks []Task) (activeModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		a.mode = modeNormal
		return a, nil
	case "j", "down":
		if a.statusCursor < len(AllStatuses)-1 {
			a.statusCursor++
		}
	case "k", "up":
		if a.statusCursor > 0 {
			a.statusCursor--
		}
	case "enter":
		a.mode = modeNormal
		task := tasks[a.cursor]
		status := AllStatuses[a.statusCursor]
		return a, func() tea.Msg {
			if err := a.taskStore.SetStatus(task.ID, status); err != nil {
				return errMsg{err}
			}
			return a.loadTasks()
		}
	}
	return a, nil
}

func statusIndex(s Status) int {
	for i, st := range AllStatuses {
		if st == s {
			return i
		}
	}
	return 0
}

func (a activeModel) View(tasks []Task, activities map[int64][]Activity) string {
	var b strings.Builder

	if len(tasks) == 0 {
		b.WriteString("No tasks yet. Press 'a' to add one.\n")
		return b.String()
	}

	for i, t := range tasks {
		cursor := "  "
		if i == a.cursor {
			cursor = "> "
		}
		line := fmt.Sprintf("%s%s %s", cursor, t.Status.Icon(), t.Title)
		if i == a.cursor {
			line = selectedItemStyle.Render(line)
		}
		fmt.Fprintf(&b, "%s\n", line)

		for _, act := range activities[t.ID] {
			date := act.CreatedAt.Local().Format("Jan 2 3:04pm")
			fmt.Fprintf(&b, "%s\n", secondaryTextStyle.Render(fmt.Sprintf("    · %s  %s", act.Note, date)))
		}
	}

	return b.String()
}
