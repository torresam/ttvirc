package main

import (
	"crypto/tls"
	"flag"
	"os"
	"strings"
	"time"

	set "github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/torresam/ttvirc/controllers"
	"github.com/torresam/ttvirc/twitch"
	"gopkg.in/irc.v3"
)

func main() {
	debugEnabled := flag.Bool("debug", false, "Whether debug mode should be enabled")
	flag.Parse()

	if *debugEnabled {
		log.SetLevel(log.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetOutput(colorable.NewColorableStdout())

	conf := &tls.Config{}
	conn, err := tls.Dial("tcp", "irc.chat.twitch.tv:6697", conf)
	if err != nil {
		log.Fatalln(err)
	}
	channels := set.NewSet()
	config := irc.ClientConfig{
		Nick:    os.Getenv("TWITCH_USER"),
		Pass:    os.Getenv("TWITCH_OAUTH"),
		User:    os.Getenv("TWITCH_USER"),
		Name:    os.Getenv("TWITCH_USER"),
		Handler: irc.HandlerFunc(twitch.MessageHandler(&channels)),
	}

	// Create the client
	log.Debug("Creating client")
	client := irc.NewClient(conn, config)
	client.CapRequest("twitch.tv/commands", true)
	client.CapRequest("twitch.tv/tags", true)
	client.CapRequest("twitch.tv/membership", true)

	// Blocking channel to pass messages from HTTP endpoint to IRC client
	httpMessages := make(chan string)

	// Listen for HTTP messages on a separate Goroutine to pick up user input
	go func(messages chan string, channels *set.Set) {
		r := gin.Default()
		r.GET("/channels", controllers.ListChannels(channels))
		r.PUT("/join/:channel", controllers.JoinChannel(messages))
		r.PUT("/leave/:channel", controllers.LeaveChannel(messages))
		r.Run(":8000")
	}(httpMessages, &channels)

	// Goroutine channel that listens for messages and writes them directly into IRC
	// Only send messages with known commands
	go func(messages chan string, client *irc.Client) {
		log.Debug("Starting HTTP Message Backend Goroutine")
		count := 0
		lastJoin := time.Now()
		delayWindow := 30 * time.Second
		for {
			msg := <-messages
			log.Debug("Backend IRC routine received message: ", msg)
			if strings.HasPrefix(msg, "JOIN") {
				count++
				log.Debug("Joining channel: ", msg)
				client.Write(msg)
			} else if strings.HasPrefix(msg, "PART") {
				count++
				log.Debug("Leaving channel: ", msg)
				client.Write(msg)
			} else {
				log.Warn("Unrecognized IRC command received")
			}
			// Reset message count if it's been 30s since last message
			// TODO: move this to something that works better
			// The oldest message could be >30s but clearing count is based on latest only
			log.Debug("Time since last IRC command: ", time.Since(lastJoin))
			if time.Since(lastJoin) >= delayWindow {
				count = 0
			}
			lastJoin = time.Now()

			// https://dev.twitch.tv/docs/irc/guide#command--message-limits
			// Comply with the rate limit for JOIN/LEAVEs
			log.Debug("Message count in time period: ", count)
			if count >= 20 {
				log.Warn("Nearing rate limit for number of commands, sleeping for 30s")
				time.Sleep(delayWindow)
				count = 0
			}
		}
	}(httpMessages, client)

	log.Debug("Running client")
	err = client.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
