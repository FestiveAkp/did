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
}

func newFooter(a activeModel) footer {
	return footer{
		mode:         a.mode,
		input:        a.input.View(),
		allStatuses:  AllStatuses,
		statusCursor: a.statusCursor,
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
		fmt.Fprintf(&b, "%s\n", secondaryTextStyle.Render("↑/k up | ↓/j down | Enter select | Esc cancel"))
	default:
		fmt.Fprintf(&b, "%s\n", secondaryTextStyle.Render("1-4 switch view | ↑/k ↓/j navigate | s status | a add | n activity | d delete | q quit"))
	}

	return strings.TrimSuffix(b.String(), "\n")
}
