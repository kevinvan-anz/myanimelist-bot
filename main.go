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

	anime, response, option := c.Anime.Details(ctx, 51009)

	fmt.Println(anime, response, option)
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
