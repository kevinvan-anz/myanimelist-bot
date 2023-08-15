package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/nstratos/go-myanimelist/mal"
)

func main() {
	clientIDUsername := os.Getenv("MALCLIENTID")
	publicInfoClient := &http.Client{
		// Create client ID from https://myanimelist.net/apiconfig.
		Transport: &clientIDTransport{ClientID: clientIDUsername},
	}

	c := mal.NewClient(publicInfoClient)

	ctx := context.Background()

	anime, _, err := c.Anime.Details(ctx, 51009,
		mal.Fields{
			"mean",
			"rank",
			"start_season",
			"broadcast",
			"popularity",
			"studios",
			"num_episodes",
			"average_episode_duration",
		})
	if err != nil {
		fmt.Printf("Unable to find anime: %v", err.Error())
		return
	}

	// Converting JST to AEST
	jstTimeStr := "2023-05-10 15:30:00"
	jstTime, aestTime, err := convertJSTToAEST(jstTimeStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Sample output for anime details
	fmt.Printf("%s\n", anime.Title)
	fmt.Printf("Score: %v\n", anime.Mean)
	fmt.Printf("Ranking: #%d\n", anime.Rank)
	fmt.Printf("Popularity: #%d\n", anime.Popularity)
	fmt.Printf("Premier: %d %s\n", anime.StartSeason.Year, anime.StartSeason.Season)
	// TO DO: Convert JST value to AEST
	fmt.Printf("Broadcast: %v at %v JST+1\n", anime.Broadcast.DayOfTheWeek, anime.Broadcast.StartTime)
	for _, s := range anime.Studios {
		fmt.Printf("Studio: %s\n", s.Name)
	}
	fmt.Printf("Episodes: %d\n", anime.NumEpisodes)
	fmt.Printf("Episode Duration: %d minutes\n", anime.AverageEpisodeDuration/60)

	fmt.Println("JST Time:", jstTime)
	fmt.Println("AEST Time:", aestTime)
}

type clientIDTransport struct {
	Transport http.RoundTripper
	ClientID  string
}

func (c *clientIDTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}
	req.Header.Add("X-MAL-CLIENT-ID", c.ClientID)
	return c.Transport.RoundTrip(req)
}

func convertJSTToAEST(jstTimeStr string) (time.Time, time.Time, error) {
	layout := "2006-01-02 15:04:05"

	// Parse JST time
	jstLocation, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	jstTime, err := time.ParseInLocation(layout, jstTimeStr, jstLocation)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Convert JST to AEST
	aestLocation, err := time.LoadLocation("Australia/Melbourne")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	aestTime := jstTime.In(aestLocation)

	return jstTime, aestTime, nil
}
