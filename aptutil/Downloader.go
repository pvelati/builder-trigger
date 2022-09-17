package aptutil

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/mod/semver"
)

type PackagesIndex map[string](map[string]string)

type Downloader struct {
	cache map[string](PackagesIndex)
	lock  sync.Mutex
}

func (downloader *Downloader) ParseIndexUrl(indexUrl string) PackagesIndex {
	downloader.lock.Lock()
	defer downloader.lock.Unlock()

	if len(downloader.cache) == 0 {
		downloader.cache = make(map[string]PackagesIndex)
	} else if len(downloader.cache[indexUrl]) > 0 {
		return downloader.cache[indexUrl]
	}

	log.Println("fetching " + indexUrl)

	resp, err := http.DefaultClient.Get(indexUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var downloadedIndexContent string
	switch strings.ToLower(resp.Header.Get("content-type")) {
	case "application/x-gzip":
		log.Println("reading gunzip content")
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			panic(err)
		}
		defer gzReader.Close()
		downloadedIndexContentBytes, err := ioutil.ReadAll(gzReader)
		if err != nil {
			panic(err)
		}
		downloadedIndexContent = string(downloadedIndexContentBytes)
	default:
		log.Println("reading content")
		downloadedIndexContentBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		downloadedIndexContent = string(downloadedIndexContentBytes)
	}

	allPackages := PackagesIndex{}
	{
		currentPackageLinees := []string{}

		for strings.Contains(downloadedIndexContent, "\n") {
			nextN := strings.Index(downloadedIndexContent, "\n")
			nextLine := downloadedIndexContent[0:nextN]
			downloadedIndexContent = downloadedIndexContent[nextN+1:]
			if nextLine == "" {
				packageName, packageVars := processLinees(currentPackageLinees)
				if len(allPackages[packageName]) != 0 {

					vPackageVars := canonicalize(packageVars["Version"])
					vCurrentlyInMap := canonicalize(allPackages[packageName]["Version"])

					if semver.Compare(vPackageVars, vCurrentlyInMap) > 0 {
						allPackages[packageName] = packageVars
					}
				} else {
					allPackages[packageName] = packageVars
				}
				currentPackageLinees = []string{}
			} else {
				currentPackageLinees = append(currentPackageLinees, nextLine)
			}
		}
	}
	downloader.cache[indexUrl] = allPackages
	return allPackages
}

func canonicalize(v string) string {
	if semver.IsValid(v) {
		return v
	}
	regVersion := regexp.MustCompile(`^(\d+)\.(\d+)(\.\d+|)(.*)`)

	matches := regVersion.FindStringSubmatch(v)
	if len(matches) == 0 {
		panic(fmt.Errorf("no matches for"))
	}

	newTmpVersion := "v"

	if major := matches[1]; major != "" {
		newTmpVersion += major
	} else {
		newTmpVersion += "0"
	}
	newTmpVersion += "."

	if minor := matches[2]; minor != "" {
		newTmpVersion += minor
	} else {
		newTmpVersion += "0"
	}

	if patch := matches[3]; patch != "" {
		newTmpVersion += patch
	} else {
		newTmpVersion += ".0"
	}
	// newTmpVersion += matches[4]

	if !semver.IsValid(newTmpVersion) {
		panic(fmt.Errorf("invalid tmp version %#v", newTmpVersion))
	}

	return semver.Canonical(newTmpVersion)
}
