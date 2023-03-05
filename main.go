package main

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

func main() {
	fmt.Printf("Init.\n")
	var wg sync.WaitGroup
	errChannel := make(chan error, 3)
	names := make([]*string, 3)
	namesClone := make([]*string, len(names))

	casa := "casa"
	bob := "bob"
	asa := "asa"

	names[0] = &casa
	names[1] = &asa
	names[2] = &bob

	// 3 goroutines to do some doWork
	createGoroutineToDoWork(&wg, &errChannel, 1, names)

	copy(namesClone, names)
	sort.Slice(namesClone, func(i, j int) bool {
		return *namesClone[i] < *namesClone[j]
	})
	createGoroutineToDoWork(&wg, &errChannel, 2, namesClone)
	createGoroutineToDoWork(&wg, &errChannel, 3, namesClone)

	waitFinishAllGoroutines(&wg, &errChannel)
	if result := <-errChannel; result != nil {
		fmt.Printf("Finish with error. %s \n", result.Error())
	}
	fmt.Printf("Finish.\n")

}

func waitFinishAllGoroutines(wg *sync.WaitGroup, errChannel *chan error) {
	fmt.Printf("Verify goroutines\n")
	wg.Wait()
	close(*errChannel)
	fmt.Printf("Finish all goroutines\n")

}

func createGoroutineToDoWork(wg *sync.WaitGroup, errChannel *chan error, id int, names []*string) {
	wg.Add(1)
	go func() {
		fmt.Printf("Init gourotines %d .\n", id)
		defer wg.Done()
		err := doWork(id, names)
		if err != nil {
			*errChannel <- fmt.Errorf("error from goroutine %d: %w", id, err)
		}
		fmt.Printf("Finish gourotines %d .\n", id)
	}()
}

func doWork(id int, names []*string) error {
	// do some work...
	fmt.Printf("Routine %d --- Names: %s, %s, %s.\n", id, *names[0], *names[1], *names[2])
	if id == 3 {
		return errors.New("test Error")
	}
	return nil
}
