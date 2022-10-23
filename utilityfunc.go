package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/pvelati/builder-trigger/aptutil"
)

func getLastReleaseInGitHubRepo(aptRepo string, buildArch string, compareFunctionIsLess func(v1 string, v2 string) bool) func() string {

	type Tag struct {
		TagName string `json:"name"`
	}

	return func() string {
		targetUrl := `https://api.github.com/repos/` + aptRepo + `/tags`

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

		var tags []Tag
		if err := json.Unmarshal(respBytes, &tags); err != nil {
			panic(err)
		}

		var filterTags []Tag
		for _, tag := range tags {
			if contain := strings.Contains(tag.TagName, buildArch); contain {
				filterTags = append(filterTags, tag)
			}
		}
		if len(filterTags) == 0 {
			fmt.Printf("no tags found, new release")
		}

		sort.SliceStable(filterTags, func(i, j int) bool {
			return filterTags[i].TagName > filterTags[j].TagName
		})

		log.Println(string(respBytes))

		// target_tag := filterTags[0].TagName

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
