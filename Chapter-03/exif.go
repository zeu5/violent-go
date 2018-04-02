// Given a URL. The resulting page is parsed for img tags
// and exif data is extracted for all of them

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

type exifTags struct {
	Artist, Model  string
	GeoLat, GeoLng float64
	Time           time.Time
}

func parseExif(f io.Reader) (exifTags, error) {
	tags := exifTags{}

	exif.RegisterParsers(mknote.All...)
	exifData, err := exif.Decode(f)
	if err != nil {
		return tags, err
	}

	artist, err := exifData.Get(exif.Artist)
	if err == nil {
		tags.Artist, _ = artist.StringVal()
	}
	model, err := exifData.Get(exif.Model)
	if err == nil {
		tags.Model, _ = model.StringVal()
	}
	lat, long, err := exifData.LatLong()
	if err == nil {
		tags.GeoLat = lat
		tags.GeoLng = long
	}
	time, err := exifData.DateTime()
	if err == nil {
		tags.Time = time
	}

	return tags, nil
}

func getImageURLs(url string) ([]string, error) {
	var imageurls []string

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return imageurls, err
	}
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		imageurl, exists := s.Attr("src")
		if exists {
			imageurls = append(imageurls, imageurl)
		}
	})

	return imageurls, nil
}

func getExifURL(url string) (exifTags, error) {
	exiftags := exifTags{}

	response, err := http.Get(url)
	if err != nil {
		return exiftags, err
	}
	defer response.Body.Close()
	exiftags, err = parseExif(response.Body)
	if err != nil {
		return exiftags, err
	}
	return exiftags, nil
}

func parseArgs() string {
	if len(os.Args) < 2 {
		fmt.Println("Usage exif <url>")
		os.Exit(1)
	}
	return os.Args[1]
}

func main() {
	url := parseArgs()
	fmt.Println("Fetched URL : ", url)
	imageurls, err := getImageURLs(url)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Fetched image tags : ")
	fmt.Println(imageurls)

	if len(imageurls) == 0 {
		fmt.Println("Could not read any image tags!")
		return
	}

	in := make(chan string, len(imageurls))
	out := make(chan exifTags, len(imageurls))

	for _, imageurl := range imageurls {
		in <- imageurl
		go func(in <-chan string, out chan<- exifTags) {
			imageurl, more := <-in
			fmt.Println("Fetching Image : ", imageurl)
			exifdata, err := getExifURL(imageurl)
			if err == nil {
				out <- exifdata
			}
			if !more {
				close(out)
			}
		}(in, out)
	}
	close(in)

	for exifdata := range out {
		fmt.Printf("%+v", exifdata)
	}
}
