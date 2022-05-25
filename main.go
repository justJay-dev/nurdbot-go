package main

import (
	"fmt"
	"os"

	"github.com/gempir/go-twitch-irc/v3"

	"github.com/fatih/color"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {

	// or client := twitch.NewAnonymousClient() for an anonymous user (no write capabilities)
	client := twitch.NewClient(os.Getenv("TWITCH_USERNAME"), os.Getenv("TWITCH_OAUTH"))

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(color.RedString(message.User.Name)+" =>", message.Message)
	})

	client.Join("pnJay")

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
