package main

import "github.com/charmbracelet/lipgloss"

// renders the UI based on the data in the model
func (m model) View() string {
	// render a horizontal table
	leftCol := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderRight(true).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Right,
				"DB",
				"Bucket",
			),
		)

	bucketNameRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	bucketName := "n/a"
	if m.bucketName != "" {
		bucketNameRenderer.Foreground(lipgloss.Color("51"))
		bucketName = m.bucketName
	}
	rightCol := lipgloss.
		NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color("51")).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left, m.bdb.Db.Path(),
				bucketNameRenderer.Render(bucketName),
			),
		)

	helpView := m.help.View(m.keys)

	return lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("240")).
		Border(lipgloss.RoundedBorder()).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left, leftCol, rightCol),
		) + "\n" + baseStyle.Render(m.table.View()) + "\n" + helpView
}
