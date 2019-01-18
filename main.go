package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/irotoris/jobkickqd/jobkickqd"
)

func main() {
	i := 0
	for {
		fmt.Println("Start a command...")
		jobID := "sampleJob" + strconv.Itoa(i)
		//cmdString := "for i in {1..3}; do sleep 1; echo " + jobID + ":  `date \"+%Y-%m-%d %H:%M:%S\"`; done"
		cmdString := "echo " + jobID + ":  `date \"+%Y-%m-%d %H:%M:%S\"`"
		envs := []string{"ENV=dev", "EDITOR=vim"}
		timeoutDuration := 60 * time.Second
		j := jobkickqd.NewJob(jobID, cmdString, envs, timeoutDuration)
		ctx := context.Background()

		errChan := make(chan error, 1)
		go func() {
			err := j.Execute(ctx)
			fmt.Println("Log: " + j.ExecutionLog)
			errChan <- err
		}()

		err := <-errChan
		if err != nil {
			fmt.Print("err:", err)
		}
		fmt.Println("Finished a command.")

		time.Sleep(1 * time.Second)
		i++
		if i > 0 {
			break
		}
	}
}
