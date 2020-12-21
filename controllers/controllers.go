package controllers

import (
	"net/http"

	set "github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

// ListChannels -
// 				Returns handler that dumps current channel list
// 				Takes in a pointer to the set of current channels
func ListChannels(channels *set.Set) gin.HandlerFunc {
	return func(context *gin.Context) {
		log.Debug("Current channel list: ", (*channels).String())
		context.JSON(http.StatusOK, *channels)
	}
}

// JoinChannel 	-
//				Attempts to join the channel specified in the path variable
//				200 response indicates request was accepted
//				Need to check /channels to verify
func JoinChannel(messages chan string) gin.HandlerFunc {
	return func(context *gin.Context) {
		channel := context.Param("channel")
		messages <- "JOIN #" + channel
		context.String(http.StatusOK, "Attempting to join: %s", channel)
	}
}

// LeaveChannel -
//				Attempts to leave the channel specified in the path variable
//				200 response indicates request was accepted
//				Need to check /channels to verify
func LeaveChannel(messages chan string) gin.HandlerFunc {
	return func(context *gin.Context) {
		channel := context.Param("channel")
		messages <- "PART #" + channel
		context.String(http.StatusOK, "Attempting to leave: %s", channel)
	}
}
