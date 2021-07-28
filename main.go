package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	searchLine    = "Go"
	maxGoroutines = 10
)

var urls = []string{
	"https://golang.org/",
	"https://golang.org/",
}

func worker(jobs <-chan int, results chan<- int) {
	for j := range jobs {
		url := urls[j-1]

		if url == "" {
			log.Println("empty url")
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		count := bytes.Count(body, []byte(searchLine))
		resp.Body.Close()

		fmt.Printf("Count for %s: %d\n", url, count)

		results <- count
	}
}

func main() {
	if searchLine == "" {
		fmt.Println("empty search string")
		return
	}

	countGO := len(urls)

	jobs := make(chan int, countGO)
	results := make(chan int, countGO)

	for w := 1; w <= maxGoroutines; w++ {
		go worker(jobs, results)
	}

	for j := 1; j <= countGO; j++ {
		jobs <- j
	}
	close(jobs)

	var all int
	for a := 1; a <= countGO; a++ {
		all += <-results
	}
	fmt.Printf("Total: %d\n", all)
}
