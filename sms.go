package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go/twiml"
)

func HandleSMS(ginCtx *gin.Context) {

	// Log out ginCtx.Request.Body to see the incoming request
	body, err := ioutil.ReadAll(ginCtx.Request.Body)
	if err != nil {
		ginCtx.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Println(string(body))

	photo, _ := getPhoto()
	customCaption, _ := getCaption(photo.Caption)

	messageBody := twiml.MessagingBody{
		Message: customCaption + "\n\n" + "Photo by " + photo.Photographer + " " + "\n" + photo.PhotographerUrl + ".",
	}
	messageMedia := twiml.MessagingMedia{
		Url: photo.Url,
	}
	message := &twiml.MessagingMessage{
		InnerElements: []twiml.Element{messageBody, messageMedia},
	}

	twimlResult, err := twiml.Messages([]twiml.Element{message})
	if err != nil {
		ginCtx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ginCtx.Header("Content-Type", "text/xml")
	ginCtx.String(http.StatusOK, twimlResult)
}
