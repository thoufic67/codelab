package main

//
import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

var MAXINT int64 = 10000000
var THREAD_LIMIT = 10
var currentNumber int64 = 2
var primes int32 = 0

func checkPrimeNumber(num int64) bool {
	if num < 2 {
		return false
	}
	sq_root := int64(math.Sqrt(float64(num)))
	var i int64 = 2

	for i = 2; i <= sq_root; i++ {
		if num%i == 0 {
			return false
		}
	}

	return true
}

func work(thread int, start time.Time, wg *sync.WaitGroup) {
	for {
		x := atomic.AddInt64(&currentNumber, 1)
		if x > MAXINT {
			break
		}
		if checkPrimeNumber(x) {
			atomic.AddInt32(&primes, 1)
		}
	}
	fmt.Println("Thread: ", thread, " took ", time.Since(start))
	defer wg.Done()
}

func fairThreadedCheck() {
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i <= THREAD_LIMIT; i++ {
		wg.Add(1)
		go work(i, time.Now(), &wg)
	}

	wg.Wait()
	fmt.Println("In Fair Threaded Check number of primes:", primes, " and took: ", time.Since(start))
}

func regularCheck() {
	start := time.Now()
	var regularCheckedPrime int32 = 0
	var i int64
	for i = 2; i <= MAXINT; i++ {
		if checkPrimeNumber(i) {
			atomic.AddInt32(&regularCheckedPrime, 1)
		}
	}
	fmt.Println("In Regular Check number of primes:", regularCheckedPrime, " and took: ", time.Since(start))
}

func main() {
	fairThreadedCheck()
	regularCheck()
}
