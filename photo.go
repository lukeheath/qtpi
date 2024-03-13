package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// pexelsAuthToken is the authorization token for the Pexels API, loaded from environment variables.
var pexelsAuthToken = os.Getenv("PEXELS_AUTH_TOKEN")

// Photo defines the structure for a photo obtained from the Pexels API, including various fields such as the URL,
// photographer name, photographer's URL, different sources of the photo, and a caption.
type Photo struct {
	Url             string `json:"url"`
	Photographer    string `json:"photographer"`
	PhotographerUrl string `json:"photographer_url"`
	Src             struct {
		Original string `json:"original"`
		Portrait string `json:"portrait"`
	} `json:"src"`
	Caption string `json:"alt"`
}

// PexelsResponse defines the structure of the response received from the Pexels API. It includes an array of Photo objects.
type PexelsResponse struct {
	Photos []Photo `json:"photos"`
}

// getPhoto fetches a random photo from the Pexels API based on the "cute animals" query.
// It returns a Photo object and an error. In case of an error or if no photos are returned by the API,
// a default Photo is returned.
func getPhoto(message string) (Photo, error) {
	// defaultPhoto is returned in case of any error during the API request or when no photos are found.
	defaultPhoto := Photo{
		Url:             "https://images.pexels.com/photos/50577/hedgehog-animal-baby-cute-50577.jpeg",
		Photographer:    "Pixabay",
		PhotographerUrl: "https://www.pexels.com/@pixabay",
		Caption:         "Hedgehog",
	}

	// Get search term from the message and use it to query the Pexels API.
	searchTerm := "cute animals"
	if message != "" {
		searchTerm, _ = getSearchTerm(message)
	}

	// Replace all spaces for the API request.
	searchTerm = strings.ReplaceAll(searchTerm, " ", "+")

	// Construct the request to the Pexels API with the required headers.
	req, err := http.NewRequest("GET", "https://api.pexels.com/v1/search?per_page=1&query="+searchTerm, nil)
	if err != nil {
		return defaultPhoto, err
	}
	req.Header.Set("Authorization", pexelsAuthToken)

	// Perform the HTTP request.
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return defaultPhoto, err
	}
	// Ensure the response body is closed after the function returns.
	defer response.Body.Close()

	// Check for a successful response status code.
	if response.StatusCode != http.StatusOK {
		return defaultPhoto, fmt.Errorf("API request failed with status code %d", response.StatusCode)
	}

	// Read and unmarshal the response body.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return defaultPhoto, err
	}

	var pexelsRes PexelsResponse
	err = json.Unmarshal(body, &pexelsRes)
	if err != nil {
		return defaultPhoto, err
	}

	// Return the first photo if available.
	if len(pexelsRes.Photos) > 0 {
		photo := pexelsRes.Photos[0]
		// Update the photo's URL to the portrait version for consistency.
		photo.Url = photo.Src.Portrait
		return photo, nil
	}

	// Return the default photo if no photos are found.
	return defaultPhoto, nil
}
