package main

import (
	"context"
	"fmt"

	"./job"
)

func main() {
	fmt.Println("Start a command...")
	jobID := "sampleJob"
	cmdString := "date \"+%Y-%m-%d\""
	envs := []string{"ENV=dev", "EDITOR=vim"}
	attr := make(map[string]string)
	job := job.NewJob(jobID, cmdString, envs, attr)
	ctx := context.Background()
	job.Execute(ctx)
	fmt.Println("Finished a command.")
}
