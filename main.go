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
	// Requesting MAL API
	clientIDUsername := os.Getenv("MALCLIENTID")
	publicInfoClient := &http.Client{
		// Create client ID from https://myanimelist.net/apiconfig.
		Transport: &clientIDTransport{ClientID: clientIDUsername},
	}

	c := mal.NewClient(publicInfoClient)

	ctx := context.Background()

	// Retrieving anime data fields
	anime, _, err := c.Anime.Details(ctx, 51179,
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

	animeData := processAnimeData(*anime)

	// Sample output for anime details
	fmt.Printf("%s\n", animeData.Title)
	fmt.Printf("Score: %v\n", animeData.Mean)
	fmt.Printf("Ranking: #%d\n", animeData.Rank)
	fmt.Printf("Popularity: #%d\n", animeData.Popularity)
	fmt.Printf("Premier: %s\n", animeData.StartSeason)
	// TO DO: Convert JST value to AEST
	fmt.Printf("Broadcast: %s\n", animeData.Broadcast)
	for _, studio := range animeData.Studios {
		fmt.Printf("Studio: %s\n", studio)
	}
	fmt.Printf("Episodes: %d\n", anime.NumEpisodes)
	fmt.Printf("Episode Duration: %d minutes\n", animeData.AverageEpisodeMinutes)

	fmt.Println("JST Time:", jstTime)
	fmt.Println("AEST Time:", aestTime)
}

type clientIDTransport struct {
	Transport http.RoundTripper
	ClientID  string
}

// RoundTrip - Authentication for API client ID
func (c *clientIDTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}
	req.Header.Add("X-MAL-CLIENT-ID", c.ClientID)
	return c.Transport.RoundTrip(req)
}

// AnimeData list
type AnimeData struct {
	Title                 string
	Synopsis              string
	Mean                  float64
	Rank                  int
	Popularity            int
	StartSeason           string
	Broadcast             string
	Studios               []string
	NumEpisodes           int
	AverageEpisodeMinutes int
}

// processAnimeData - Processing data from AnimeData struct
func processAnimeData(anime mal.Anime) AnimeData {
	animeData := AnimeData{
		Title:       anime.Title,
		Synopsis:    anime.Synopsis,
		Mean:        anime.Mean,
		Rank:        anime.Rank,
		Popularity:  anime.Popularity,
		StartSeason: fmt.Sprintf("Premier: %d %s", anime.StartSeason.Year, anime.StartSeason.Season),
		Broadcast:   fmt.Sprintf("Broadcast: %v at %v JST+1", anime.Broadcast.DayOfTheWeek, anime.Broadcast.StartTime),
		// Need explanation
		Studios:               make([]string, len(anime.Studios)),
		NumEpisodes:           anime.NumEpisodes,
		AverageEpisodeMinutes: anime.AverageEpisodeDuration / 60,
	}

	for i, s := range anime.Studios {
		animeData.Studios[i] = s.Name
	}

	return animeData
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
