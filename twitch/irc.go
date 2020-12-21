package twitch

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	set "github.com/deckarep/golang-set"
	"gopkg.in/irc.v3"
)

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// MessageHandler -
//                Handles different commands based on Twitch IRC
func MessageHandler(channels *set.Set) irc.HandlerFunc {
	return func(c *irc.Client, m *irc.Message) {
		if isNumeric(m.Command) {
			log.Debug(m)
		} else if m.Command == "CAP" {
			log.Debug(m)
		} else if m.Command == "PING" {
			// Twitch IRC wants you to send this back when they ping you
			// https://dev.twitch.tv/docs/irc/guide#connecting-to-twitch-irc
			resp, _ := irc.ParseMessage("PONG :tmi.twitch.tv")
			log.Debug(resp)
			c.WriteMessage(resp)
		} else if m.Command == "PRIVMSG" {
			log.Debugf("Channel: %s, User: %s, --- %v", m.Params[0], m.User, strings.Join(m.Params[1:], " "))
		} else if m.Command == "CLEARCHAT" {
			// This is a ban
			log.Debugf("Channel: %s, User: %s was banned", m.Params[0], m.Params[1])
		} else if m.Command == "USERNOTICE" {
		} else if m.Command == "USERSTATE" {
		} else if m.Command == "CLEARMSG" {
		} else if m.Command == "JOIN" {
			// Server will ACK back your JOIN messages on success
			// Use this to confirm we joined the channel successfully
			if m.User == c.CurrentNick() {
				log.Debug("Joined the following channel: ", m.Params[0])
				(*channels).Add(m.Params[0])
			}
		} else if m.Command == "GLOBALUSERSTATE" {
		} else if m.Command == "NOTICE" {
		} else if m.Command == "HOSTTARGET" {
		} else if m.Command == "ROOMSTATE" {
		} else if m.Command == "RECONNECT" {
		} else if m.Command == "PART" {
			// Server will ACK back your PART messages on success
			// Handle leaving channel, ignore leaves that aren't us
			if m.User == c.CurrentNick() {
				log.Debug("Leaving the following channel: ", m.Params[0])
				(*channels).Remove(m.Params[0])
			}
		} else {
			log.Debug("Unknown: ", m)
		}
	}
}
