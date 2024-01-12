package generation

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/operator-framework/operator-registry/alpha/template/composite"
	"github.com/operator-framework/operator-registry/pkg/image/containerdregistry"
)

type TapFitterCompositeTemplateReader struct {
	CompositePath string
	CatalogPath   string
	Provider      string
}

func (p *TapFitterCompositeTemplateReader) Ingest(c context.Context) ([]*GenerateDevfile, error) {

	cacheDir, err := os.MkdirTemp("", "tap-fitter-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(cacheDir)

	reg, err := containerdregistry.NewRegistry(
		containerdregistry.WithCacheDir(cacheDir),
	)
	if err != nil {
		return nil, err
	}
	defer reg.Destroy()

	// operator author's 'composite.yaml' file
	compositeReader, err := os.Open(p.CompositePath)
	if err != nil {
		return nil, fmt.Errorf("opening composite config file %q: %v", p.CompositePath, err)
	}
	defer compositeReader.Close()

	// catalog maintainer's 'catalogs.yaml' file
	tempCatalog, err := composite.FetchCatalogConfig(p.CatalogPath, http.DefaultClient)
	if err != nil {
		return nil, err
	}
	defer tempCatalog.Close()

	template := composite.NewTemplate(
		composite.WithCatalogFile(tempCatalog),
		composite.WithContributionFile(compositeReader),
		composite.WithRegistry(reg),
	)
	specs, err := template.Parse()
	if err != nil {
		return nil, err
	}

	generators := make([]*GenerateDevfile, 0)
	for _, catalog := range specs.CatalogSpec.Catalogs {
		workdir := catalog.Destination.WorkingDir
		// Ensure all paths exist with wxr-xr-x perms
		if err := os.MkdirAll(workdir, 0751); err != nil {
			return nil, err
		}
		writerFile, err := os.OpenFile(filepath.Join(workdir, "devfile.yaml"), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		// NOTE:  cannot defer closing the writer, since we delegate it to the caller

		generator := &GenerateDevfile{
			IndexDir:    "catalog",
			Name:        catalog.Name,
			BuildCTX:    catalog.Name,
			Provider:    p.Provider,
			Writer:      writerFile,
			CleanupFunc: writerFile.Close,
		}
		generators = append(generators, generator)
	}

	return generators, nil
}
