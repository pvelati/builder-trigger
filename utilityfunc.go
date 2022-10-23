package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/pvelati/builder-trigger/aptutil"
)

func getLastReleaseInGitHubRepo(aptRepo string, filterRegexp string, compareFunctionIsLess func(v1 string, v2 string) bool) func() string {

	type oneReleaseItem struct {
		TagName string `json:"tag_name"`
	}

	return func() string {
		targetUrl := `https://api.github.com/repos/` + aptRepo + `/releases`

		log.Println(targetUrl)

		httpReq, err := http.NewRequest(http.MethodPost, targetUrl, nil)
		if err != nil {
			panic(err)
		}

		httpReq.Header.Set("Authorization", "Bearer "+githubToken())
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != 200 {
			panic(fmt.Errorf("invalid status %s", resp.Status))
		}

		var items []oneReleaseItem
		if err := json.Unmarshal(respBytes, &items); err != nil {
			panic(err)
		}

		log.Println(string(respBytes))

		return ""
	}
}

func getLastPveKernelInDestinationAptRepo(
	aptDownloader *aptutil.Downloader,
	codename string,
	arch string,
	variant string,
) func() string {
	return func() string {
		allPackages := aptDownloader.ParseIndexUrl("https://pvelati.github.io/apt-repository/debian/dists/" + codename + "/main/binary-" + arch + "/Packages")

		regExpr := `(pve-kernel-([^\s]*)-pve-` + regexp.QuoteMeta(variant) + `)`
		packageSelected := aptutil.GetOneDependsByRegex(
			allPackages["pve-kernel-5.15-"+variant],
			regExpr,
		)
		return regexp.MustCompile(regExpr).FindStringSubmatch(packageSelected)[2]
	}
}
func getLastPveKernel5_15(
	aptDownloader *aptutil.Downloader,
	codename string,
	arch string,
) func() string {
	return func() string {
		allPackages := aptDownloader.ParseIndexUrl("http://download.proxmox.com/debian/pve/dists/" + codename + "/pve-no-subscription/binary-" + arch + "/Packages")

		regExpr := `(pve-kernel-([^\s]*)-pve)`
		packageSelected := aptutil.GetOneDependsByRegex(
			allPackages["pve-kernel-5.15"],
			regExpr,
		)
		return regexp.MustCompile(regExpr).FindStringSubmatch(packageSelected)[2]
	}
}

// func getLastDebianMainKernel(
// 	aptDownloader *aptutil.Downloader,
// 	codename string,
// 	arch string,
// ) func() string {
// 	return func() string {
// 		allPackages := aptDownloader.ParseIndexUrl("http://deb.debian.org/debian/dists/" + codename + "/main/binary-" + arch + "/Packages.gz")

// 		regExpr := `(linux-image-([^\s]*)-` + arch + `)`
// 		packageSelected := aptutil.GetOneDependsByRegex(
// 			allPackages["linux-image-"+arch],
// 			regExpr,
// 		)
// 		return regexp.MustCompile(regExpr).FindStringSubmatch(packageSelected)[2]
// 	}
// }

// func getLastDebianMainCloudKernel(
// 	aptDownloader *aptutil.Downloader,
// 	codename string,
// 	arch string,
// ) func() string {
// 	return func() string {
// 		allPackages := aptDownloader.ParseIndexUrl("http://deb.debian.org/debian/dists/" + codename + "/main/binary-" + arch + "/Packages.gz")

// 		regExpr := `(linux-image-([^\s]*-cloud)-` + arch + `)`
// 		packageSelected := aptutil.GetOneDependsByRegex(
// 			allPackages["linux-image-cloud-"+arch],
// 			regExpr,
// 		)
// 		return regexp.MustCompile(regExpr).FindStringSubmatch(packageSelected)[2]
// 	}
// }
