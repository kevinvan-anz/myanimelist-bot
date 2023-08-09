package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

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
			"ID",
			"start_season",
			"studios",
			"num_episodes",
			"average_episode_duration",
		})
	if err != nil {
		fmt.Printf("Unable to find anime: %v", err.Error())
		return
	}

	fmt.Printf("%s\n", anime.Title)
	fmt.Printf("ID: %d\n", anime.ID)
	fmt.Printf("Premier: %d %s\n", anime.StartSeason.Year, anime.StartSeason.Season)
	fmt.Printf("Studio: %v\n", anime.Studios)
	fmt.Printf("Episodes: %d\n", anime.NumEpisodes)
	fmt.Printf("Episode Duration: %d minutes\n", anime.AverageEpisodeDuration/60)
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
