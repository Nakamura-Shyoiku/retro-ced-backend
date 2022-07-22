package trending

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/spf13/viper"
)

// MediumPost format
type MediumPost struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	ImageURL    string `json:"imageUrl"`
}

// MediumPosts posts
type MediumPosts struct {
	Posts []MediumPost `json:"data"`
}

// GetPosts service
type GetPosts struct {
	posts *MediumPosts
}

// NewGetPosts instance
func NewGetPosts() *GetPosts {
	n := new(GetPosts)
	n.posts = new(MediumPosts)
	return n
}

// Data returns medium posts
func (g *GetPosts) Data() *MediumPosts {
	return g.posts
}

// Do tasks
func (g *GetPosts) Do() (err error) {
	if err = g.getPosts(); err != nil {
		return err
	}

	return nil
}

func (g *GetPosts) getPosts() error {
	// Build the request
	url := fmt.Sprintf(
		"https://api.medium.com/v1/users/%s/publications",
		viper.GetString("medium.userID"),
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WithError(err).Error("failed to create requst")
		return err
	}
	// Add authentication
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", viper.GetString("medium.key")))

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("failed to do get medium posts")
		return err
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	// var records []MediumPost

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&g.posts); err != nil {
		log.WithError(err).Error("failed to unmarshal medium posts")
		return err
	}

	return nil
}
