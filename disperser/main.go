package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var counter int
var mu sync.Mutex

func sendRequest(addr string, count int) {
	jsonStr := fmt.Sprintf("{\"count\":%d}", count)

	buffer := bytes.NewBuffer([]byte(jsonStr))

	res, err := http.Post("http://"+addr, "application/json", buffer)

	if err != nil {
		log.Fatalf("Error when sending request %d to %s: %s", count, addr, err.Error())
	}

	if res.StatusCode == 200 {
		log.Printf("Success when sending request %d to %s", count, addr)
	}
}

func main() {
	var nodeAddrs []string

	addrsEnv := os.Getenv("NODE_ADDRS")

	if addrsEnv != "" {
		nodeAddrs = strings.Split(addrsEnv, " ")
	}

	for _, addr := range nodeAddrs {
		go func(addr string) {

			ticker := time.NewTicker(time.Second)

			for range ticker.C {
				mu.Lock()
				currentCounter := counter
				counter++
				mu.Unlock()

				sendRequest(addr, currentCounter)
			}
		}(addr)
	}

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, os.Interrupt)

	<-sigCh
}
