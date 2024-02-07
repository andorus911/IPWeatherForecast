package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
  "encoding/json"
)

type OpenMeteoForecast struct {
  Latitude float32 `json:"latitude"`
  Longitude float32 `json:"longitude"`
  GenerationTime_ms float32 `json:"generationtime_ms"`
  UTC_Offset_Seconds int `json:"utc_offset_seconds"`
  TimeZone string `json:"timezone"`
  TimeZoneAbbr string `json:"timezone_abbreviation"`
  Elevation float32 `json:"elevation"`
  CurrentUnits struct {
    Time string `json:"time"`
    Interval string `json:"interval"`
    Temp_2m string `json:"temperature_2m"`
  } `json:"current_units"`
  Current struct {
    Time string `json:"time"`// date?
    Interval int `json:"interval"`
    Temp_2m float32 `json:"temperature_2m"`
  } `json:"current"`
}

func main() {
	// get weather forecast via coords
  tr := &http.Transport{
    MaxIdleConns:       10,
    IdleConnTimeout:    30 * time.Second,
    DisableCompression: true,
  }
  client := &http.Client{Transport: tr}
  resp, err := client.Get("https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&current=temperature_2m&start_date=2024-02-07&end_date=2024-02-07")
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