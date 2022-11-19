run:
	go run main.go

generate-profiles:
	go run main.go &&\
	go tool pprof -dot .profiling/cpu.prof | dot -Tsvg > .profiling/cpu-graph.svg &&\
	go tool pprof -dot .profiling/mem.prof | dot -Tsvg > .profiling/mem-graph.svg