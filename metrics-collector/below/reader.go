package below

import "os/exec"

func readEthtoolMetrics() string {
	cmd := exec.Command("below", "dump", "ethtool-queue", "-b", "1 min ago", "--detail", "-O", "openmetrics")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	ethtool_metrics := string(out)
	return ethtool_metrics
}

func readNetworkMetrics() []string {
	cmd := exec.Command("below", "dump", "iface", "-b", "1 min ago", "--detail", "-O", "openmetrics")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	iface_metrics := string(out)
	return []string{iface_metrics, readEthtoolMetrics()}
}

func ReadMetrics() []string {
	return readNetworkMetrics()
}
