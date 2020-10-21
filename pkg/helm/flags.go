package helm

import (
	"flag"
)

var (
	KindCluster     string
	HelmChartPath   string
	HelmChartValues string
)

func init() {
	flag.StringVar(&KindCluster, "kind-cluster", "", "define pre-existing cluster to use for tests")
	flag.StringVar(&HelmChartPath, "helm-chart-path", "", "define path to local helm chart to use for tests")
	flag.StringVar(&HelmChartValues, "helm-chart-values", "", "define path to local helm chart values file to use for tests")
}
