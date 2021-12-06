package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, ch chan Tweet, wg *sync.WaitGroup) (tweets []*Tweet) {
	defer wg.Done()
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(ch)
			break
		}
		ch <- *tweet
	}
	return
}

func consumer(ch chan Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range ch {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
			continue
		}
		fmt.Println(t.Username, "\tdoes not tweet about golang")
	}

}

func main() {
	start := time.Now()
	stream := GetMockStream()

	ch := make(chan Tweet)

	// Producer
	wg := sync.WaitGroup{}
	wg.Add(2)
	go producer(stream, ch, &wg)
	// Consumer
	go consumer(ch, &wg)
	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
