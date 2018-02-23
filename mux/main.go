package main

import (
	"net/http"
	"sync"
	"fmt"
)

type FetchResult struct {
	url string
	status int
}

func fetch(url string) <-chan FetchResult {
	ch := make(chan FetchResult)
	go func() {
		defer close(ch)
		resp, err := http.Head(url)
		status := http.StatusBadGateway
		if err == nil {
			status = resp.StatusCode
		}

		ch <- FetchResult{url, status}
	}()

	return ch
}

func merge(cs ...<-chan FetchResult) <-chan FetchResult {
	var wg sync.WaitGroup
	out := make(chan FetchResult)

	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan FetchResult) {
			for n := range c {
				out <- n
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	resultChan := merge(
		fetch(`https://ispringsolutions.com`),
		fetch(`https://ispringonline.ru`),
		fetch(`https://ispringlearn.com`),
		fetch(`https://ispring.ru`),
		fetch(`https://ispringcloud.ru`),
		fetch(`https://ispringcloud.com`),
		fetch(`https://blog.ispring.ru`),
	)

	for res := range resultChan {
		fmt.Printf("Server: %s, status: %d\n", res.url, res.status)
	}
}
