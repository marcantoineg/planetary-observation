run:
	cd src && go run main.go

generate-profiles:
	cd src &&\
	go run  main.go -cpuProfiler -memProfiler &&\
	go tool pprof -dot .profiling/cpu.prof | dot -Tsvg > .profiling/cpu-graph.svg &&\
	go tool pprof -dot .profiling/mem.prof | dot -Tsvg > .profiling/mem-graph.svg
