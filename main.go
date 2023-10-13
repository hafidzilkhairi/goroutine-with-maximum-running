package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// i've created several versions to provide you with more options
	fmt.Println("======start of v1=====")
	v1()
	fmt.Println("======end of v1=====")
	fmt.Println("======start of v2=====")
	v2()
	fmt.Println("======end of v2=====")
}

const (
	maxData        = 11
	maxConcurrency = 3
)

func v2() {
	var (
		concurrencyRunningKeeper = make(chan struct{}, maxConcurrency)
		resultChan               = make(chan int)
		result                   []int
		wg                       = sync.WaitGroup{}
	)

	go func() {
		for resultData := range resultChan {
			wg.Done()
			fmt.Printf("Receiving the value %d	at %s\n", resultData, time.Now().Format(time.DateTime))
			result = append(result, resultData)
		}
	}()

	for i := 0; i < maxData; i++ {
		wg.Add(1)
		go func(idx int) {
			defer func() {
				// consume the keeper to allow the next goroutine to run
				<-concurrencyRunningKeeper
			}()
			// sending a value to the keeper to signal that a goroutine is starting
			// This line will block and wait for its value to be consumed
			concurrencyRunningKeeper <- struct{}{}

			time.Sleep(1 * time.Second)
			resultChan <- idx
		}(i)
	}

	// wg.Wait() wait for all goroutines to complete
	wg.Wait()
	close(resultChan)

	fmt.Print("Finally, the results are:")
	for _, val := range result {
		fmt.Print(" ", val)
	}
	fmt.Printf("\nLength of the result: %d\n", len(result))
}

func v1() {

	var (
		done                     = make(chan struct{})
		concurrencyRunningKeeper = make(chan struct{}, maxConcurrency)
		resultChan               = make(chan int)
		result                   []int
	)

	go func() {
		for i := 0; i < maxData; i++ {
			// waiting for any goroutine to complete
			<-done

			// consume the keeper to allow the next goroutine to run
			<-concurrencyRunningKeeper
		}
		// closing resultChan to indicate that all goroutines have completed
		close(resultChan)
	}()

	for i := 0; i < maxData; i++ {
		go func(idx int) {
			defer func() {
				// sending a value to the 'done' channel to signal that a goroutine has completed
				done <- struct{}{}
			}()
			// sending a value to the keeper to signal that a goroutine is starting
			// This line will block and wait for its value to be consumed
			concurrencyRunningKeeper <- struct{}{}

			time.Sleep(1 * time.Second)
			resultChan <- idx
		}(i)
	}

	// This 'for' loop will iterate over 'resultChan' until it is closed
	for resultData := range resultChan {
		fmt.Printf("Receiving the value %d	at %s\n", resultData, time.Now().Format(time.DateTime))
		result = append(result, resultData)
	}

	fmt.Print("Finally, the results are:")
	for _, val := range result {
		fmt.Print(" ", val)
	}
	fmt.Printf("\nLength of the result: %d\n", len(result))
}
