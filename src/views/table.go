package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type TableMode string
type TableModel struct {
	table         table.Model
	mode          TableMode
	selectedRows  []int
	bottomMessage string
	csvTable      *TableModel
	paginator     paginator.Model
	allRows       []table.Row
}

const maxColN = 8
const (
	CSV       = "csv"
	ColSelect = "colSelect"
)

func NewTableModel(t table.Model, p paginator.Model, mode TableMode, items []table.Row) TableModel {
	return TableModel{
		table:     t,
		mode:      mode,
		paginator: p,
		allRows:   items,
	}
}

func (m TableModel) Init() tea.Cmd {
	if m.mode == ColSelect {
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
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "f":
			if m.mode == ColSelect {
				return handleCsvTableGeneration(&m)
			}
		case "enter", "e":
			if m.mode == ColSelect {
				handleColSelection(&m)
			}

		case "d":
			if m.mode == CSV {
				m.changePage(true, 1)
			}

		case "a":
			if m.mode == CSV {
				m.changePage(false, 1)
			}

		case "D":
			if m.mode == CSV {
				m.changePage(true, 10)
			}

		case "A":
			if m.mode == CSV {
				m.changePage(false, 10)
			}

		case "p":
			if m.mode == ColSelect {
				m.showInfo(fmt.Sprintf("selected rows: %v", m.selectedRows))
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m TableModel) View() string {
	var sb strings.Builder

	sb.WriteString(baseStyle.Render(m.table.View()) + "\n")
	sb.WriteString(m.bottomMessage + "\n")
	sb.WriteString(fmt.Sprintf("Col: %d/%d\n", m.table.Cursor()+1, len(m.table.Rows())))
	sb.WriteString(fmt.Sprintf("Page: %s\n", m.paginator.View()))
	return sb.String()
}
