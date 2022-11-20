package views

import (
	"fmt"
	"planetary-observation/data"

	tea "github.com/charmbracelet/bubbletea"
)

func handleColSelection(m *TableModel) {
	if len(m.selectedRows) < maxColN {
		selectedRow := m.table.SelectedRow()
		if selectedRow[0] == "[ ]" {
			selectedRow[0] = "[x]"
			m.selectedRows = append(m.selectedRows, m.table.Cursor())
		} else {
			selectedRow[0] = "[ ]"
			m.selectedRows = removeIntFromSlice(m.selectedRows, m.table.Cursor())
		}

		rows := m.table.Rows()
		rows[m.table.Cursor()] = selectedRow
		m.table.SetRows(rows)
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

func handleCsvTableGeneration(m *TableModel) (tea.Model, tea.Cmd) {
	if len(m.selectedRows) > 0 {
		csvTable := CreateTableForCSVData(
			m.getSelectedRowKeys(),
			data.FilterCSVColumns(data.LoadPlanetarySystems(false), m.selectedRows),
		)
		m.csvTable = &csvTable
		return m.csvTable.Update(nil)
	} else {
		m.showWarning("You must select at least one column to display.")
		return m, nil
	}
}
