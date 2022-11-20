package views

import (
	"math"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	colsHorizontalPadding  = 12
	tableHorizontalPadding = 10
	tableVerticalPadding   = 10
	rowsPerPage            = 100
)

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
	maxColWidth := (w) / len(fields)

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
	otherCols := (w - firstColWidth) / 2
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
	width, height, err := terminal.GetSize(0)
	if err != nil {
		panic(err)
	}

	return width, height
}

func createTable(cols []table.Column, rows []table.Row, mode TableMode) TableModel {
	w, h := getScreenSize()

	t := table.New(
		table.WithKeyMap(table.KeyMap{
			LineUp:     createBinding("w", "up"),
			LineDown:   createBinding("s", "down"),
			PageUp:     createBinding("W"),
			PageDown:   createBinding("S"),
			GotoTop:    createBinding("g"),
			GotoBottom: createBinding("G"),
		}),
		table.WithColumns(cols),
		table.WithRows(rows[:int(math.Min(float64(len(rows)), 100))]),
		table.WithFocused(true),
		table.WithWidth(w),
		table.WithHeight(h-tableVerticalPadding),
	)
	t.SetStyles(GetTableStyles())

	p := paginator.New()
	p.PerPage = rowsPerPage
	p.SetTotalPages(len(rows))

	return NewTableModel(t, p, mode, rows)
}

func createBinding(binding ...string) key.Binding {
	b := key.NewBinding(func(b *key.Binding) { b.SetKeys(binding...) })
	return b
}

func (m *TableModel) showInfo(message string) {
	m.showMessage(message, &infoMessageStyle)
}

func (m *TableModel) showWarning(message string) {
	m.showMessage(message, &warningMessageStyle)
}

func (m *TableModel) showError(message string) {
	m.showMessage(message, &errorMessageStyle)
}

func (m *TableModel) showMessage(message string, style *lipgloss.Style) {
	finalMessage := message
	if style != nil {
		finalMessage = style.Render(message)
	}

	m.bottomMessage = finalMessage
}

func removeIntFromSlice(slice []int, valueToRemove int) []int {
	if len(slice) == 0 {
		return slice
	}

	var index int = -1
	for i, v := range slice {
		if v == valueToRemove {
			index = i
			break
		}
	}

	if index == -1 {
		return slice
	} else if len(slice) <= index+1 {
		return slice[:index]
	}

	return append(slice[:index], slice[index+1:]...)
}

func (m *TableModel) getSelectedRowKeys() []string {
	rows := m.table.Rows()

	keys := []string{}
	for _, index := range m.selectedRows {
		keys = append(keys, rows[index][1])
	}

	return keys
}

func (m *TableModel) changePage(incrementPage bool, n int) {
	for i := 0; i < n; i++ {
		if incrementPage {
			m.paginator.NextPage()
		} else {
			m.paginator.PrevPage()
		}
	}
	start, end := m.paginator.GetSliceBounds(len(m.allRows))
	m.table.SetRows(m.allRows[start:end])
}
