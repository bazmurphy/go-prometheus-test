package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	log.Println("starting client...")

	client := &http.Client{}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-signals:
				log.Println("signal received, stopping client...")
				return
			default:
				makeRequest(client)
			}
		}
	}()

	wg.Wait()
}

func makeRequest(client *http.Client) {
	request, err := http.NewRequest("GET", "http://server:8080/", nil) // internal docker network (not localhost)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	log.Println("making request...")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer response.Body.Close()

	log.Println("response received. status:", response.Status)
}
