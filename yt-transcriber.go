package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/kkdai/youtube/v2"
)

func main() {
	r := gin.Default()
	r.GET("/transcript", getTranscript)

	// PORT environment variable is provided by Cloud Run.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Print("Hello from Cloud Run! The container started successfully and is listening for HTTP requests on $PORT")
	log.Printf("Listening on port %s", port)
	r.Run(fmt.Sprintf(":%s", port))
}

func getTranscript(c *gin.Context) {
	youtubeURL := c.Query("url")
	if youtubeURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing YouTube URL"})
		return
	}

	videoID, err := extractVideoID(youtubeURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid YouTube URL"})
		return
	}

	client := youtube.Client{}
	video, err := client.GetVideo(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video info"})
		return
	}

	captions, err := client.GetTranscript(video, "en")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch captions"})
		return
	}

	transcript := ""
	for _, caption := range captions {
		transcript += caption.Text + " "
	}

	c.JSON(http.StatusOK, gin.H{"transcript": transcript})
}

func extractVideoID(url string) (string, error) {
	regex := regexp.MustCompile(`(?:v=|\/)([0-9A-Za-z_-]{11}).*`)
	matches := regex.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("invalid YouTube URL")
}
