package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"time"
)

func getLocalIPs() ([]net.IP, []net.IPMask, error) {
	var ips []net.IP
	var masks []net.IPMask
	ifaces, err := net.Interfaces()
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.To4() != nil && !v.IP.To4().IsLoopback() {
					fmt.Println(v.Mask)
				}
			}
		}
	}
	if err != nil {
		log.Panic().Msg(err.Error())
		return nil, nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Panic().Msg(err.Error())
			return nil, nil, err
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.To4() != nil && !v.IP.To4().IsLoopback() {
					ips = append(ips, v.IP)
					masks = append(masks, v.Mask)
				}
			}
		}
	}

	return ips, masks, err
}

func getPrinterIP(parsedConfig configuration) string {
	ips, _, err := getLocalIPs()

	if err != nil {
		log.Panic().Msg(err.Error())
	}

	fmt.Println(ips)

	for i, s := range parsedConfig.Printers.Buddy {
		if testConnection(s.Address) {
			version, _, _, _, _, _, _, _ := getBuddyResponse(s)
			parsedConfig.Printers.Buddy[i].Reachable = true
			parsedConfig.Printers.Buddy[i].Type = version.Hostname
		} else {
			parsedConfig.Printers.Buddy[i].Reachable = false
			log.Error().Msg(s.Address + " is not reachable")
		}
	}

	return ""
}

func testConnection(s string) bool {
	req, _ := http.NewRequest("GET", "http://"+s+"/", nil)
	client := &http.Client{Timeout: time.Duration(config.Exporter.ScrapeTimeout) * time.Second}
	r, e := client.Do(req)
	return e == nil && r.StatusCode == 200
}

func probeConfigFile(parsedConfig configuration) configuration {
	for i, s := range parsedConfig.Printers.Buddy {
		if testConnection(s.Address) {
			parsedConfig.Printers.Buddy[i].Reachable = true
		} else {
			parsedConfig.Printers.Buddy[i].Reachable = false
			log.Error().Msg(s.Address + " is not reachable")
		}
	}
	return parsedConfig
}
