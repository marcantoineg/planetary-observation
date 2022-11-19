package views

import (
	"math"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh/terminal"
)

const horizontalPadding = 10

func GetTableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	return s
}

func CreateTableForCSVData(fields []string, data [][]string) TableModel {
	if len(fields) != len(data[0]) {
		panic("data must have the same number of columns than then fields")
	}
	return createTable(createCSVColumns(fields), createRows(data), CSV)
}

func CreateTableForColSelection(fields map[string]string) TableModel {
	return createTable(createColSelectionColumns(), createColSelectionRows(fields), ColSelect)
}

func createCSVColumns(fields []string) []table.Column {
	w, _ := getScreenSize()

	cols := []table.Column{}
	maxColWidth := (w - horizontalPadding) / len(fields)

	for _, field := range fields[:int(math.Min(float64(len(fields)), maxColN))] {
		cols = append(cols, table.Column{Title: field, Width: maxColWidth})
	}

	return cols
}

func createRows(data [][]string) []table.Row {
	rows := []table.Row{}
	for _, d := range data {
		rows = append(rows, d[:int(math.Min(float64(len(d)), maxColN))])
	}

	return rows
}

func createColSelectionColumns() []table.Column {
	w, _ := getScreenSize()
	firstColWidth := 8
	otherCols := (w - firstColWidth - horizontalPadding) / 2
	return []table.Column{
		{Title: "selected", Width: firstColWidth},
		{Title: "col name", Width: otherCols},
		{Title: "col desc", Width: otherCols},
	}
}

func createColSelectionRows(data map[string]string) []table.Row {
	rows := []table.Row{}

	for k, v := range data {
		rows = append(rows, table.Row{"[ ]", k, v})
	}

	return rows
}

// fun getScreenSize return a tuple containing the width and the height respectively.
func getScreenSize() (int, int) {
	width, height, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	return width, height
}

func createTable(cols []table.Column, rows []table.Row, mode TableMode) TableModel {
	w, h := getScreenSize()

	t := table.New(
		table.WithKeyMap(table.KeyMap{
			LineUp:       createBinding("w", "up"),
			LineDown:     createBinding("s", "down"),
			PageUp:       createBinding("W"),
			PageDown:     createBinding("S"),
			HalfPageUp:   createBinding("a"),
			HalfPageDown: createBinding("d"),
			GotoTop:      createBinding("g"),
			GotoBottom:   createBinding("G"),
		}),
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithWidth(w-4),
		table.WithHeight(h-10),
	)

	t.SetStyles(GetTableStyles())

	return TableModel{Table: t, Mode: mode}
}

func createBinding(binding ...string) key.Binding {
	b := key.NewBinding(func(b *key.Binding) { b.SetKeys(binding...) })
	return b
}
