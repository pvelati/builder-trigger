package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func executeTasks(tasks map[string]TaskInfoType) {
	lastExecutionJsonFilename := "last-execution.json"

	previousTasks := map[string]TaskInfoType{}

	{
		previousBytes, err := ioutil.ReadFile(lastExecutionJsonFilename)
		if err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}
		} else {
			log.Println("loading " + lastExecutionJsonFilename)
			if err := json.Unmarshal(previousBytes, &previousTasks); err != nil {
				panic(err)
			}
		}
	}

	for taskName, taskInfo := range tasks {
		log.Println("task name: " + taskName)
		runCheck := true
		if previousTasks[taskName].CheckTime != nil {
			lastCheckAge := time.Since(*previousTasks[taskName].CheckTime)
			log.Println("last check age:", lastCheckAge.Truncate(time.Second))
			runCheck = lastCheckAge > taskInfo.CheckInterval
		}
		if runCheck {
			log.Println("checking")
			taskInfo.CurrentVersion = taskInfo.LastUpstreamVersion()
			t := time.Now()
			taskInfo.CheckTime = &t

			if previousTasks[taskName].CurrentVersion != taskInfo.CurrentVersion {
				taskInfo.VersionChangeNotify(taskInfo.CurrentVersion)
			}

			tasks[taskName] = taskInfo
		} else {
			log.Println("skipping check")
			tasks[taskName] = previousTasks[taskName]
		}
	}

	{
		previousBytes, err := json.MarshalIndent(tasks, "", "\t")
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(lastExecutionJsonFilename, previousBytes, os.FileMode(0o644)); err != nil {
			panic(err)
		}
	}

}
