package structured_test

import (
	"log"
	"os"

	structured "github.com/poy/go-structured-writer"
)

func ExampleWithTime() {
	stucturedWriter := structured.New(os.Stdout,
		structured.WithTimestamp(),
		structured.WithCallSite(),
	)
	log.SetOutput(stucturedWriter)

	log.Println("Some helpful tracing")
}
