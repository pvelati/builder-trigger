package main

import (
	"log"
)

func executeTasks(
	tasks map[string]TaskInfoType,
	filterByName string,
	filterByTags []string,
) {
	if filterByName != "" {
		log.Println("Execute only task with name:", filterByName)
	}
	if len(filterByTags) > 0 {
		log.Println("Include only task with tags:", filterByTags)
	}

	for taskName, taskInfo := range tasks {
		queueTask := true
		if filterByName != "" {
			queueTask = queueTask && taskName == filterByName
		}
		if len(filterByTags) > 0 {
			tagFound := false
			for _, t1 := range taskInfo.Tags {
				for _, t2 := range filterByTags {
					tagFound = tagFound || t1 == t2
				}
			}
			queueTask = queueTask && tagFound
		}

		if queueTask {
			log.Println("running task name: " + taskName)
			isRunnig := taskInfo.IsTaskRunning()
			if !isRunnig {
				log.Println("checking")
				// targetVersion := taskInfo.ObtainLastTargetVersion()

				// currentVersion := taskInfo.ObtainLastReleasedVersion()

				// log.Println("targetVersion: " + targetVersion)
				// log.Println("currentVersion: " + currentVersion)

				// if currentVersion != targetVersion {
				//      taskInfo.VersionChangeNotify(targetVersion)
				// }

				// tasks[taskName] = taskInfo
			}
		} else {
			log.Println("skipping task name: " + taskName)
		}
	}
}
