package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NasaResponse struct {
	Explanation string `json:"explanation"`
	Hdurl       string `json:"hdurl"`
	MediaType   string `json:"media_type"`
	Title       string `json:"title"`
	Url         string `json:"url"`
}

func main() {
	today := time.Now().Format("2006-01-02")
	nasaResponse := question_nasa(today)

	send_image(nasaResponse)
}

func send_image(NasaResponse NasaResponse) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	imageURL := tgbotapi.FileURL(NasaResponse.Hdurl)
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(imageURL))
	photo.Caption = NasaResponse.Title

	_, err = bot.Send(photo)
	if err != nil {
		log.Panic(err)
	}
}

func question_nasa(date string) NasaResponse {
	client := http.Client{}
	api_key := os.Getenv("NASA_API_KEY")
	url := fmt.Sprintf("https://api.nasa.gov/planetary/apod?date=%s&api_key=%s", date, api_key)

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var parsedNasaResponse NasaResponse
	err = json.Unmarshal(body, &parsedNasaResponse)

	if err != nil {
		fmt.Printf("There was an error decoding the json. err = %s", err)
	}

	return parsedNasaResponse
}
