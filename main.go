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
		hasNewVersion, version, err := deps.CheckNewVersionAvailable(dep)
		if err != nil {
			logs.ErrorLogger.Println(err)
		}

		actualBlobName, err := deps.EvalPlaceholders(dep.NamePattern, dep)
		if err != nil {
			logs.ErrorLogger.Println(err)
		}

		if hasNewVersion {
			logs.InfoLogger.Println("Needs bump: " + actualBlobName + ", From: " + dep.CurVersion + ", To: " +  version)
		} else {
			logs.InfoLogger.Println("Up-to-date: " + actualBlobName + ", Version: " + dep.CurVersion)
		}
	}
	
}
