package pipeline

import (
	"fmt"
	"sync"
)

func createStageOne(done chan struct{}, numOfJobs int) <-chan int {
	stageOneChannel := make(chan int)
	go func(nJobs int) {
		fmt.Println("STAGE 1: Push input numbers 1 to 10 to Stage 1 Channel")
		for i := 0; i < nJobs; i++ {
			fmt.Printf("STAGE 1: Pushing input %v to Stage 1 channel \n", i)
			select {
			// push the job details into a channel
			case stageOneChannel <- i:
			case <-done:
				fmt.Println("STAGE 1: Graceful exit of Stage 1")
				return
			}
		}
		close(stageOneChannel)
	}(numOfJobs)
	return stageOneChannel
}

func createStageTwo(go_id int, done chan struct{}, stageOneChannel <-chan int) <-chan int {
	stageTwoChannel := make(chan int)
	go func() {
		for n := range stageOneChannel {
			fmt.Printf("STAGE 2: Performed operation on input %v by go routine %v \n", n, go_id)
			select {
			// Perform the operation and put the result into a channel
			case stageTwoChannel <- n + 20:
			case <-done:
				fmt.Println("STAGE 2: Graceful exit of Stage 2")
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
				fmt.Println("STAGE 3: Graceful exit of Stage 3")
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
	fmt.Printf("Running %v Go routines to perform operation\n", numOfParallelProcessors)
	for i := 0; i < numOfParallelProcessors; i++ {
		// Create go routine for each processing
		stageTwoChannels = append(stageTwoChannels, createStageTwo(i, done, stageOneChannel))
	}

	stageThreeChannel := createStageThree(done, stageTwoChannels)

	for o := range stageThreeChannel {
		fmt.Printf("Result after STAGE 3 = %v on input %v \n", o, o-20)
	}
	fmt.Println()
}
