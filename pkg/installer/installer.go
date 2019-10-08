package installer

import (
	"encoding/json"
	"fmt"
	"io"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/manifest"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/renderutil"
	"strings"
)

type Installer struct {
	Cellery   *chart.Chart
	IstioInit *chart.Chart
	Istio     *chart.Chart
	Writer    io.Writer
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

func (i *Installer) InstallCellery(cfg Config) error {
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	m, err := renderutil.Render(i.Cellery, &chart.Config{Raw: string(b)}, renderutil.Options{})
	if err != nil {
		return err
	}
	manifests := manifest.SplitManifests(m)
	for _, v := range SortByKind(manifests, InstallOrder) {
		fmt.Fprintf(i.Writer, "---\n# Source: %s\n", v.Name)
		fmt.Fprintln(i.Writer, v.Content)
	}
	return nil
}

func (i *Installer) UninstallCellery(cfg Config) error {
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	m, err := renderutil.Render(i.Cellery, &chart.Config{Raw: string(b)}, renderutil.Options{})
	if err != nil {
		return err
	}
	manifests := manifest.SplitManifests(m)
	for _, v := range SortByKind(manifests, UninstallOrder) {
		fmt.Fprintf(i.Writer, "---\n# Source: %s\n", v.Name)
		fmt.Fprintln(i.Writer, v.Content)
	}
	return nil
}

func (i *Installer) InstallIstio() error {
	for _, f := range i.IstioInit.Files {
		if strings.HasSuffix(f.TypeUrl, ".yaml") || strings.HasSuffix(f.TypeUrl, ".yaml") {
			fmt.Fprintf(i.Writer, "---\n# Source: %s\n", f.TypeUrl)
			fmt.Fprintln(i.Writer, string(f.Value))
		}
	}

	//b, err := json.Marshal(cfg)
	//if err != nil {
	//	return err
	//}

	m, err := renderutil.Render(i.Istio, &chart.Config{Raw: "{}"}, renderutil.Options{
		ReleaseOptions: chartutil.ReleaseOptions{
			Name:      "do",
			Namespace: "istio-system",
		},
	})
	if err != nil {
		return err
	}
	manifests := manifest.SplitManifests(m)
	for _, v := range SortByKind(manifests, InstallOrder) {
		if strings.HasSuffix(v.Name, ".yaml") || strings.HasSuffix(v.Name, ".yaml") {
			fmt.Fprintf(i.Writer, "---\n# Source: %s\n", v.Name)
			fmt.Fprintln(i.Writer, v.Content)
		}
	}
	return nil
}
