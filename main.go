package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	gtfs "github.com/go-gtfs-server/com.google.transit.realtime"
	"google.golang.org/protobuf/proto"
)

func getTrip(feed *gtfs.FeedMessage, id string, currentStop string, direction uint32) {
	fmt.Println("SHOWING TRIPS FOR LINE:", id)
	fmt.Println("Looking for stop: ", currentStop)
	now := time.Now().Unix()

	for i := 0; i < len(feed.Entity); i++ {
		entity := feed.Entity[i].TripUpdate
		var routeDir uint32
		if entity.Trip.DirectionId != nil {
			routeDir = *entity.Trip.DirectionId
		} else {
			routeDir = 0
		}
		if *entity.Trip.RouteId == id && routeDir == direction { // loop through the realtime trips to find one on our route
			fmt.Println("Bus direction: ", routeDir)
			count := 0
			stopsAway := 0
			fmt.Println("Trip ID:", *entity.Trip.TripId)
			// fmt.Println(feed.Entity[i])
			lastStop := ""
			for j := 0; j < len(entity.StopTimeUpdate); j++ {
				tripUpdate := entity.StopTimeUpdate[j]
				fmt.Print(*tripUpdate.StopId, " ")
				if tripUpdate.Departure != nil {
					departureTime := *tripUpdate.Departure.Time
					fmt.Print("Estimated Departure Time:", departureTime)
					if now > departureTime {
						fmt.Print("*")
					} else {
						diff := departureTime - now
						minutes := diff / 60
						fmt.Print(" Departing in ", minutes, " minutes")
						lastStop = *tripUpdate.StopId
					}
				}

				if lastStop != "" && lastStop != currentStop {
					count += 1
				} else {
					stopsAway = count
				}

				fmt.Println()
			}

			fmt.Println("This bus is ", stopsAway, " stops away from your location")
			fmt.Println()
		}
	}
}

func main() {
	resp, err := http.Get("https://tmgtfsprd.sorttrpcloud.com/TMGTFSRealTimeWebService/tripupdate/tripupdates.pb")

	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	feed := &gtfs.FeedMessage{}
	if err := proto.Unmarshal(body, feed); err != nil {
		log.Fatalln("Failed to parse trip update")
	}

	getTrip(feed, "17", "CLI2659e", 1)

}
