package routes

import (
	"net/http"
	"rest-api-practice/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "couldn't parse event id"})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "couldn't fetch event"})
		return
	}

	err = event.Register(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "couldn't register user for event"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "registered"})
}

func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "couldn't parse event id"})
		return
	}

	var event models.Event
	event.ID = eventId

	err = event.Unregister(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "couldn't unregister event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "unregistered"})
}
