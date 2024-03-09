package main

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go/client"
	"github.com/twilio/twilio-go/twiml"
)

var twilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")

func HandleSMS(ginCtx *gin.Context) {

	if twilioAuthToken == "" {
		ginCtx.String(http.StatusInternalServerError, "Twilio Auth Token not configured")
		return
	}

	requestValidator := client.NewRequestValidator(twilioAuthToken)

	// Since we need to read the body for validation and possibly again later,
	// read it into a buffer and then replace the request body with a new reader.
	var buf bytes.Buffer
	tee := io.TeeReader(ginCtx.Request.Body, &buf)
	body, err := io.ReadAll(tee)
	if err != nil {
		ginCtx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Replace the request body so it can be read again later.
	ginCtx.Request.Body = io.NopCloser(&buf)

	// Construct the full URL
	fullURL := "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI

	// Parse the query string
	queryString, err := url.ParseQuery(string(body))
	if err != nil {
		ginCtx.String(http.StatusBadRequest, "Error parsing request body")
		return
	}

	// Get Twilio signature from the request headers
	twilioSig := ginCtx.GetHeader("X-Twilio-Signature")

	// Convert url.Values to map
	queryStringMap := make(map[string]string)
	for key, values := range queryString {
		queryStringMap[key] = values[0]
	}

	// Validate the request using the request validator
	if !requestValidator.Validate(fullURL, queryStringMap, twilioSig) {
		ginCtx.String(http.StatusForbidden, "")
		return
	}

	photo, _ := getPhoto()
	customCaption, _ := getCaption(photo.Caption, queryString.Get("Body"))

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
