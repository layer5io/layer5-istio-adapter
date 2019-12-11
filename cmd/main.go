package main

import (
	"fmt"
	"os"

	"istio.io/istio/mixer/adapter/layer5"
)

func main() {
	addr := ""
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	s, err := layer5.NewLayer5Adapter(addr)
	if err != nil {
		fmt.Printf("unable to start server: %v", err)
		os.Exit(-1)
	}

	shutdown := make(chan error, 1)
	go func() {
		s.Run(shutdown)
	}()
	_ = <-shutdown
}
