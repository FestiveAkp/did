package main

func (m model) viewTodoBody() string {
	return secondaryTextStyle.Render("Todo — coming soon") + "\n" +
		secondaryTextStyle.Render("All tasks with Todo status.") + "\n"
}
