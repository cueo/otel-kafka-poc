package otel

import "testing"

var openmetrics_txt = `# TYPE ethtool_queue_rx_bytes_per_sec_bytes_per_second gauge
	# UNIT ethtool_queue_rx_bytes_per_sec_bytes_per_second bytes_per_second
	ethtool_queue_rx_bytes_per_sec_bytes_per_second{hostname="ip-172-31-24-129",interface="ens5",queue="0"} 15 1692917890
	# TYPE ethtool_queue_tx_bytes_per_sec_bytes_per_second gauge
	# UNIT ethtool_queue_tx_bytes_per_sec_bytes_per_second bytes_per_second
	ethtool_queue_tx_bytes_per_sec_bytes_per_second{hostname="ip-172-31-24-129",interface="ens5",queue="0"} 0 1692917890
	# TYPE ethtool_queue_rx_count_per_sec gauge
	ethtool_queue_rx_count_per_sec{hostname="ip-172-31-24-129",interface="ens5",queue="0"} 0 1692917890
	# TYPE ethtool_queue_tx_count_per_sec gauge
	ethtool_queue_tx_count_per_sec{hostname="ip-172-31-24-129",interface="ens5",queue="0"} 0 1692917890
	# TYPE ethtool_queue_tx_missed_tx counter
	ethtool_queue_tx_missed_tx{hostname="ip-172-31-24-129",interface="ens5",queue="0"} 0 1692917890
	# TYPE ethtool_queue_tx_unmask_interrupt counter
	ethtool_queue_tx_unmask_interrupt{hostname="ip-172-31-24-129",interface="ens5",queue="0"} 277601 1692917890
	# TYPE ethtool_queue_rx_bytes_per_sec_bytes_per_second gauge
	# UNIT ethtool_queue_rx_bytes_per_sec_bytes_per_second bytes_per_second
	ethtool_queue_rx_bytes_per_sec_bytes_per_second{hostname="ip-172-31-24-129",interface="ens5",queue="1"} 0 1692917890
	# TYPE ethtool_queue_tx_bytes_per_sec_bytes_per_second gauge
	# UNIT ethtool_queue_tx_bytes_per_sec_bytes_per_second bytes_per_second
	ethtool_queue_tx_bytes_per_sec_bytes_per_second{hostname="ip-172-31-24-129",interface="ens5",queue="1"} 18 1692917890
	# TYPE ethtool_queue_rx_count_per_sec gauge
	ethtool_queue_rx_count_per_sec{hostname="ip-172-31-24-129",interface="ens5",queue="1"} 0 1692917890
	# TYPE ethtool_queue_tx_count_per_sec gauge
	ethtool_queue_tx_count_per_sec{hostname="ip-172-31-24-129",interface="ens5",queue="1"} 0 1692917890
	# TYPE ethtool_queue_tx_missed_tx counter
	ethtool_queue_tx_missed_tx{hostname="ip-172-31-24-129",interface="ens5",queue="1"} 0 1692917890
	# TYPE ethtool_queue_tx_unmask_interrupt counter
	ethtool_queue_tx_unmask_interrupt{hostname="ip-172-31-24-129",interface="ens5",queue="1"} 568982 1692917890
# EOF`

func TestParser(t *testing.T) {
	metrics := parseOpenMetrics(openmetrics_txt)

	if len(metrics) != 12 {
		t.Fatalf("Expected 12 metrics, got %d", len(metrics))
	}
	if metrics[0].Name != "ethtool_queue_rx_bytes_per_sec_bytes_per_second" {
		t.Fatalf("Expected metric name to be ethtool_queue_rx_bytes_per_sec_bytes_per_second, got %s", metrics[0].Name)
	}
	if metrics[0].Value != 15.0 {
		t.Fatalf("Expected metric value to be 15.0, got %f", metrics[0].Value)
	}
	if metrics[0].Labels["hostname"] != "ip-172-31-24-129" {
		t.Errorf("Expected metric hostname to be ip-172-31-24-129, got %s", metrics[0].Labels["hostname"])
	}
	if metrics[0].Labels["interface"] != "ens5" {
		t.Errorf("Expected metric interface to be ens5, got %s", metrics[0].Labels["interface"])
	}
	if metrics[0].Labels["queue"] != "0" {
		t.Errorf("Expected metric queue to be 0, got %s", metrics[0].Labels["queue"])
	}
}

func TestBelowCollector(t *testing.T) {
	_, err := collectBelowMetrics()
	if err != nil {
		t.Fatalf("Error collecting below metrics: %s", err)
	}
}
