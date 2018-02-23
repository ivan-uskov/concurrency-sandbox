package main

import (
	"net/http"
	"fmt"
	"sync"
)

type FetchResult struct {
	url string
	status int
}

type Fetcher func() FetchResult

func getFetcher(url string) Fetcher {
	return func() FetchResult {
		resp, err := http.Head(url)
		status := http.StatusBadGateway
		if err == nil {
			status = resp.StatusCode
		}

		return FetchResult{url, status}
	}
}

func runPool(jobs <-chan Fetcher, results chan<- FetchResult, size int) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(size)

	for i := 0; i < size; i++ {
		go func(wg *sync.WaitGroup) {
			for j := range jobs {
				results <- j()
			}
			wg.Done()
		}(wg)
	}

	return wg
}

func fetchWithPool(urls []string, poolSize int) <-chan FetchResult {
	jobs := make(chan Fetcher, poolSize)
	results := make(chan FetchResult, poolSize)
	wg := runPool(jobs, results, poolSize)

	go func() {
		for _, url := range urls {
			jobs <- getFetcher(url)
		}
		close(jobs)

		wg.Wait()
		close(results)
	}()
	return results
}

func main() {
	urls := []string{
		`https://ispringsolutions.com`,
		`https://ispringonline.ru`,
		`https://ispring.ru`,
		`https://ispringsolutions.com`,
		`https://ispringonline.ru`,
		`https://ispringlearn.com`,
		`https://ispring.ru`,
		`https://ispringcloud.ru`,
		`https://ispringcloud.com`,
		`https://blog.ispring.ru`,
	}

	for res := range fetchWithPool(urls, 3) {
		fmt.Printf("Server: %s, status: %d\n", res.url, res.status)
	}
}