package main

import (
	geoip "SSHMonitor/internal/GeoIP"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	"github.com/joho/godotenv"
)

var WEBHOOK string

type payload struct {
	Content     *string   `json:"content"`
	Embeds      []*embed  `json:"embeds"`
	Attachments []*string `json:"attachments"`
}

type author struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type embed struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Color       int     `json:"color"`
	Author      *author `json:"author"`
	Timestamp   string  `json:"timestamp"`
}

func main() {

	logFile, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Set the log output to the log file
	log.SetOutput(logFile)

	err = godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	WEBHOOK = os.Getenv("WEBHOOK_URL")
	if WEBHOOK == "" {
		log.Fatalln("Empty webhook")
	}

	// Set up the tail to read the log file
	t, err := tail.TailFile("/var/log/auth.log", tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: 0, Whence: 2}})
	if err != nil {
		log.Fatal(err)
	}
	defer t.Cleanup()

	for line := range t.Lines {
		if strings.Contains(line.Text, "Accepted") || strings.Contains(line.Text, "Failed") {
			// You can use regular expressions to extract relevant information from the log line
			// For example, to extract IP addresses and usernames
			pattern := regexp.MustCompile(`(\S+) from (\S+)`)

			matches := pattern.FindStringSubmatch(line.Text)

			if len(matches) > 2 {
				username := matches[1]
				ipAddress := matches[2]

				if strings.Contains(line.Text, "Accepted") {
					// Handle successful login
				} else {
					postWithAddr(ipAddress, username)
				}
			}
		}

	}

}

func postWithAddr(addr, user string) {
	gip := &geoip.GeoIP{IP: addr}

	switch gip.Get() {
	case http.StatusBadRequest:
	case http.StatusNotFound:
	case http.StatusInternalServerError:
		log.Printf("Failed to post for %v\n", gip)
		return

	case http.StatusTooManyRequests:
		log.Println("IP Limit reached")
		return
	}

	e := failEmbed(user, gip)
	p := payload{nil, nil, nil}
	p.Embeds = append(p.Embeds, e)

	body, _ := json.Marshal(p)

	resp, err := http.Post(WEBHOOK, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Panic(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		log.Println("Post Failed with code", resp.StatusCode)
	}

}

func failEmbed(user string, gip *geoip.GeoIP) *embed {
	timestamp := time.Now().UTC() // Get the current UTC time

	// Format the timestamp as desired
	formattedTimestamp := timestamp.Format("2006-01-02T15:04:05.000Z")

	author := &author{
		Name: "ExTBH",
		Url:  "https://extbh.dev",
	}
	desc := fmt.Sprintf("**User:** %s\n**from:** %s :flag_%s:\n**Country:** %s\n**City:** %s", user, gip.IP, strings.ToLower(gip.CountryCode), gip.CountryName, gip.CityName)

	e := &embed{
		Title:       "New Failed Attempt",
		Description: desc,
		Color:       16734296,
		Author:      author,
		Timestamp:   formattedTimestamp,
	}

	return e

}
