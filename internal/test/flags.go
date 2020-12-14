package test

import (
	"flag"
)

var (
	KindCluster     string
	HelmChartPath   string
	HelmChartValues string
	WithKind        bool
	DexConfigPath   string
)

func init() {
	flag.BoolVar(&WithKind, "with-kind", false, "define wether to use kind cluster for tests")
	flag.StringVar(&KindCluster, "kind-cluster", "", "define pre-existing cluster to use for tests")
	flag.StringVar(&HelmChartPath, "helm-chart-path", "", "define path to local helm chart to use for tests")
	flag.StringVar(&HelmChartValues, "helm-chart-values", "", "define path to local helm chart values file to use for tests")
	flag.StringVar(&DexConfigPath, "dex-conf-path", "", "define path to local dex config file to use for tests")
}
