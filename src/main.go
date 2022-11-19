package main

import (
	"fmt"
	"log"
	"os"
	"planetary-observation/data"
	"planetary-observation/views"
	"runtime"
	"runtime/pprof"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	f, err := os.Create(".profiling/cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// for i := 0; i < 100_000; i++ {
	// 	squared := i * i
	// 	println(doSomething(squared))
	// }

	fields := data.LoadFields(false)
	// planetarySystems := data.LoadPlanetarySystems()

	launchTeaApp(views.CreateTableForColSelection(fields))
	// launchTeaApp(views.CreateTableForCSVData(planetarySystems))

	runMemoryProfiler()
}

func runMemoryProfiler() {
	f, err := os.Create(".profiling/mem.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	runtime.GC() // for up-to-date stats
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func doSomething(squared int) int {
	return squared + squared/3
}

func launchTeaApp(m tea.Model) {
	if err := tea.NewProgram(m).Start(); err != nil {
		panic(fmt.Sprintf("an error occured running the TUI app.\n%s", err))
	}
}
