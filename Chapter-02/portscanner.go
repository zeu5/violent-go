package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strings"
)

func parseArgs() (string, []string) {

	var host, ports string
	flag.StringVar(&host, "h", "", "Host to scan for")
	flag.StringVar(&ports, "p", "", "Ports (comma seperated) to scan for")

	if host == "" || ports == "" {
		panic("Usage scanner -h <ip> -p <ports>")
	}

	return host, strings.Split(ports, ",")
}

func scan(ip, port string) (string, error) {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return "", err
	}
	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	return data, nil
}

func main() {
	host, ports := parseArgs()
	fmt.Println("Scanning for host: " + host)
	for _, port := range ports {
		go func(host, port string) {
			if banner, err := scan(host, port); err == nil {
				fmt.Println("Could connect on port : " + port + ", with banner :" + banner)
			}
		}(host, port)
	}

	fmt.Println("Done!")
}
