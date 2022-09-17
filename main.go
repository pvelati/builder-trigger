package main

import (
	"github.com/pvelati/builder-trigger/aptutil"
	"regexp"
	"time"
)

func main() {
	aptDownloader := &aptutil.Downloader{}

	tasks := map[string]TaskInfoType{}

	tasks["pve-kernel5-15"] = TaskInfoType{
		LastUpstreamVersion: getLastPveKernel5_15(aptDownloader),
		CheckInterval:       10 * time.Minute,
		VersionChangeNotify: makeDefaultWebhook("pvelati/github-actions-sandbox", "bullseye", "amd64"),
	}

	for _, codename := range []string{"bullseye"} {
		for _, arch := range []string{"amd64"} {
			tasks["kernel-"+codename+"-"+arch] = TaskInfoType{
				LastUpstreamVersion: getLastDebianMainKernel(aptDownloader, codename, arch),
				CheckInterval:       10 * time.Minute,
				VersionChangeNotify: makeDefaultWebhook("pvelati/github-actions-sandbox", codename, arch),
			}
			tasks["kernel-"+codename+"-"+arch+"-cloud"] = TaskInfoType{
				LastUpstreamVersion: getLastDebianMainCloudKernel(aptDownloader, codename, arch),
				CheckInterval:       10 * time.Minute,
				VersionChangeNotify: makeDefaultWebhook("pvelati/github-actions-sandbox", codename, arch),
			}
		}
	}
	executeTasks(tasks)
}

func getLastPveKernel5_15(aptDownloader *aptutil.Downloader) func() string {
	return func() string {
		allPackages := aptDownloader.ParseIndexUrl("http://download.proxmox.com/debian/pve/dists/bullseye/pve-no-subscription/binary-amd64/Packages")

		regExpr := `(pve-kernel-([^\s]*)-pve)`
		packageSelected := aptutil.GetOneDependsByRegex(
			allPackages["pve-kernel-5.15"],
			regExpr,
		)
		return regexp.MustCompile(regExpr).FindStringSubmatch(packageSelected)[2]
	}
}

func getLastDebianMainKernel(
	aptDownloader *aptutil.Downloader,
	codename string,
	arch string,
) func() string {
	return func() string {
		allPackages := aptDownloader.ParseIndexUrl("http://deb.debian.org/debian/dists/" + codename + "/main/binary-" + arch + "/Packages.gz")

		regExpr := `(linux-image-([^\s]*)-` + arch + `)`
		packageSelected := aptutil.GetOneDependsByRegex(
			allPackages["linux-image-"+arch],
			regExpr,
		)
		return regexp.MustCompile(regExpr).FindStringSubmatch(packageSelected)[2]
	}
}

func getLastDebianMainCloudKernel(
	aptDownloader *aptutil.Downloader,
	codename string,
	arch string,
) func() string {
	return func() string {
		allPackages := aptDownloader.ParseIndexUrl("http://deb.debian.org/debian/dists/" + codename + "/main/binary-" + arch + "/Packages.gz")

		regExpr := `(linux-image-([^\s]*-cloud)-` + arch + `)`
		packageSelected := aptutil.GetOneDependsByRegex(
			allPackages["linux-image-cloud-"+arch],
			regExpr,
		)
		return regexp.MustCompile(regExpr).FindStringSubmatch(packageSelected)[2]
	}
}
