package main

import (
	"concurrentHTTP/boundedparallelism/pipeline"
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	fmt.Println("hello world")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Running TestPipeline")
		pipeline.TestPipeline()
	})

	done := make(chan struct{})
	go func() {
		http.ListenAndServe(":6060", nil)
	}()
	<-done
}
