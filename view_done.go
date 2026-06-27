package main

func (m model) viewDoneBody() string {
	return secondaryTextStyle.Render("Done — coming soon") + "\n" +
		secondaryTextStyle.Render("All tasks with Done status.") + "\n"
}
