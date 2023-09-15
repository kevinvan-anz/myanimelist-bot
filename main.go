package main

import (
	"context"
	"flag"
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

	animeID, err := readAnimeID()
	if err != nil {
		fmt.Errorf("invalid animeID input from MAL, error: %v", err)
		return
	}

	// Retrieving anime data fields
	anime, _, err := c.Anime.Details(ctx, animeID,
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
		fmt.Printf("Unable to find anime using input, error: %v", err.Error())
		return
	}

	animeData := processAnimeData(*anime)
	printAnimeData(animeData)

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
	jstStartTime := &anime.Broadcast.StartTime
	_, aestTime, err := convertJSTToAEST(*jstStartTime)
	if err != nil {
		fmt.Println("Error:", err)
		return AnimeData{}
	}

	animeData := AnimeData{
		Title:                 anime.Title,
		Synopsis:              anime.Synopsis,
		Mean:                  anime.Mean,
		Rank:                  anime.Rank,
		Popularity:            anime.Popularity,
		StartSeason:           fmt.Sprintf("%d %s", anime.StartSeason.Year, anime.StartSeason.Season),
		Broadcast:             fmt.Sprintf("%v at %v AEST", anime.Broadcast.DayOfTheWeek, aestTime),
		Studios:               make([]string, len(anime.Studios)),
		NumEpisodes:           anime.NumEpisodes,
		AverageEpisodeMinutes: anime.AverageEpisodeDuration / 60,
	}

	for i, s := range anime.Studios {
		animeData.Studios[i] = s.Name
	}

	return animeData
}

func printAnimeData(animeData AnimeData) {
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
	fmt.Printf("Episodes: %d\n", animeData.NumEpisodes)
	fmt.Printf("Episode Duration: %d minutes\n\n", animeData.AverageEpisodeMinutes)

	fmt.Printf("OBJECT PRINT: %+v\n", animeData)
}

func readAnimeID() (int, error) {
	// Define flag for Anime selection
	var animeIDFlag = flag.Int("animeID", 0, "Anime ID")

	// Parse the flag to the command line into the defined flags
	flag.Parse()

	// Check if the "animeID" flag was provided
	if *animeIDFlag == 0 {
		fmt.Println("Usage: program --animeID <animeID>")
		os.Exit(1)
	}

	// Validate the anime ID
	animeID := *animeIDFlag
	if animeID <= 0 {
		return 0, fmt.Errorf("invalid anime ID: %d", animeID)
	}

	return animeID, nil
}

func convertJSTToAEST(jstTimeStr string) (time.Time, time.Time, error) {
	layout := "2006-01-02 15:04"
	jstTimeStr = "2023-01-01 " + jstTimeStr

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
