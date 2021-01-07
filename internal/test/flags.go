package test

import (
	"flag"
)

var (
	HelmChartPath   string
	HelmChartValues string
	DexConfigPath   string
)

func init() {
	flag.StringVar(&HelmChartPath, "helm-chart-path", "", "define path to local helm chart to use for tests")
	flag.StringVar(&HelmChartValues, "helm-chart-values", "", "define path to local helm chart values file to use for tests")
	flag.StringVar(&DexConfigPath, "dex-conf-path", "", "define path to local dex config file to use for tests")
}
