package main

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// handles incoming events and updates the model accordingly
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":

			// Clear state
			m.bucketName = ""

			rows := []table.Row{}
			for _, bucketName := range m.bdb.ListBuckets() {
				stats := m.bdb.DescribeBucket(bucketName)
				rows = append(rows, table.Row{bucketName, strconv.Itoa(stats.KeyN), strconv.Itoa(stats.BucketN)})
			}

			columns := []table.Column{
				{Title: "Bucket", Width: 20},
				{Title: "KeyN", Width: 7},
				{Title: "BucketN", Width: 7},
			}

			m.table = table.New(
				table.WithColumns(columns),
				table.WithRows(rows),
				table.WithFocused(true),
				table.WithHeight(10),
			)
			m.table.SetStyles(getTableStyles())
			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.bucketName != "" {
				return m, tea.Println("Bucket already selected: ", m.bucketName)
			}

			// Update state with selected bucket
			m.bucketName = m.table.SelectedRow()[0]

			rows := []table.Row{}

			// List key/value pairs
			for _, kv := range m.bdb.ListKV(m.bucketName) {
				for k, v := range kv {
					rows = append(rows, table.Row{k, v.(string)})
				}
			}
			columns := []table.Column{
				{Title: "Key", Width: 30},
				{Title: "Value", Width: 30},
			}

			m.table = table.New(
				table.WithColumns(columns),
				table.WithRows(rows),
				table.WithFocused(true),
				table.WithHeight(10),
			)
			m.table.SetStyles(getTableStyles())
			m.table.SetStyles(getTableStyles())

			return m, nil
		}

	}

	m.table, cmd = m.table.Update(msg)

	// fmt.Println("BUCKET", bucketName)
	return m, cmd
}
