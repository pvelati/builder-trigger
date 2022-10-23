package main

type TaskInfoType struct {
	IsTaskRunning             func() bool
	ObtainLastTargetVersion   func() string
	ObtainLastReleasedVersion func() string
	VersionChangeNotify       func(version string)

	Tags []string
}
