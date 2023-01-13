package cli

import (
	"flag"
	"io/ioutil"
	"gopkg.in/yaml.v3"

	"github.com/fejnartal/dependabosh/deps"
)

var depsYaml string
var repoPath string

func init() {
	flag.StringVar(&depsYaml, "deps", "", "[mandatory] - dependabosh.yml where dependencies are defined")
	flag.StringVar(&repoPath, "repo", "", "path to the repository where dependencies need to be bumped")
}

func ParseArgs() (deps.Dependabosh, error) {
	flag.Parse()

	yfilePath := depsYaml
	yfile, err := ioutil.ReadFile(yfilePath)

	if err != nil {
		return deps.Dependabosh{}, err
	}
	var dependabosh deps.Dependabosh
	err = yaml.Unmarshal(yfile, &dependabosh)
	if err != nil {
		return deps.Dependabosh{}, err
	}

	return dependabosh, nil
}

func PrintDefaults() {
	flag.PrintDefaults()
}
