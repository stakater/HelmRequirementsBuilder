package main

import (
	"io/ioutil"
	"os"
	"strings"

	arg "github.com/alexflint/go-arg"
	"github.com/stakater/HelmRequirementsBuilder/pkg/log"
	"github.com/stakater/HelmRequirementsBuilder/pkg/requirements"
)

const (
	RequirementsFile = "requirements.yaml"
)

var (
	logger = log.New()
	args   = Args{}
)

type Args struct {
	Path              string
	ChartName         string `arg:"required"`
	Version           string `arg:"required"`
	RepositoryURL     string `arg:"required"`
	Alias             string
	CreateIfNotExists bool
}

func main() {

	// Create requirements.yaml if it doesn't exist
	if _, err := os.Stat(args.Path); os.IsNotExist(err) && args.CreateIfNotExists {
		if err := ioutil.WriteFile(args.Path, nil, 0644); err != nil {
			logger.Error(err)
		}
	}

	req, err := requirements.Read(args.Path)
	if err != nil {
		logger.Error(err)
		return
	}

	index, exists := getDependencyIndexForArgs(req.Dependencies)

	if exists {
		req.Dependencies[index].Version = args.Version
	} else {
		// Add new entry
		req.Dependencies = append(req.Dependencies, requirements.Dependency{
			Alias:      args.Alias,
			Name:       args.ChartName,
			Repository: args.RepositoryURL,
			Version:    args.Version,
		})
	}

	err = requirements.Write(args.Path, req)
	if err != nil {
		logger.Error(err)
	}

}

func init() {
	parseArgs()
}

func parseArgs() {
	// Default values
	args.CreateIfNotExists = false

	arg.MustParse(&args)

	// Populate alias
	if args.Alias == "" {
		args.Alias = strings.ToLower(args.ChartName)
	}

	if args.Path == "" {
		args.Path = "./" + RequirementsFile
	} else {
		args.Path = args.Path + "/" + RequirementsFile
	}
}

func getDependencyIndexForArgs(dependencies []requirements.Dependency) (int, bool) {
	for index := 0; index < len(dependencies); index++ {
		if dependencies[index].Name == args.ChartName &&
			dependencies[index].Repository == args.RepositoryURL &&
			dependencies[index].Alias == args.Alias {
			return index, true
		}
	}

	return -1, false
}
