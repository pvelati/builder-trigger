package main

import "github.com/pvelati/builder-trigger/aptutil"

func buildTasksList() map[string]TaskInfoType {
	aptDownloader := &aptutil.Downloader{}

	tasks := map[string]TaskInfoType{}

	for _, codename := range []string{"bullseye"} {
		for _, arch := range []string{"amd64"} {
			tasks["pve-kernel5-15"] = TaskInfoType{
				IsTaskRunning: func() bool {
					return false
				},
				ObtainLastTargetVersion:   getLastPveKernel5_15(aptDownloader, codename, arch),
				ObtainLastReleasedVersion: getLastReleaseInGitHubRepo("pvelati/kernel-builder-pve", "123", func(v1, v2 string) bool { return false }),
				VersionChangeNotify:       makeDefaultWebhook("pvelati/kernel-builder-pve", "bullseye", "amd64"),

				Tags: []string{"apt"},
			}
			tasks["pve-kernel5-15-apt"] = TaskInfoType{
				IsTaskRunning: func() bool {
					return false
				},
				ObtainLastTargetVersion:   getLastReleaseInGitHubRepo("pvelati/kernel-builder-pve", "123", func(v1, v2 string) bool { return false }),
				ObtainLastReleasedVersion: getLastPveKernelInDestinationAptRepo(aptDownloader, codename, arch, "silvermont"),
				VersionChangeNotify:       makeDefaultWebhook("pvelati/apt-repository", "bullseye", "amd64"),

				Tags: []string{"repository"},
			}
		}
	}

	// NOT USED FOR NOW
	// for _, codename := range []string{"bullseye"} {
	// 	for _, arch := range []string{"amd64"} {
	// 		tasks["kernel-"+codename+"-"+arch] = TaskInfoType{
	// 			ObtainLastTargetVersion: getLastDebianMainKernel(aptDownloader, codename, arch),
	// 			VersionChangeNotify: makeDefaultWebhook("pvelati/kernel-builder-pve", codename, arch),
	// 		}
	// 		tasks["kernel-"+codename+"-"+arch+"-cloud"] = TaskInfoType{
	// 			ObtainLastTargetVersion: getLastDebianMainCloudKernel(aptDownloader, codename, arch),
	// 			VersionChangeNotify: makeDefaultWebhook("pvelati/github-actions-sandbox", codename, arch),
	// 		}
	// 	}
	// }
	return tasks
}
