package main

import (
	"os"

	"github.com/codegangsta/cli"
	logging "github.com/op/go-logging"
)

func main() {
	app := cli.NewApp()
	app.Name = "mdp"
	app.Usage = "Generate organized html documentation from markdown files with a toml preamble"
	var outDir, inDir string
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "out",
			Value:       "./public",
			Usage:       "Directory to output the generated html pages",
			Destination: &outDir,
		},
		cli.StringFlag{
			Name:        "in",
			Value:       ".",
			Usage:       "Directory to read the project inside.",
			Destination: &inDir,
		},
	}

	app.Action = func(c *cli.Context) {
		log := logging.MustGetLogger("mdp")
		log.Infof("Parsing directory %v", inDir)
		node, err := ParseTopDirectory(inDir)
		if err != nil {
			log.Fatalf("Error parsing directory: %v", err)
		}
		// Create output directory if it doesn't exist
		if _, err = os.Stat(outDir); err != nil {
			if os.IsNotExist(err) {
				os.Mkdir(outDir, 0744)
			} else {
				log.Fatalf("Unknown error inspecting output directory: %v", err)
			}
		}
		err = RenderNode(outDir, node)
		if err != nil {
			log.Fatalf("Error rendering output: %v", err)
		}
	}

	app.Run(os.Args)
}
