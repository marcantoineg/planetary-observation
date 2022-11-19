package main

import (
	"flag"
	"fmt"
	"planetary-observation/data"
	"planetary-observation/profiling"
	"planetary-observation/views"

	tea "github.com/charmbracelet/bubbletea"
)

var helpFlagValue bool

func main() {
	if handleHelpFlag() {
		// kills the program if help message was printed.
		return
	}

	if profiling.RunCpuProfiler() {
		defer profiling.StopCpuProfiler()
	}

	fields := data.LoadFields(false)
	// planetarySystems := data.LoadPlanetarySystems()

	launchTeaApp(views.CreateTableForColSelection(fields))
	// launchTeaApp(views.CreateTableForCSVData(planetarySystems))

	profiling.RunMemoryProfiler()
}

func launchTeaApp(m tea.Model) {
	if err := tea.NewProgram(m).Start(); err != nil {
		panic(fmt.Sprintf("an error occured running the TUI app.\n%s", err))
	}
}

// handleFlags creates, parse and handle the `-help` flag and retuns true if help flag is active.
func handleHelpFlag() bool {
	flag.BoolVar(&helpFlagValue, "h", false, "Show program information.")
	flag.BoolVar(&helpFlagValue, "help", false, "Show program information.")
	flag.Parse()
	if helpFlagValue {
		println("\nplanetary-observation - A simple TUI to view CSV data.\n")
		println("-h - Show program information")
		profiling.PrintFlagsHelp()
		return true
	}
	return false
}
