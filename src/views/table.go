package views

import (
	"fmt"
	"planetary-observation/data"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TableMode string

type TableModel struct {
	Table         table.Model
	Mode          TableMode
	selectedRows  []int
	bottomMessage string
	csvTable      *TableModel
}

const maxColN = 10
const (
	CSV       = "csv"
	ColSelect = "colSelect"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Color("240"))

var infoMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#489eff"))
var warningMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffed68"))
var errorMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f94848"))

func (m TableModel) Init() tea.Cmd {
	if m.Mode == ColSelect {
		m.selectedRows = []int{}
		m.csvTable = nil
	}
	return nil
}

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.Table.Focused() {
				m.Table.Blur()
			} else {
				m.Table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "f":
			if len(m.selectedRows) > 0 {
				csvTable := CreateTableForCSVData(
					m.getSelectedRowKeys(),
					data.FilterCSVColumns(data.LoadPlanetarySystems(false)[:100], m.selectedRows),
				)
				m.csvTable = &csvTable
				return m.csvTable.Update(nil)
			} else {
				m.showWarning("You must select at least one column to display.")
			}
		case "enter", "e":
			if m.Mode == ColSelect {
				if len(m.selectedRows) < maxColN {
					selectedRow := m.Table.SelectedRow()
					if selectedRow[0] == "[ ]" {
						selectedRow[0] = "[x]"
						m.selectedRows = append(m.selectedRows, m.Table.Cursor())
					} else {
						selectedRow[0] = "[ ]"
						m.selectedRows = removeIntFromSlice(m.selectedRows, m.Table.Cursor())
					}

					rows := m.Table.Rows()
					rows[m.Table.Cursor()] = selectedRow
					m.Table.SetRows(rows)
					m.showInfo(fmt.Sprintf("%d/%d row selected", len(m.selectedRows), maxColN))
				} else {
					m.showWarning(
						fmt.Sprintf(
							"%d/%d row selected. You have already selected the maximum number of rows.",
							len(m.selectedRows),
							maxColN,
						),
					)
				}
			}

		case "p":
			if m.Mode == ColSelect {
				m.showInfo(fmt.Sprintf("selected rows: %v", m.selectedRows))
			}
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m TableModel) View() string {
	if m.csvTable != nil {
		return m.csvTable.View()
	}
	return baseStyle.Render(m.Table.View()) + "\n" + m.bottomMessage
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
	rows := m.Table.Rows()

	keys := []string{}
	for _, index := range m.selectedRows {
		keys = append(keys, rows[index][1])
	}

	return keys
}
