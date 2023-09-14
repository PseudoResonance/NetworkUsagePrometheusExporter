package main

import (
	"flag"
	"fmt"
	"time"
	"os/exec"
	"bufio"
	"log"
	"strings"
	"regexp"

	"net/http"
)

var cacheTime time.Time
var cache string
var cacheDuration time.Duration

var interfaceName *regexp.Regexp

func main() {
	var cacheExpiry int
	var host string
	var port int

	flag.IntVar(&cacheExpiry, "t", 30, "Cache timeout period in seconds (Default 30).")
	flag.StringVar(&host, "h", "", "Webserver host (Default empty).")
	flag.IntVar(&port, "p", 15835, "Webserver port (Default 15835).")

	flag.Parse()
	cacheDuration, _ = time.ParseDuration(fmt.Sprintf("%ds", cacheExpiry))

	interfaceName = regexp.MustCompile("^[0-9]+: (.*?): .*?state (.*?) .*$")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<!DOCTYPE html><title>Network Usage Exporter</title><h1>Network Usage Exporter</h1><p><a href=\"/metrics\">Metrics</a></p>")
	})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Error fetching metrics:", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		fmt.Fprintf(w, GetCache())
	})
	StartServer(host, port)
}

func StartServer(host string, port int) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Error with HTTP server:", err)
			StartServer(host, port)
		}
	}()
	
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		panic(err)
	}
}

func GetCache() string {
	if (time.Now().Sub(cacheTime) > cacheDuration) {
		UpdateCache()
	}
	return cache
}

func UpdateCache() {
	cache = DataToString(GetData())
	cacheTime = time.Now()
}

func DataToString(data map[string]map[string]string) string {
	output := ""

	output += "# HELP interface_rx_bytes_total RX bytes sent over the interface" + "\n"
	output += "# TYPE interface_rx_bytes_total counter" + "\n"
	for key, _ := range data {
		output += "interface_rx_bytes_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["rx_bytes"] + "\n"
	}
	output += "# HELP interface_rx_packets_total RX packets sent over the interface" + "\n"
	output += "# TYPE interface_rx_packets_total counter" + "\n"
	for key, _ := range data {
		output += "interface_rx_packets_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["rx_packets"] + "\n"
	}
	output += "# HELP interface_rx_errors_total RX packet errors on the interface" + "\n"
	output += "# TYPE interface_rx_errors_total counter" + "\n"
	for key, _ := range data {
		output += "interface_rx_errors_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["rx_errors"] + "\n"
	}
	output += "# HELP interface_rx_dropped_total RX packets dropped on the interface" + "\n"
	output += "# TYPE interface_rx_dropped_total counter" + "\n"
	for key, _ := range data {
		output += "interface_rx_dropped_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["rx_dropped"] + "\n"
	}
	output += "# HELP interface_rx_overrun_total RX packets overrun on the interface" + "\n"
	output += "# TYPE interface_rx_overrun_total counter" + "\n"
	for key, _ := range data {
		output += "interface_rx_overrun_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["rx_overrun"] + "\n"
	}
	output += "# HELP interface_rx_mcast_total RX mcast packets sent over the interface" + "\n"
	output += "# TYPE interface_rx_mcast_total counter" + "\n"
	for key, _ := range data {
		output += "interface_rx_mcast_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["rx_mcast"] + "\n"
	}
	
	output += "# HELP interface_tx_bytes_total TX bytes sent over the interface" + "\n"
	output += "# TYPE interface_tx_bytes_total counter" + "\n"
	for key, _ := range data {
		output += "interface_tx_bytes_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["tx_bytes"] + "\n"
	}
	output += "# HELP interface_tx_packets_total TX packets sent over the interface" + "\n"
	output += "# TYPE interface_tx_packets_total counter" + "\n"
	for key, _ := range data {
		output += "interface_tx_packets_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["tx_packets"] + "\n"
	}
	output += "# HELP interface_tx_errors_total TX packet errors on the interface" + "\n"
	output += "# TYPE interface_tx_errors_total counter" + "\n"
	for key, _ := range data {
		output += "interface_tx_errors_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["tx_errors"] + "\n"
	}
	output += "# HELP interface_tx_dropped_total TX packets dropped on the interface" + "\n"
	output += "# TYPE interface_tx_dropped_total counter" + "\n"
	for key, _ := range data {
		output += "interface_tx_dropped_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["tx_dropped"] + "\n"
	}
	output += "# HELP interface_tx_carrier_total TX packets with lost carrier on the interface" + "\n"
	output += "# TYPE interface_tx_carrier_total counter" + "\n"
	for key, _ := range data {
		output += "interface_tx_carrier_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["tx_carrier"] + "\n"
	}
	output += "# HELP interface_tx_collsns_total TX packets collisions on the interface" + "\n"
	output += "# TYPE interface_tx_collsns_total counter" + "\n"
	for key, _ := range data {
		output += "interface_tx_collsns_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["tx_collsns"] + "\n"
	}

	output += "# HELP interface_rx_errors_length_total RX packets with invalid length on the interface" + "\n"
	output += "# TYPE interface_rx_errors_length_total counter" + "\n"
	for key, _ := range data {
		output += "interface_rx_errors_length_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["rx_errors_length"] + "\n"
	}
	output += "# HELP interface_rx_errors_crc_total RX packets with invalid CRC on the interface" + "\n"
	output += "# TYPE interface_rx_errors_crc_total counter" + "\n"
	for key, _ := range data {
		output += "interface_rx_errors_crc_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["rx_errors_crc"] + "\n"
	}
	output += "# HELP interface_rx_errors_frame_total RX packets with invalid frame alignment on the interface" + "\n"
	output += "# TYPE interface_rx_errors_frame_total counter" + "\n"
	for key, _ := range data {
		output += "interface_rx_errors_frame_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["rx_errors_frame"] + "\n"
	}
	output += "# HELP interface_rx_errors_fifo_total RX packets dropped because of FIFO errors on the interface" + "\n"
	output += "# TYPE interface_rx_errors_fifo_total counter" + "\n"
	for key, _ := range data {
		output += "interface_rx_errors_fifo_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["rx_errors_fifo"] + "\n"
	}
	output += "# HELP interface_rx_errors_missed_total RX packets missed on the interface" + "\n"
	output += "# TYPE interface_rx_errors_missed_total counter" + "\n"
	for key, _ := range data {
		output += "interface_rx_errors_missed_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["rx_errors_missed"] + "\n"
	}

	output += "# HELP interface_tx_errors_aborted_total TX packets aborted on the interface" + "\n"
	output += "# TYPE interface_tx_errors_aborted_total counter" + "\n"
	for key, _ := range data {
		output += "interface_tx_errors_aborted_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["tx_errors_aborted"] + "\n"
	}
	output += "# HELP interface_tx_errors_fifo_total TX packets dropped because of FIFO errors on the interface" + "\n"
	output += "# TYPE interface_tx_errors_fifo_total counter" + "\n"
	for key, _ := range data {
		output += "interface_tx_errors_fifo_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["tx_errors_fifo"] + "\n"
	}
	output += "# HELP interface_tx_errors_window_total TX packets with frame collisions on the interface" + "\n"
	output += "# TYPE interface_tx_errors_window_total counter" + "\n"
	for key, _ := range data {
		output += "interface_tx_errors_window_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["tx_errors_window"] + "\n"
	}
	output += "# HELP interface_tx_errors_heartbeat_total TX packets dropped because of heartbeat errors on the interface" + "\n"
	output += "# TYPE interface_tx_errors_heartbeat_total counter" + "\n"
	for key, _ := range data {
		output += "interface_tx_errors_heartbeat_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["tx_errors_heartbeat"] + "\n"
	}
	output += "# HELP interface_tx_errors_transns_total TX transns packets on the interface" + "\n"
	output += "# TYPE interface_tx_errors_transns_total counter" + "\n"
	for key, _ := range data {
		output += "interface_tx_errors_transns_total{interface=\"" + key + "\",up=\"" + data[key]["up"] + "\"} " + data[key]["tx_errors_transns"] + "\n"
	}

	return output
}

func GetData() map[string]map[string]string {
	// Return data
	var data = make(map[string]map[string]string)
	// Run IP command to get data
	out, err := exec.Command("ip", "-s", "-s", "link", "show").Output()
	if err != nil {
		// Panic if unable to get output of command
		panic(err)
	}
	// Convert output into string and read line by line
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	var curInter string // Current interface name
	var nextRx, nextRxErrors, nextTx, nextTxErrors bool // If next line is RX or TX
	for scanner.Scan() {
		line := scanner.Text()
		if (len(line) > 0) {
			if (!strings.HasPrefix(line, " ")) { // If line doesn't start with whitespace, it is a header line (interface name)
				// Reset current interface name
				curInter = ""
				curInterUp := "0"
				// Match for interface name and up status
				matches := interfaceName.FindStringSubmatch(line)
				curInter = strings.Split(matches[1], "@")[0] // Set interface name and strip anything after @
				state := matches[2]
				if (state != "DOWN") { // If state isn't down, assume it is up
					curInterUp = "1"
				}
				data[curInter] = map[string]string {
					"up": curInterUp,
				}
			} else if (nextRx) {
				// Split line by spaces into individual values
				vals := strings.Fields(line)
				data[curInter]["rx_bytes"] = vals[0]
				data[curInter]["rx_packets"] = vals[1]
				data[curInter]["rx_errors"] = vals[2]
				data[curInter]["rx_dropped"] = vals[3]
				data[curInter]["rx_overrun"] = vals[4]
				data[curInter]["rx_mcast"] = vals[5]
				nextRx = false
			} else if (nextTx) {
				// Split line by spaces into individual values
				vals := strings.Fields(line)
				data[curInter]["tx_bytes"] = vals[0]
				data[curInter]["tx_packets"] = vals[1]
				data[curInter]["tx_errors"] = vals[2]
				data[curInter]["tx_dropped"] = vals[3]
				data[curInter]["tx_carrier"] = vals[4]
				data[curInter]["tx_collsns"] = vals[5]
				nextTx = false
			} else if (nextRxErrors) {
				// Split line by spaces into individual values
				vals := strings.Fields(line)
				data[curInter]["rx_errors_length"] = vals[0]
				data[curInter]["rx_errors_crc"] = vals[1]
				data[curInter]["rx_errors_frame"] = vals[2]
				data[curInter]["rx_errors_fifo"] = vals[3]
				data[curInter]["rx_errors_missed"] = vals[4]
				nextRxErrors = false
			} else if (nextTxErrors) {
				// Split line by spaces into individual values
				vals := strings.Fields(line)
				data[curInter]["tx_errors_aborted"] = vals[0]
				data[curInter]["tx_errors_fifo"] = vals[1]
				data[curInter]["tx_errors_window"] = vals[2]
				data[curInter]["tx_errors_heartbeat"] = vals[3]
				data[curInter]["tx_errors_transns"] = vals[4]
				nextTxErrors = false
			} else if (strings.HasPrefix(line, "    RX: ")) { // If line is RX header, next line will be RX data
				nextRx = true
			} else if (strings.HasPrefix(line, "    TX: ")) { // If line is TX header, next line will be TX data
				nextTx = true
			} else if (strings.HasPrefix(line, "    RX errors: ")) { // If line is RX errors header, next line will be RX errors data
				nextRxErrors = true
			} else if (strings.HasPrefix(line, "    TX errors: ")) { // If line is TX errors header, next line will be TX errors data
				nextTxErrors = true
			}
		}
	}
	return data
}