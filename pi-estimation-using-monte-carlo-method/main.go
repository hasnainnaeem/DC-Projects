package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type MonteCarloPi struct {
	totalNumberOfTosses int64
	numberInCircle      int64
	threadCount         int
	segment             int64
}

func main() {
	var piEstimate float64

	pi := &MonteCarloPi{}
	pi.getUserInput()
	pi.calculatePi()

	piEstimate = 4 * float64(pi.numberInCircle) / float64(pi.totalNumberOfTosses)
	fmt.Printf("Estimated pi = %f\n", piEstimate)
}

func (pi *MonteCarloPi) getUserInput() {
	fmt.Print("Enter the number of threads and the total number of tosses: ")
	_, err := fmt.Scanf("%d %d", &pi.threadCount, &pi.totalNumberOfTosses)
	if err != nil {
		fmt.Println("Invalid input. Please enter valid numbers.")
	}
}

func (pi *MonteCarloPi) calculatePi() {
	rand.Seed(time.Now().UnixNano())

	pi.segment = pi.totalNumberOfTosses / int64(pi.threadCount)

	var wg sync.WaitGroup

	for thread := 0; thread < pi.threadCount; thread++ {
		wg.Add(1)
		go pi.threadToss(thread, &wg)
	}

	wg.Wait()
}

func (pi *MonteCarloPi) threadToss(rank int, wg *sync.WaitGroup) {
	var x, y, distanceSquared float64
	var numberInCircleInThread int64

	defer wg.Done()

	for toss := int64(0); toss < pi.segment; toss++ {
		x = randDouble()
		y = randDouble()
		distanceSquared = x*x + y*y
		if distanceSquared <= 1 {
			numberInCircleInThread++
		}
	}

	// Use atomic operation to update numberInCircle
	atomic.AddInt64(&pi.numberInCircle, numberInCircleInThread)
}

// Returns a random number in the range -1 to 1
func randDouble() float64 {
	return 2 * (rand.Float64() - 0.5)
}
