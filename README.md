# Go-Graphics
Collection of computer graphics algorithms implemented purely in Go.

## Run
Requires an installation of go 1.21 or higher. All dependencies can be manually updated using the go get command:
```shell
go get .
```

The project provides a CLI to simplify rendering of scenes. Use the following command to show a help.
```shell
go run cmd/cli/main.go -h
```
A few demo scenes are provided in the `configs` directory and can be rendered like this;
```shell
go run cmd/cli/main.go -f config/bunny.json
```

## Compare
The implementation can be compared to the [pt](https://github.com/fogleman/pt) implementation. Select a scene in `cmd/compare/main.go` by removing the comment and run it using:
```shell
go run cmd/compare/main.go
```

## Benchmark
For the best benchmark result, use the `bench` package:
```shell
go install golang.design/x/bench@latest
cd benchmark
bench -name BenchmarkBunny
```
