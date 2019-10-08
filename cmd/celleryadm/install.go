package main

import (
	"bytes"
	"github.com/cellery-io/cellery/pkg/client/helm"
	"github.com/cellery-io/cellery/pkg/installer"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
)

func newInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "install cellery runtime components in your kubernetes cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			c1, err := helm.GetHelmArchive("https://storage.googleapis.com/cellery-io/helm/cellery-runtime-0.1.0.tgz")
			c2, err := helm.GetHelmArchive("https://storage.googleapis.com/cellery-io/helm/istio-init-1.2.2.tgz")
			c3, err := helm.GetHelmArchive("https://storage.googleapis.com/cellery-io/helm/istio-1.2.2.tgz")
			if err != nil {
				return err
			}
			h := installer.Config{}
			h.Controller.Enabled = true
			var buf bytes.Buffer

			i := &installer.Installer{
				Cellery:   c1,
				IstioInit: c2,
				Istio:     c3,
				Writer:    &buf,
			}
			i.InstallIstio()

			KubectlApply(&buf)
			return nil
		},
	}
}

func KubectlApply(r io.Reader) error {
	cmd := exec.Command("kubectl", "apply", "--dry-run", "-f", "-")
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
