package main

import (
	"fmt"

	"github.com/release-engineering/tap-fitter/pkg/generation"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	compositePath string
	catalogPath   string
	provider      string
)

func newRootCmd() (*cobra.Command, error) {
	var rootCmd = &cobra.Command{
		Use:   "tap-fitter",
		Short: "tap-fitter takes composite templates and outputs corresponding devfiles",
		Long:  `tap-fitter reads a composite template and outputs corresponding devfiles to prepare a repository for a catalog production pipeline.
Note: because templates may use relative paths to related artifacts, tap-fitter should be run in the same location as input templates.`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if catalogPath == "" || compositePath == "" {
				return fmt.Errorf("both 'catalog-path' and 'composite-path' flags are required")
			}
			p := generation.TapFitterCompositeTemplateReader{
				CompositePath: compositePath,
				CatalogPath:   catalogPath,
				Provider:      provider,
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
	rootCmd.Flags().StringVarP(&compositePath, "composite-path", "t", "", "[REQUIRED] the path to the composite template used for configuration")
	rootCmd.Flags().StringVarP(&catalogPath, "catalog-path", "c", "", "[REQUIRED] the path/URL to the catalog template used for configuration")
	rootCmd.Flags().StringVarP(&provider, "provider", "p", "tap-fitter", "the provider of the catalog")

	return rootCmd, nil
}

func main() {
	cmd, err := newRootCmd()
	if err != nil {
		logrus.Fatalf("Error constructing tap-fitter command: ", err)
	}
	if err := cmd.Execute(); err != nil {
		logrus.Fatalf("Error executing tap-fitter command: ", err)
	}
}
