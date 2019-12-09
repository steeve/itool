package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/simulatelocation"
)

var locationPlayCmdFlags = struct {
	duration time.Duration
	loop     bool
}{}

func init() {
	locationPlayCmd.Flags().DurationVarP(&locationPlayCmdFlags.duration, "duration", "d", 1*time.Second, "duration between points")
	locationPlayCmd.Flags().BoolVarP(&locationPlayCmdFlags.loop, "loop", "l", false, "loop endlessly")

	locationCmd.AddCommand(locationSetCmd)
	locationCmd.AddCommand(locationResetCmd)
	locationCmd.AddCommand(locationPlayCmd)

	rootCmd.AddCommand(locationCmd)
}

var locationCmd = &cobra.Command{
	Use:   "location",
	Short: "Simulate location",
}

var locationSetCmd = &cobra.Command{
	Use:   "set LATITUDE LONGITUDE",
	Short: "Set location",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := simulatelocation.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		latitude, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		longitude, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		if err := client.SetLocation(latitude, longitude); err != nil {
			log.Fatal(err)
		}
	},
}

var locationResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset location",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := simulatelocation.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		if err := client.ResetLocation(); err != nil {
			log.Fatal(err)
		}
	},
}

var locationPlayCmd = &cobra.Command{
	Use:   "play FILE.gpx ...",
	Short: "Play locations from one or more GPX files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// client := simulatelocation.NewClient(getUDID())
		doLoop := true
		for doLoop {
			doLoop = locationPlayCmdFlags.loop
			for _, arg := range args {
				log.Println("Opening", arg)
				data, err := ioutil.ReadFile(arg)
				if err != nil {
					log.Fatal(err)
				}
				gpxData := &struct {
					WayPoints []struct {
						Latitude  float64 `xml:"lat,attr"`
						Longitude float64 `xml:"lon,attr"`
					} `xml:"wpt"`
					Tracks []struct {
						Points []struct {
							Latitude  float64 `xml:"lat,attr"`
							Longitude float64 `xml:"lon,attr"`
						} `xml:"trkpt"`
					} `xml:"trk"`
				}{}
				if err := xml.Unmarshal(data, gpxData); err != nil {
					log.Fatal(err)
				}
				for i, waypoint := range gpxData.WayPoints {
					log.Printf("waypoint:%d/%d latitude:%f longitude:%f\n", i+1, len(gpxData.WayPoints), waypoint.Latitude, waypoint.Longitude)
					time.Sleep(locationPlayCmdFlags.duration)
				}
				for i, track := range gpxData.Tracks {
					for j, point := range track.Points {
						log.Printf("track:%d %d/%d: latitude:%f longitude:%f\n", i+1, j+1, len(track.Points), point.Latitude, point.Longitude)
						time.Sleep(locationPlayCmdFlags.duration)
					}
				}
			}
		}
	},
}
