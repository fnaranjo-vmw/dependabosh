package main

import (
	"github.com/fejnartal/dependabosh/cli"
	"github.com/fejnartal/dependabosh/deps"
	"github.com/fejnartal/dependabosh/logs"
)

func main() {
	dependabosh, err := cli.ParseArgs()
	if err != nil {
		logs.ErrorLogger.Println(err)
		cli.PrintDefaults()
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
