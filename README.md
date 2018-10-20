# go-structured-writer
[![GoDoc][go-doc-badge]][go-doc] [![travis][travis-badge]][travis]

Structured logging via standard log

```go
package main

import (
	"log"
	"os"

	structured "github.com/poy/go-structured-writer"
)

func main() {
	// Structured logs are slow and should be able to be off by default.
	if os.Getenv("TRACE_LOGS") == "true"{
		stucturedWriter := structured.New(os.Stdout,
			// Include timestamp
			structured.WithTimestamp(),

			// Include call site information (FileName:LineNumber)
			structured.WithCallSite(),

			// Custom name and value. This example is "level:debug".
			structured.WithFieldFunc("level", func([]byte) (interface{}, error) {
				return "debug", nil
			}),
		)

		// Remove extra flags for timestamp as the timestamp field has it.
		log.SetFlags(0)

		// Set the structured writer as the log writer. Anywhere the global log
		// is used, structured logs will be emitted.
		log.SetOutput(stucturedWriter)
	}

	// Write some logs
	log.Println("Some helpful tracing")
	log.Println("Some other helpful tracing")
	log.Println("wow... more helpful tracing")
}
```

if `TRACE_LOGS` is set to `true`, then the output will be:
```
{"callsite":"test/main.go:35","level":"debug","msg":"Some helpful tracing","timestamp":1513055977608015185}
{"callsite":"test/main.go:36","level":"debug","msg":"Some other helpful tracing","timestamp":1513055977608093398}
{"callsite":"test/main.go:37","level":"debug","msg":"wow... more helpful tracing","timestamp":1513055977608118570}
```

Otherwise:
```
2017/12/11 21:55:47 Some helpful tracing
2017/12/11 21:55:47 Some other helpful tracing
2017/12/11 21:55:47 wow... more helpful tracing
```


[go-doc-badge]:             https://godoc.org/github.com/poy/go-structured-writer?status.svg
[go-doc]:                   https://godoc.org/github.com/poy/go-structured-writer
[travis-badge]:             https://travis-ci.org/poy/go-structured-writer.svg?branch=master
[travis]:                   https://travis-ci.org/poy/go-structured-writer?branch=master
