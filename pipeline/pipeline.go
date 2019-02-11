package pipeline

import (
	"fmt"
	"sync"
)

func createStageOne(done chan struct{}, numOfJobs int) <-chan int {
	stageOneChannel := make(chan int)
	go func(nJobs int) {
		for i := 0; i < nJobs; i++ {
			select {
			// push the job details into a channel
			case stageOneChannel <- i:
			case <-done:
				fmt.Println("Graceful exit of Stage 1")
				return
			}
		}
		close(stageOneChannel)
	}(numOfJobs)
	return stageOneChannel
}

func createStageTwo(done chan struct{}, stageOneChannel <-chan int) <-chan int {
	stageTwoChannel := make(chan int)
	go func() {
		for n := range stageOneChannel {
			select {
			// Perform the operation and put the result into a channel
			case stageTwoChannel <- n + 20:
			case <-done:
				fmt.Println("Graceful exit of Stage 2")
				return
			}
		}
		close(stageTwoChannel)
	}()
	return stageTwoChannel
}

func createStageThree(done chan struct{}, stageTwoChannels []<-chan int) <-chan int {
	var wg sync.WaitGroup
	stageThreeChannel := make(chan int)
	doOperation := func(stageTwoChannel <-chan int) {
		for out := range stageTwoChannel {
			select {
			// read the value from second stage's channel to final stage
			case stageThreeChannel <- out:
			case <-done:
				fmt.Println("Graceful exit of Stage 3")
				return
			}
		}
		wg.Done()
	}

	wg.Add(len(stageTwoChannels))
	// Create go routines equal to number of channels created by stage 2
	for _, stageTwoChannel := range stageTwoChannels {
		go doOperation(stageTwoChannel)
	}

	go func() {
		// wait for all go routine to complete
		wg.Wait()
		// close the thord stage channel
		close(stageThreeChannel)
	}()
	return stageThreeChannel
}

// TestPipeline runs an iteration on pipeline
func TestPipeline() {

	done := make(chan struct{})
	defer close(done)

	numOfJobs := 10
	stageOneChannel := createStageOne(done, numOfJobs)

	var stageTwoChannels []<-chan int

	numOfParallelProcessors := 5
	for i := 0; i < numOfParallelProcessors; i++ {
		// Create go routine for each processing
		stageTwoChannels = append(stageTwoChannels, createStageTwo(done, stageOneChannel))
	}

	stageThreeChannel := createStageThree(done, stageTwoChannels)

	for o := range stageThreeChannel {
		fmt.Printf("%v ", o)
	}
	fmt.Println()
}
