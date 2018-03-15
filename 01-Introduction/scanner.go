package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const ipPrefix string = "172.16."

var ports = [...]int{21, 22, 25, 80, 110, 443}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func scan(ips []string, banners []string) {
	fmt.Println("Scanning")
	for _, ip := range ips {
		for _, port := range ports {
			banner, err := getBanner(ip, port)
			if err == nil {
				if ok := checkBanners(banner, banners); ok {
					fmt.Println("Found server : " + ip + " with banner : " + banner)
				}
			}
		}
	}
	fmt.Println("Done!")
}

func readLines(f *os.File) []string {
	// Read lines of the given file
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func parseArgs() (map[string]string, error) {
	if len(os.Args) == 1 {
		return nil, errors.New("Usage: scanner banners_file")
	}
	var args = make(map[string]string)
	args["file"] = os.Args[1]
	return args, nil
}

func generateIPs() []string {
	var ips []string
	if len(strings.Split(ipPrefix, ".")) == 3 {
		for i := 50; i < 70; i++ {
			for j := 120; j < 170; j++ {
				ips = append(ips, ipPrefix+strconv.Itoa(i)+"."+strconv.Itoa(j))
			}
		}
	} else if len(strings.Split(ipPrefix, ".")) == 4 {
		for j := 120; j < 170; j++ {
			ips = append(ips, ipPrefix+"."+strconv.Itoa(j))
		}
	}
	return ips
}

func getBanner(ip string, port int) (string, error) {
	conn, err := net.Dial("tcp", "ip:"+strconv.Itoa(port))
	if err != nil {
		return "", err
	}
	banner, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	return banner, nil
}

func checkBanners(banner string, bannerList []string) bool {
	for _, bannerSub := range bannerList {
		if strings.Contains(banner, bannerSub) {
			return true
		}
	}
	return false
}

func main() {

	args, err := parseArgs()
	check(err)

	f, err := os.Open(args["file"])
	check(err)
	defer f.Close()
	fmt.Printf("Reading from file: %s\n", args["file"])

	banners := readLines(f)
	fmt.Println("Fetched valid banners")

	ips := generateIPs()
	fmt.Println("Generated IPs to scan through")

	scan(ips, banners)
}
