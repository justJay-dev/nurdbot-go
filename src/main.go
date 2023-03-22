package main

import (
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

	"gorm.io/gorm"

	"github.com/glebarez/sqlite"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	client := twitch.NewClient(os.Getenv("TWITCH_USERNAME"), os.Getenv("TWITCH_OAUTH"))

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
		gorm.Model
		ID      uint      `gorm:"primaryKey" json:"id"`
		Message string    `json:"message"`
		User    string    `json:"user"`
		Channel string    `json:"channel"`
		Time    time.Time `json:"time"`
	}
	//do database stuff.
	db, dbErr := gorm.Open(sqlite.Open("./database/nurdbot.db"), &gorm.Config{})
	if dbErr != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&LogMessage{})

	writeLog := func(message twitch.PrivateMessage) {
		log := LogMessage{
			Message: message.Message,
			User:    message.User.Name,
			Channel: message.Channel,
			Time:    message.Time,
		}
		db.Create(&log)
	}
	nurdbotSay := func(message twitch.PrivateMessage, output string) {
		//terminal output
		fmt.Println(color.BlueString("nurdbot:"), output)
		//write the reply to the log, this is technically incorrect on timing but close enough.
		db.Create(&LogMessage{Message: output, User: "nurdbot", Channel: message.Channel, Time: message.Time})
		//actually write it in chat.
		client.Say(message.Channel, output)
	}

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(color.YellowString(message.User.Name+":"), message.Message)
		writeLog(message)

		if message.Message == "!joke" {
			nurdbotSay(message,
				readFromTextFile("./replies/jokes.txt"))
		}

		if message.Message == "!github" {
			nurdbotSay(message, "https://github.com/justJay-dev/nurdbot-go")
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
			resp, respErr := http.Get("https://beta.decapi.me/twitch/uptime/" + message.Channel)

			if respErr != nil {
				fmt.Println(respErr)
				nurdbotSay(message, "Error getting uptime :(")
			}

			body, bodyErr := ioutil.ReadAll(resp.Body)

			if bodyErr != nil {
				fmt.Println(bodyErr)
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

	twitchClientErr := client.Connect()
	if twitchClientErr != nil {
		panic(twitchClientErr)
	}
}
