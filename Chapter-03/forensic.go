package main

func getMacCodes() []string {
	const (
		netkey = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\NetworkList\Signatures\Unmanaged`
	)
	// k, err := registry.OpenKey(registry.LOCAL_MACHINE, netkey, registry.QUERY_VALUE)
	// if err != nil {
	// 	fmt.Println("Could not read from windows registry!")
	// 	os.Exit(1)
	// }
	// defer k.Close()
	// Need to write this on a windows machine. Not able to figure out the values
	return []string{""}
}

func getGeoLoc(maccodes []string) []GeoCode {
	var geocodes []GeoCode
	in := make(chan string, len(maccodes))
	out := make(chan GeoCode, len(maccodes))

	for _, mac := range maccodes {
		in <- mac
		go func(in <-chan string, out chan<- GeoCode) {
			mac, more := <-in
			geocode, err := WiggleGet(mac)
			if err == nil {
				out <- geocode
			}
			if !more {
				close(out)
			}
		}(in, out)
	}
	close(in)

	for geocode := range out {
		geocodes = append(geocodes, geocode)
	}
	return geocodes
}

// func main() {
// 	for _, loc := range getGeoLoc(getMacCodes()) {
// 		fmt.Printf("You are connected to network : %s\n", loc.ssid)
// 		fmt.Printf("Latitude : %f Longitude : %f\n", loc.lat, loc.lng)
// 	}
// }
