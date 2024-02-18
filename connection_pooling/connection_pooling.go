package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	connectionLimit = 5
	q               []Connection
	requests        = 100
	wg              sync.WaitGroup
	mu              sync.Mutex
)

type Connection struct {
	connectionId string
}

func releaseConnection(connection Connection) {
	mu.Lock()
	q = append([]Connection{connection}, q...)
	mu.Unlock()
}

func getConnection() (Connection, bool) {
	mu.Lock()
	defer mu.Unlock()
	if len(q) == 0 {
		return Connection{}, false
	}
	conn := q[0]
	q = q[1:]

	return conn, true
}

func makeConnections() {
	for i := 1; i <= connectionLimit; i++ {
		connection := Connection{connectionId: fmt.Sprintf("connection%d", i)}
		q = append(q, connection)
	}
}

func doOperation(request int) {
	defer wg.Done()
	var connection Connection
	var available bool
	start := time.Now()

	for available != true {
		connection, available = getConnection()
	}

	time.Sleep(5 * time.Millisecond)
	fmt.Println("Request ", request, " acquired: ", connection.connectionId, " in: ", time.Since(start))
	releaseConnection(connection)

}

func main() {
	start := time.Now()
	makeConnections()
	fmt.Println(q)
	for req := 0; req < requests; req++ {
		wg.Add(1)
		go doOperation(req)

	}
	wg.Wait()
	fmt.Println("Completed in :", time.Since(start))
}
