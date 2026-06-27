package main

func (m model) viewTimelineBody() string {
	return secondaryTextStyle.Render("Timeline — coming soon") + "\n" +
		secondaryTextStyle.Render("Reverse chronological activity feed. Filter by: yesterday · past week · past month · past year") + "\n"
}
