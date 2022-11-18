package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v3"
)

type Dependabosh struct {
	Updates []Dependency `yaml:"updates"`
}

type Dependency struct {
	NamePattern string `yaml:"name_pattern"`
	Constraints string `yaml:"constraints"`
	CurVersion  string `yaml:"cur_version"`
	VersRegexp  string `yaml:"vers_regexp"`
	VersURL     string `yaml:"vers_url"`
	BlobURL     string `yaml:"blob_url"`
	SrcURL      string `yaml:"src_url"`
}

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	DebugLogger   *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "", 0) 
	DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger.SetOutput(ioutil.Discard)
	WarningLogger.SetOutput(ioutil.Discard)
}

func main() {
	yfilePath := os.Args[1]
	yfile, err := ioutil.ReadFile(yfilePath)

	if err != nil {
		panic(err)
	}
	var dependabosh Dependabosh
	err = yaml.Unmarshal(yfile, &dependabosh)
	if err != nil {
		panic(err)
	}
	for _, dep := range dependabosh.Updates {
		InfoLogger.Println("###################################")
		InfoLogger.Println("Name Pattern: " + dep.NamePattern)
		InfoLogger.Println("Cur Version: " + dep.CurVersion)
		InfoLogger.Println("Constraints: " + dep.Constraints)
		InfoLogger.Println("Version Regex: " + dep.VersRegexp)
		InfoLogger.Println("Versions URL: " + dep.VersURL)
		InfoLogger.Println("Blob URL: " + dep.BlobURL)
		InfoLogger.Println("SRC URL: " + dep.SrcURL)
		latestVersion, err := checkLatestVersion(dep)
		if err != nil {
			ErrorLogger.Println(err)
			continue
			InfoLogger.Println("###################################")
		}
		InfoLogger.Println("Latest Version: " + latestVersion)
		InfoLogger.Println("###################################")
	}
	
}

func findAllVersionsWithCaptureGroups(textContainingVersions string, regex *regexp.Regexp) []string {
	var versionCaptureGroup int

	versionCaptureGroup = 0
	for i, name := range regex.SubexpNames() {
		DebugLogger.Println("Capture Group " + name + " has index " + strconv.Itoa(i))
		if name == "version" {
			versionCaptureGroup = i
		}
	}


	allSubmatches := regex.FindAllStringSubmatch(textContainingVersions, -1)
	allVersions := make([]string, len(allSubmatches))
	for _, submatch := range allSubmatches {
		if submatch[versionCaptureGroup] != "" {
			allVersions = append(allVersions, submatch[versionCaptureGroup])
			DebugLogger.Println("Found version in input data: " + submatch[versionCaptureGroup] + ".")
		}
	}
	return allVersions
}

func checkLatestVersion(dep Dependency) (string, error) {
	vregexp := regexp.MustCompile(dep.VersRegexp)
	latestVersionsData := fetchLatestVersionsData(dep.VersURL)
	detectedVersions := findAllVersionsWithCaptureGroups(latestVersionsData, vregexp)

	_, err := semver.NewVersion(dep.CurVersion)
	if err != nil {
		ErrorLogger.Println("Vers " + dep.CurVersion + " can't be interpreted as a semver.")
		return "", semver.ErrInvalidSemVer
	}

	candidateVersion := dep.CurVersion
	for _, version := range detectedVersions {
		if version == "" {
			// FIXME
			continue
		}
		DebugLogger.Println("Comparing current candidate " + candidateVersion + " with " + version + ".")
		canonicalCandidate, err := semver.NewVersion(candidateVersion)
		if err != nil {
			ErrorLogger.Println("Vers " + candidateVersion + " can't be interpreted as a semver.")
			continue
		}
		canonicalVersion, err := semver.NewVersion(version)
		if err != nil {
			ErrorLogger.Println("Vers " + version + " can't be interpreted as a semver.")
			return "", semver.ErrInvalidSemVer
		}

		greaterThanCandidateConstraint, _ := semver.NewConstraint("> " + canonicalCandidate.String())
		versionInRangeConstraint, err := semver.NewConstraint(dep.Constraints)
		if err != nil {
			ErrorLogger.Println("Impossible to parse constraints: " + dep.Constraints)
		}

		if !greaterThanCandidateConstraint.Check(canonicalVersion) {
			DebugLogger.Println("Rejecting version " + version + " as it's older than current candidate " + candidateVersion)
		} else {
			if !versionInRangeConstraint.Check(canonicalVersion) {
				DebugLogger.Println("Rejecting version " + version + " despite being newer than current candidate " + candidateVersion + " because it doesn't pass constraints")
			} else {
				DebugLogger.Println("Adopting version " + version + " as it's newer than current candidate " + candidateVersion + " and passes constraints")
				candidateVersion = version
			}
		}
	}
	return candidateVersion, nil
}

func fetchLatestVersionsData(url string) string {
	resp, err := http.Get(url)
	if err != nil {
	   panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}

