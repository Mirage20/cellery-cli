package helm

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/manifest"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/renderutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
}

func Some() {
	c, err := GetHelmArchive("https://storage.googleapis.com/cellery-io/helm/cellery-runtime-0.1.0.tgz")
	if err != nil {
		log.Fatal(err)
	}
	str := `
mysql:
  enabled: false

controller:
  enabled: false

apim:
  enabled: false

idp:
  enabled: false

observability-portal:
    enabled: false

sp-worker:
    enabled: false
`

	m, err := renderutil.Render(c, &chart.Config{Raw: str}, renderutil.Options{})
	if err != nil {
		log.Fatal(err)
	}
	manifests := manifest.SplitManifests(m)
	for _, v := range manifests {
		fmt.Println(v.Name)
		fmt.Println(v.Content)
		fmt.Println("---")
	}
}

type Render struct {
	Enabled bool `json:"enabled"`
}

type Config struct {
	MySql               Render `json:"mysql"`
	Controller          Render `json:"controller"`
	ApiManager          Render `json:"apim"`
	Idp                 Render `json:"idp"`
	ObservabilityPortal Render `json:"observability-portal"`
	SpWorker            Render `json:"sp-worker"`
}

func GenerateManifest(cfg Config) []manifest.Manifest {
	b, err := json.Marshal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	c, err := GetHelmArchive("https://storage.googleapis.com/cellery-io/helm/cellery-runtime-0.1.0.tgz")
	if err != nil {
		log.Fatal(err)
	}

	m, err := renderutil.Render(c, &chart.Config{Raw: string(b)}, renderutil.Options{})
	if err != nil {
		log.Fatal(err)
	}
	manifests := manifest.SplitManifests(m)
	for _, v := range manifests {
		fmt.Println(v.Name)
		fmt.Println(v.Content)
		fmt.Println("---")
	}
	return manifests
}

// Returns the Helm chart archive located at the given URI (can be either an http(s) address or a file path)
func GetHelmArchive(chartArchiveUri string) (*chart.Chart, error) {

	// Download chart archive
	chartFile, err := GetResource(chartArchiveUri)
	if err != nil {
		return nil, err
	}
	//noinspection GoUnhandledErrorResult
	defer chartFile.Close()

	// Check chart requirements to make sure all dependencies are present in /charts
	helmChart, err := chartutil.LoadArchive(chartFile)
	if err != nil {
		return nil, errors.Wrapf(err, "loading chart archive")
	}
	return helmChart, err
}

func GetResource(uri string) (io.ReadCloser, error) {
	var file io.ReadCloser
	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		resp, err := http.Get(uri)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, errors.Errorf("http GET returned status %d", resp.StatusCode)
		}

		file = resp.Body
	} else {
		path, err := filepath.Abs(uri)
		if err != nil {
			return nil, errors.Wrapf(err, "getting absolute path for %v", uri)
		}

		f, err := os.Open(path)
		if err != nil {
			return nil, errors.Wrapf(err, "opening file %v", path)
		}
		file = f
	}

	// Write the body to file
	return file, nil
}
