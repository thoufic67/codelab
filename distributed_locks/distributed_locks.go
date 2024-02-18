package main

//
import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	instances             = make([]int, 5)
	minRequiredLock int32 = 3
	FREE                  = 0
	connections           = 50
	wg              sync.WaitGroup
	mu              sync.Mutex
	recursed        int32
)

func acquireDistributedLock(connection int) {
	var acquired int32 = 0
	for instance := 0; instance < len(instances); instance++ {
		mu.Lock()
		if instances[instance] == FREE {
			instances[instance] = connection
		}
		mu.Unlock()
	}

	for instance := 0; instance < len(instances); instance++ {
		if instances[instance] == connection {
			acquired++
		}
	}

	if acquired >= minRequiredLock {
		time.Sleep(15 * time.Nanosecond)
		// fmt.Println(acquired, connection, instances)
		releaseDistributedLock(connection, false)
		return
	} else {
		atomic.AddInt32(&recursed, 1)
		releaseDistributedLock(connection, true)
	}
}

func releaseDistributedLock(connection int, relock bool) {
	for instance := 0; instance < len(instances); instance++ {
		mu.Lock()
		if instances[instance] == connection {
			instances[instance] = FREE
		}
		mu.Unlock()
	}
	if relock {
		defer acquireDistributedLock(connection)
	}
}

func operation(connection int) {
	start := time.Now()
	defer wg.Done()
	acquireDistributedLock(connection)
	fmt.Println("Connection ", connection, " Completed in :", time.Since(start))
}

func main() {
	start := time.Now()

	for connection := 1; connection <= connections; connection++ {
		wg.Add(1)
		go operation(connection)
	}

	wg.Wait()
	fmt.Println("Recursed: ", recursed, "Total time : ", time.Since(start))
}
