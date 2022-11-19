package profiling

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuProfile = flag.Bool("cpuProfiler", false, "write cpu profile to file")
var memProfile = flag.Bool("memProfiler", false, "write memory profile to file")

func PrintFlagsHelp() {
	println("-cpuProfiler - write cpu profile to file")
	println("-memProfiler - write memeory profile to file")
}

var cpuFile *os.File

func RunCpuProfiler() bool {
	if *cpuProfile {
		cpuFile, err := os.Create(".profiling/cpu.prof")
		if err != nil {
			panic(err)
		}
		println("i am running")
		pprof.StartCPUProfile(cpuFile)

		return true
	}
	return false
}

func StopCpuProfiler() {
	println("generating cpu profiler data")
	cpuFile.Close()
	pprof.StopCPUProfile()
}

func RunMemoryProfiler() {
	if *memProfile {
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
}
