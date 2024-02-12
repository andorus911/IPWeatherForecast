package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

type OpenMeteoForecast struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	CurrentUnits         struct {
		Time          string `json:"time"`
		Interval      string `json:"interval"`
		Temperature2M string `json:"temperature_2m"`
	} `json:"current_units"`
	Current struct {
		Time          cTime   `json:"time"`
		Interval      int     `json:"interval"`
		Temperature2M float64 `json:"temperature_2m"`
	} `json:"current"`
}

type Coordinates struct {
	lat float64
	lon float64
}

func main() {
	// primitive cli for now
	for {
		var choice int

		fmt.Printf("Choose an option:\n1. Get forecast\n2. Exit\n")
		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			// to add location determinator
			fmt.Printf("Enter coordinates (lat, lon): ")
			var c Coordinates
			fmt.Scanf("%f %f", &c.lat, &c.lon)
			getOpenMeteoForecast(c)
		case 2:
			return
		default:
			fmt.Println("Incorrect option!")
		}
	}
}

// get weather forecast via coords
func getOpenMeteoForecast(c Coordinates) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{Transport: tr}

	ts := time.Now().Format("2006-01-02")
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m&start_date=%s&end_date=%s", c.lat, c.lon, ts, ts)

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	var forecast OpenMeteoForecast
	err = json.Unmarshal(body, &forecast)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", forecast)
}

// it is for future
func getLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}
