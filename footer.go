package main

import (
	"fmt"
	"strings"
)

type footer struct {
	mode         mode
	input        string
	allStatuses  []Status
	statusCursor int
	filter       taskFilter
}

func newFooter(a tasksModel, filter taskFilter) footer {
	return footer{
		mode:         a.mode,
		input:        a.input.View(),
		allStatuses:  AllStatuses,
		statusCursor: a.statusCursor,
		filter:       filter,
	}
}

func (f footer) View() string {
	var b strings.Builder

	switch f.mode {
	case modeAdding:
		fmt.Fprintf(&b, "New task: %s\n", f.input)
		fmt.Fprintf(&b, "%s\n", secondaryTextStyle.Render("enter save | esc cancel"))
	case modeAddingActivity:
		fmt.Fprintf(&b, "New activity: %s\n", f.input)
		fmt.Fprintf(&b, "%s\n", secondaryTextStyle.Render("enter save | esc cancel"))
	case modePickingStatus:
		b.WriteString("Set status:\n")
		for i, st := range f.allStatuses {
			cursor := "  "
			if i == f.statusCursor {
				cursor = "> "
			}
			line := fmt.Sprintf("%s%s %s", cursor, st.Icon(), st.Label())
			if i == f.statusCursor {
				line = selectedItemStyle.Render(line)
			}
			fmt.Fprintf(&b, "%s\n", line)
		}
		fmt.Fprintf(&b, "%s\n", secondaryTextStyle.Render("enter select | esc cancel"))
	default:
		fmt.Fprintf(&b, "%s\n", secondaryTextStyle.Render("s set status | t create task | a create activity | d delete | q quit"))
	}

	return strings.TrimSuffix(b.String(), "\n")
}
