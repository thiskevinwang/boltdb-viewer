package main

import (
	"fmt"
	"os"
	"strconv"

	db "main/db"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	keys  keyMap
	table table.Model
	help  help.Model

	// Our BoltDB instance wrapper
	// db *bolt.DB
	bdb *db.Bolt

	// the current bolt DB bucket name
	bucketName string
}

func getTableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("51"))
	return s
}

func main() {
	bdb := db.Bolt{}
	bdb.Init(os.Args[1])
	defer bdb.Close()

	// Initialize the table
	rows := []table.Row{}
	for _, bucketName := range bdb.ListBuckets() {
		stats := bdb.DescribeBucket(bucketName)
		rows = append(rows, table.Row{bucketName, strconv.Itoa(stats.KeyN), strconv.Itoa(stats.BucketN)})
	}

	columns := []table.Column{
		{Title: "Bucket", Width: 20},
		{Title: "KeyN", Width: 7},
		{Title: "BucketN", Width: 7},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	t.SetStyles(getTableStyles())

	m := model{
		keys:       keys,
		table:      t,
		bdb:        &bdb,
		bucketName: "",
		help:       help.New(),
	}

	// Run
	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
