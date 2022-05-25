package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"io/ioutil"

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
	client := twitch.NewClient(os.Getenv("TWITCH_USERNAME"), os.Getenv("TWITCH_OAUTH"))

	nurdbotSay := func(message twitch.PrivateMessage, output string) {
		fmt.Println(color.BlueString("nurdbot:"), output)
		client.Say(message.Channel, output)
	}

	readFromTextFile := func(fileName string) string {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		responses, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		lines := strings.Split(string(responses), "\n")
		return lines[rand.Intn(len(lines))]
	}

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(color.YellowString(message.User.Name+":"), message.Message)

		if message.Message == "!joke" {
			nurdbotSay(message,
				readFromTextFile("./replies/jokes.txt"))
		}

		if message.Message == "!fart" {
			nurdbotSay(message,
				readFromTextFile("./replies/farts.txt"))
		}

		if message.Message == "!panic" {
			nurdbotSay(message,
				readFromTextFile("./replies/panic.txt"))
		}

		if message.Message == "!merch" {
			nurdbotSay(message, "https://www.nurdbot.com/uwu")
		}

		if message.Message == "!uptime" {
			resp, err := http.Get("https://beta.decapi.me/twitch/uptime/" + message.Channel)

			if err != nil {
				fmt.Println(err)
				nurdbotSay(message, "Error getting uptime :(")
			}

			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				fmt.Println(err)
				nurdbotSay(message, "Error getting uptime :(")
			}
			nurdbotSay(message, string(body))
		}

		if message.User.Name == "missqueeney" {
			// todo you are here. spongebob case
			s := "averylargeword"
			v := strings.SplitAfter(s, "")
			fmt.Println(v) // [a v e r y l a r g e w o r d]
			nurdbotSay(message, "hi queen :)")
		}

	})

	client.Join("pnJay")
	fmt.Println(color.BlueString("=== Beep boop. Nurdbot connected ==="))

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
