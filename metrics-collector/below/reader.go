package below

import "os/exec"

//func readEthtoolMetrics() string {
//	cmd := exec.Command("below", "dump", "ethtool-queue", "-b", "1 min ago", "--detail", "-O", "openmetrics")
//	out, err := cmd.Output()
//	if err != nil {
//		panic(err)
//	}
//	ethtoolMetrics := string(out)
//	return ethtoolMetrics
//}

func readNetworkMetrics() []string {
	cmd := exec.Command("below", "dump", "iface", "-b", "1 min ago", "--detail", "-O", "openmetrics")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	ifaceMetrics := string(out)
	return []string{ifaceMetrics}
}

func ReadMetrics() []string {
	return readNetworkMetrics()
}
