package main

import (
	"strings"

	arg "github.com/alexflint/go-arg"
	"github.com/stakater/RequirementsUpdater/pkg/log"
	"github.com/stakater/RequirementsUpdater/pkg/requirements"
)

const (
	RequirementsFile = "requirements.yaml"
)

var (
	logger = log.New()
)

type Args struct {
	Path          string
	ChartName     string `arg:"required"`
	Version       string `arg:"required"`
	RepositoryURL string `arg:"required"`
	Alias         string
}

func main() {
	var args Args
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

	req, err := requirements.Read(args.Path)
	if err != nil {
		logger.Error(err)
		return
	}

	replaced := false

	for index := 0; index < len(req.Dependencies); index++ {
		if req.Dependencies[index].Name == args.ChartName &&
			req.Dependencies[index].Repository == args.RepositoryURL &&
			req.Dependencies[index].Alias == args.Alias {
			req.Dependencies[index].Version = args.Version
			replaced = true
			break
		}
	}

	if !replaced {
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
