package main

import (
	"fmt"

	"github.com/oceanc80/tap-fitter/pkg/generation"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newRootCmd() (*cobra.Command, error) {
	var (
		compositePath string
		catalogPath   string
	)
	var rootCmd = &cobra.Command{
		Short: "tap-fitter",
		Long:  `tap-fitter reads a composite template to prepare a repository for a catalog production pipeline`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if catalogPath == "" || compositePath == "" {
				return fmt.Errorf("both 'catalog-path' and 'composite-path' flags are required")
			}
			p := generation.TapFitterCompositeTemplateReader{
				CompositePath: compositePath,
				CatalogPath:   catalogPath,
			}
			generators, err := p.Ingest(cmd.Context())
			if err != nil {
				return err
			}
			for _, g := range generators {

				if err := g.Generate(); err != nil {
					return err
				}
			}

			return nil
		},
	}
	f := rootCmd.Flags()
	f.StringVar(&compositePath, "composite-path", "", "the path to the composite template used for configuration (required if with-composite-template is set)")
	f.StringVar(&catalogPath, "catalog-path", "", "the path/URL to the catalog template used for configuration (required if with-composite-template is set)")

	return rootCmd, nil
}

func main() {
	cmd, err := newRootCmd()
	if err != nil {
		logrus.Panic(err)
	}
	if err := cmd.Execute(); err != nil {
		logrus.Panic(err)
	}
}
