package main

import "time"

type TaskInfoType struct {
	CurrentVersion string
	CheckTime      *time.Time

	CheckInterval       time.Duration        `json:"-"`
	LastUpstreamVersion func() string        `json:"-"`
	VersionChangeNotify func(version string) `json:"-"`
}
