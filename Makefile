run:
	cd src && go run main.go

generate-profiles:
	make run &&\
	go tool pprof -dot src/.profiling/cpu.prof | dot -Tsvg > src/.profiling/cpu-graph.svg &&\
	go tool pprof -dot src/.profiling/mem.prof | dot -Tsvg > src/.profiling/mem-graph.svg
