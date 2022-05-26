package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

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

	//write out any message into a json file named %CURRENT_DATE%.json
	type LogMessage struct {
		Message string `json:"message"`
		User    string `json:"user"`
		Channel string `json:"channel"`
		Time    string `json:"time"`
	}

	writeLog := func(message twitch.PrivateMessage) {
		log := LogMessage{
			Message: message.Message,
			User:    message.User.Name,
			Channel: message.Channel,
			Time:    message.Time.String(),
		}

		fileName := time.Now().Format("2006-01-02") + ".json"
		f, err := os.OpenFile("logs/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		data, _ := json.MarshalIndent(log, "", " ")
		f.WriteString(string(data) + ",\n")
	}

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(color.YellowString(message.User.Name+":"), message.Message)
		writeLog(message)

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
			nurdbotSay(message, "https://todo.todo.org")
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
			// take the message.Message and turn it into spongbob case
			v := strings.SplitAfter(message.Message, "")
			for i := 0; i < len(v); i++ {
				if i%2 == 0 {
					v[i] = strings.ToUpper(v[i])

				} else {
					v[i] = strings.ToLower(v[i])
				}
			}
			nurdbotSay(message, strings.Join(v, ""))
		}

	})

	client.Join("pnJay")
	fmt.Println(color.BlueString("=== Beep boop. Nurdbot connected ==="))

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
