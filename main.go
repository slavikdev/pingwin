package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	f, err := os.Open("targets.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	targets := make([]string, 0)
	for scanner.Scan() {
		targets = append(targets, scanner.Text())
	}

	var wg sync.WaitGroup
	wg.Add(len(targets))

	fmt.Println("Running load test...")
	for i := 0; i < len(targets); i++ {
		go func(id int) {
			defer wg.Done()
			client := NewClient(id, targets[id])
			client.Run()
		}(i)
	}

	wg.Wait()
	fmt.Println("Finished load test.")
}
