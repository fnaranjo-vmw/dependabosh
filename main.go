package main

import (
	"io/ioutil"
	"os"
	"gopkg.in/yaml.v3"

	"github.com/fejnartal/dependabosh/deps"
	"github.com/fejnartal/dependabosh/logs"
)

func main() {
	yfilePath := os.Args[1]
	yfile, err := ioutil.ReadFile(yfilePath)

	if err != nil {
		panic(err)
	}
	var dependabosh deps.Dependabosh
	err = yaml.Unmarshal(yfile, &dependabosh)
	if err != nil {
		panic(err)
	}
	for _, dep := range dependabosh.Updates {
		logs.InfoLogger.Println("###################################")
		logs.InfoLogger.Println("Name Pattern: " + dep.NamePattern)
		logs.InfoLogger.Println("Cur Version: " + dep.CurVersion)
		logs.InfoLogger.Println("Constraints: " + dep.Constraints)
		logs.InfoLogger.Println("Version Regex: " + dep.VersRegexp)
		logs.InfoLogger.Println("Versions URL: " + dep.VersURL)
		logs.InfoLogger.Println("Blob URL: " + dep.BlobURL)
		logs.InfoLogger.Println("SRC URL: " + dep.SrcURL)
		latestVersion, err := deps.CheckLatestVersion(dep)
		if err != nil {
			logs.ErrorLogger.Println(err)
			continue
			logs.InfoLogger.Println("###################################")
		}
		logs.InfoLogger.Println("Latest Version: " + latestVersion)
		logs.InfoLogger.Println("###################################")
	}
	
}
